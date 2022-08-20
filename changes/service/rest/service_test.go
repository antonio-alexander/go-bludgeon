package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	logic "github.com/antonio-alexander/go-bludgeon/changes/logic"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	memory "github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	service "github.com/antonio-alexander/go-bludgeon/changes/service/rest"
	internal "github.com/antonio-alexander/go-bludgeon/internal"

	internal_error "github.com/antonio-alexander/go-bludgeon/internal/errors"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	address         string        = "localhost"
	port            string        = "8082"
	shutdownTimeout time.Duration = 15 * time.Second
)

type restServerTest struct {
	server interface {
		internal.Initializer
		internal.Configurer
		internal_server.Router
	}
	meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		meta.Serializer
		internal.Shutdowner
	}
	logic interface {
		logic.Logic
		internal.Initializer
	}
	client *http.Client
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	if _, ok := envs["BLUDGEON_REST_ADDRESS"]; ok {
		address = envs["BLUDGEON_REST_ADDRESS"]
	}
	if _, ok := envs["BLUDGEON_REST_PORT"]; ok {
		port = envs["BLUDGEON_REST_PORT"]
	}
	if _, ok := envs["BLUDGEON_REST_SHUTDOWN_TIMEOUT"]; ok {
		if i, err := strconv.Atoi(envs["BLUDGEON_REST_SHUTDOWN_TIMEOUT"]); err != nil {
			shutdownTimeout = time.Duration(i) * time.Second
		}
	}
}

func new() *restServerTest {
	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	server := internal_server.New()
	server.SetUtilities(logger)
	employeeMeta := memory.New()
	employeeMeta.SetUtilities(logger)
	employeeLogic := logic.New()
	employeeLogic.SetUtilities(logger)
	employeeLogic.SetParameters(employeeMeta)
	service := service.New()
	service.SetUtilities(logger)
	service.SetParameters(server, employeeLogic)
	return &restServerTest{
		server: server,
		meta:   employeeMeta,
		logic:  employeeLogic,
		client: &http.Client{},
	}
}

func (l *restServerTest) generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func (r *restServerTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s%s", address, port, route)
	request, err := http.NewRequest(method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, -1, err
	}
	response, err := r.client.Do(request)
	if err != nil {
		return nil, -1, err
	}
	bytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	return bytes, response.StatusCode, err
}

func (r *restServerTest) pingHandler(t *testing.T, ws interface {
	WriteControl(messageType int, data []byte, deadline time.Time) error
}) func(string) error {
	return func(ping string) error {
		t.Logf("ping received: %s", ping)
		deadline := time.Now().Add(10 * time.Second)
		ws.WriteControl(websocket.PongMessage, []byte("pong"), deadline)
		return nil
	}
}

func (r *restServerTest) pongHandler(t *testing.T) func(pong string) error {
	return func(pong string) error {
		t.Logf("pong: %s", pong)
		return nil
	}
}

func (r *restServerTest) initialize(t *testing.T) {
	err := r.logic.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(&internal_server.Configuration{
		Address:         address,
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
	})
	assert.Nil(t, err)
	err = r.server.Initialize()
	assert.Nil(t, err)
	//KIM: we have to sleep here because the start for the rest
	// server isn't synchronous
	time.Sleep(2 * time.Second)
}

func (r *restServerTest) shutdown(t *testing.T) {
	r.server.Shutdown()
	r.logic.Shutdown()
	r.meta.Shutdown()
}

func (r *restServerTest) testChangeOperations(t *testing.T) {
	//generate constants
	dataId, version := r.generateId(), 1
	dataType, dataService := "employee", "employees"

	//create change
	bytes, err := json.Marshal(&data.RequestChange{
		ChangePartial: data.ChangePartial{
			DataId:          &dataId,
			DataVersion:     &version,
			DataType:        &dataType,
			DataServiceName: &dataService,
		},
	})
	assert.Nil(t, err)
	uri := data.RouteChanges
	bytes, statusCode, err := r.doRequest(uri, data.MethodChangeUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeResponse := &data.ResponseChange{}
	err = json.Unmarshal(bytes, changeResponse)
	assert.Nil(t, err)
	changeId := changeResponse.Change.Id

	//read change
	uri = fmt.Sprintf(data.RouteChangesParamf, changeId)
	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeResponse = &data.ResponseChange{}
	err = json.Unmarshal(bytes, changeResponse)
	assert.Nil(t, err)
	assert.Equal(t, changeId, changeResponse.Change.Id)
	assert.Equal(t, dataId, changeResponse.Change.DataId)
	assert.Equal(t, version, changeResponse.Change.DataVersion)
	assert.Equal(t, dataType, changeResponse.Change.DataType)
	assert.Equal(t, dataService, changeResponse.Change.DataServiceName)

	//delete change
	uri = fmt.Sprintf(data.RouteChangesParamf, changeId)
	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeDelete, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)

	//read change
	uri = fmt.Sprintf(data.RouteChangesParamf, changeId)
	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.NotEmpty(t, bytes)
	internalErr := &internal_error.Error{}
	err = json.Unmarshal(bytes, internalErr)
	assert.Nil(t, err)
	assert.NotEmpty(t, internalErr.Error)
	assert.Equal(t, meta.ErrChangeNotFound.Error(), internalErr.Error)

	//TODO: add test to acknowledge change
}

func (r *restServerTest) testChangeStreaming(t *testing.T) {
	var wg sync.WaitGroup

	//generate change
	dataId, version := r.generateId(), 1
	dataType, dataService := "employee", "employees"

	//connect to web socket
	websocketUri := fmt.Sprintf("ws://%s:%s"+data.RouteChangesWebsocket, address, port)
	ws, response, err := websocket.DefaultDialer.Dial(websocketUri, nil)
	defer response.Body.Close()
	if err == websocket.ErrBadHandshake {
		t.Logf("handshake failed with status %d", response.StatusCode)
	}
	assert.Nil(t, err)
	assert.NotNil(t, ws)
	ws.SetPingHandler(r.pingHandler(t, ws))
	ws.SetPongHandler(r.pongHandler(t))

	//start go rountes to test
	start := make(chan struct{})
	stopper := make(chan struct{})
	chChangeReceived := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-start
		bytes, err := json.Marshal(&data.RequestChange{
			ChangePartial: data.ChangePartial{
				DataId:          &dataId,
				DataVersion:     &version,
				DataType:        &dataType,
				DataServiceName: &dataService,
			},
		})
		assert.Nil(t, err)
		bytes, statusCode, err := r.doRequest(data.RouteChanges, data.MethodChangeUpsert, bytes)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.NotEmpty(t, bytes)
		changeCreated := &data.Change{}
		err = json.Unmarshal(bytes, changeCreated)
		assert.Nil(t, err)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		tRead := time.NewTicker(time.Second)
		defer tRead.Stop()
		<-start
		for {
			select {
			case <-stopper:
				return
			case <-chChangeReceived:
				return
			default:
				wrapper := &data.Wrapper{}
				err := ws.ReadJSON(wrapper)
				assert.Nil(t, err)
				switch data.MessageType(wrapper.Type) {
				case data.MessageTypeChange:
					change := &data.Change{}
					if err := json.Unmarshal(wrapper.Bytes, change); err != nil {
						break
					}
					if change.Id != "" {
						close(chChangeReceived)
					}
				}
			}
		}
	}()
	close(start)
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm message received")
	case <-chChangeReceived:
	}
	close(stopper)
	ws.Close()
	wg.Wait()
}

func (r *restServerTest) testChangeRegistration(t *testing.T) {
	// 	var changes []*data.Change
	// 	ctx := context.TODO()

	// 	//upsert change (neither registration should see this change)
	// 	dataId := r.generateId()
	// 	dataVersion, dataType := rand.Intn(1000), r.generateId()
	// 	dataServiceName, whenChanged := r.generateId(), time.Now().UnixNano()
	// 	changedBy, dataAction := "test_change_crud", r.generateId()
	// 	bytes, err := json.Marshal(&data.RequestChange{
	// 		ChangePartial: data.ChangePartial{
	// 			DataId:          &dataId,
	// 			DataVersion:     &dataVersion,
	// 			DataType:        &dataType,
	// 			DataServiceName: &dataServiceName,
	// 			DataAction:      &dataAction,
	// 			WhenChanged:     &whenChanged,
	// 			ChangedBy:       &changedBy,
	// 		},
	// 	})
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, bytes)
	// 	uri := fmt.Sprintf("http://%s:%s"+data.RouteChanges, address, port)
	// 	bytes, statusCode, err := r.doRequest(uri, data.MethodChangeUpsert, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeCreated := &data.Change{}
	// 	err = json.Unmarshal(bytes, changeCreated)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, changeCreated)
	// 	defer func(changeId string) {
	// 		uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, address, port, changeId)
	// 		r.doRequest(uri, data.MethodChangeDelete, bytes)
	// 	}(changeCreated.Id)
	// 	changes = append(changes, changeCreated)

	// 	//create registration (1)
	// 	registrationId1 := r.generateId()
	// 	bytes, err = json.Marshal(&data.RequestRegister{
	// 		RegistrationId: registrationId1,
	// 	})
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, bytes)
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistration, address, port)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodRegistrationUpsert, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	registrationResponse := &data.ResponseRegister{}
	// 	err = json.Unmarshal(bytes, registrationResponse)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, registrationId1, registrationResponse.RegistrationId)
	// 	defer func(registrationId string) {
	// 		uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamf, address, port, registrationId)
	// 		r.doRequest(uri, data.MethodRegistrationDelete, bytes)
	// 	}(registrationResponse.RegistrationId)

	// 	//validate that registration (1) doesn't include any changes
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest := &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 0)

	// 	//upsert change (to be seen by registration (1))
	// 	dataId = r.generateId()
	// 	dataVersion, dataType = rand.Intn(1000), r.generateId()
	// 	dataServiceName, whenChanged = r.generateId(), time.Now().UnixNano()
	// 	dataAction, changedBy = r.generateId(), "test_change_crud"
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChanges, address, port)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeUpsert, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeCreated = &data.Change{}
	// 	err = json.Unmarshal(bytes, changeCreated)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, changeCreated)
	// 	defer func(changeId string) {
	// 		uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, address, port, changeId)
	// 		r.doRequest(uri, data.MethodChangeDelete, bytes)
	// 	}(changeCreated.Id)
	// 	changes = append(changes, changeCreated)

	// 	//validate that registration (1) sees the second change, but not the first
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest = &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 1)
	// 	assert.Contains(t, changeDigest.Changes, changes[1])
	// 	assert.NotContains(t, changeDigest.Changes, changes[0])

	// 	//create registration
	// 	registrationId2 := r.generateId()
	// 	bytes, err = json.Marshal(&data.RequestRegister{
	// 		RegistrationId: registrationId1,
	// 	})
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, bytes)
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistration, address, port)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodRegistrationUpsert, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	registrationResponse = &data.ResponseRegister{}
	// 	err = json.Unmarshal(bytes, registrationResponse)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, registrationId2, registrationResponse.RegistrationId)
	// 	defer func(registrationId string) {
	// 		uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamf, address, port, registrationId)
	// 		r.doRequest(uri, data.MethodRegistrationDelete, bytes)
	// 	}(registrationResponse.RegistrationId)

	// 	//upsert change (to be seen by registration (1) and (2))
	// 	dataId = r.generateId()
	// 	dataVersion, dataType = rand.Intn(1000), r.generateId()
	// 	dataServiceName, whenChanged = r.generateId(), time.Now().UnixNano()
	// 	dataAction, changedBy = r.generateId(), "test_change_crud"
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChanges, address, port)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeUpsert, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeCreated = &data.Change{}
	// 	err = json.Unmarshal(bytes, changeCreated)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, changeCreated)
	// 	defer func(changeId string) {
	// 		uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, address, port, changeId)
	// 		r.doRequest(uri, data.MethodChangeDelete, bytes)
	// 	}(changeCreated.Id)
	// 	changes = append(changes, changeCreated)

	// 	//validate that registration (2) sees the third change, but not the first or second
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId2)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest = &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 1)
	// 	assert.Contains(t, changeDigest.Changes, changes[2])
	// 	assert.NotContains(t, changeDigest.Changes, changes[0])
	// 	assert.NotContains(t, changeDigest.Changes, changes[1])

	// 	//acknowledge the initial change for both services and confirm that there's no change
	// 	bytes, err = json.Marshal(&data.RequestAcknowledge{
	// 		ChangeIds: []string{changes[0].Id},
	// 	})
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationServiceIdAcknowledgef,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodRegistrationChangeAcknowledge, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	acknowledgeResponse := &data.ResponseAcknowledge{}
	// 	err = json.Unmarshal(bytes, acknowledgeResponse)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, registrationId1, acknowledgeResponse.RegistrationId)
	// 	assert.Contains(t, acknowledgeResponse.ChangeIds, changes[0].Id)
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest = &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 2)
	// 	bytes, err = json.Marshal(&data.RequestAcknowledge{
	// 		ChangeIds: []string{changes[0].Id},
	// 	})
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationServiceIdAcknowledgef,
	// 		address, port, registrationId2)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodRegistrationChangeAcknowledge, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	acknowledgeResponse = &data.ResponseAcknowledge{}
	// 	err = json.Unmarshal(bytes, acknowledgeResponse)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, registrationId1, acknowledgeResponse.RegistrationId)
	// 	assert.Contains(t, acknowledgeResponse.ChangeIds, changes[0].Id)
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId2)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest = &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 1)

	// 	//validate that initial change has been removed
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf,
	// 		address, port, changes[0].Id)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	//REVIEW: this should validate a 404 (but its probably a 500)
	// 	assert.NotEqual(t, http.StatusOK, statusCode)

	// 	//acknowledge the second change for both registrations and confirm actual changes
	// 	bytes, err = json.Marshal(&data.RequestAcknowledge{
	// 		ChangeIds: []string{changes[1].Id},
	// 	})
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationServiceIdAcknowledgef,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodRegistrationChangeAcknowledge, bytes)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	acknowledgeResponse = &data.ResponseAcknowledge{}
	// 	err = json.Unmarshal(bytes, acknowledgeResponse)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, registrationId1, acknowledgeResponse.RegistrationId)
	// 	assert.Contains(t, acknowledgeResponse.ChangeIds, changes[0].Id)
	// 	//
	// 	uri = fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
	// 		address, port, registrationId1)
	// 	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, http.StatusOK, statusCode)
	// 	assert.Empty(t, bytes)
	// 	changeDigest = &data.ChangeDigest{}
	// 	err = json.Unmarshal(bytes, changeDigest)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changeDigest.Changes, 1)
	// 	assert.Contains(t, changeDigest.Changes, changes[2])
	// 	//
	// 	err = l.RegistrationChangeAcknowledge(ctx, registrationId2, changes[1].Id)
	// 	assert.Nil(t, err)
	// 	changesRead, err = l.RegistrationChangesRead(ctx, registrationId2)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changesRead, 1)
	// 	assert.Contains(t, changesRead, changes[2])

	// 	//validate that second change has been removed
	// 	change, err = l.ChangeRead(ctx, changes[1].Id)
	// 	assert.NotNil(t, err)
	// 	assert.Nil(t, change)

	// 	//delete the initial registration, then re-create it and ensure that there are no changes
	// 	err = l.RegistrationDelete(ctx, registrationId1)
	// 	assert.Nil(t, err)
	// 	err = l.RegistrationUpsert(ctx, registrationId1)
	// 	assert.Nil(t, err)
	// 	changesRead, err = l.RegistrationChangesRead(ctx, registrationId1)
	// 	assert.Nil(t, err)
	// 	assert.Len(t, changesRead, 0)
}

func TestChangesRestService(t *testing.T) {
	r := new()

	//initialize
	r.initialize(t)
	defer func() {
		r.shutdown(t)
	}()

	//execute tests
	t.Run("Change Operations", r.testChangeOperations)
	t.Run("Change Streaming", r.testChangeStreaming)
	t.Run("Change Registration", r.testChangeRegistration)
	time.Sleep(2 * time.Second)
}
