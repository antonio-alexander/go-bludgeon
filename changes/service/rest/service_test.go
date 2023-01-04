package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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

var configServer = new(internal_server.Configuration)

type restServerTest struct {
	server interface {
		internal.Initializer
		internal.Configurer
	}
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	logic interface {
		internal.Initializer
	}
	client *http.Client
}

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configServer.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "9000"
	configServer.AllowedMethods = []string{http.MethodDelete, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodGet}
	configServer.ShutdownTimeout = 15 * time.Second
}

func newRestServerTest() *restServerTest {
	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	changesMeta := memory.New()
	changesMeta.SetUtilities(logger)
	changesLogic := logic.New()
	changesLogic.SetUtilities(logger)
	changesLogic.SetParameters(changesMeta)
	changesService := service.New()
	changesService.SetUtilities(logger)
	server := internal_server.New()
	server.SetUtilities(logger)
	changesService.SetParameters(changesLogic, server)
	server.SetParameters(changesService)
	return &restServerTest{
		server: server,
		meta:   changesMeta,
		logic:  changesLogic,
		client: &http.Client{},
	}
}

func (l *restServerTest) generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func (r *restServerTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s"+route, configServer.Address, configServer.Port)
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
	err := r.meta.Initialize()
	assert.Nil(t, err)
	err = r.meta.Configure()
	assert.Nil(t, err)
	err = r.logic.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Initialize()
	assert.Nil(t, err)
	//KIM: we have to sleep here because the start for the rest
	// server isn't synchronous
	time.Sleep(2 * time.Second)
}

func (r *restServerTest) shutdown(t *testing.T) {
	//REVIEW: why do we need to sleep?
	time.Sleep(2 * time.Second)
	r.server.Shutdown()
	r.logic.Shutdown()
	r.meta.Shutdown()
}

func (r *restServerTest) testChangeOperations(t *testing.T) {
	//generate constants
	dataId, version := r.generateId(), 1
	dataType, dataService := "employee", "employees"

	//create change
	bytes, err := json.Marshal(&data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &version,
		DataType:        &dataType,
		DataServiceName: &dataService,
	})
	assert.Nil(t, err)
	uri := data.RouteChanges
	bytes, statusCode, err := r.doRequest(uri, data.MethodChangeUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeUpserted := &data.Change{}
	err = json.Unmarshal(bytes, changeUpserted)
	assert.Nil(t, err)
	changeId := changeUpserted.Id

	//read change
	uri = fmt.Sprintf(data.RouteChangesParamf, changeId)
	bytes, statusCode, err = r.doRequest(uri, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeRead := &data.Change{}
	err = json.Unmarshal(bytes, changeRead)
	assert.Nil(t, err)
	assert.Equal(t, changeId, changeRead.Id)
	assert.Equal(t, dataId, changeRead.DataId)
	assert.Equal(t, version, changeRead.DataVersion)
	assert.Equal(t, dataType, changeRead.DataType)
	assert.Equal(t, dataService, changeRead.DataServiceName)

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
	assert.Equal(t, meta.ErrChangeNotFound.Error(), internalErr.Error())

	//TODO: add test to acknowledge change
}

func (r *restServerTest) testChangeStreaming(t *testing.T) {
	var wg sync.WaitGroup

	//generate change
	dataId, version := r.generateId(), 1
	dataType, dataService := "employee", "employees"

	//connect to web socket
	websocketUri := fmt.Sprintf("ws://%s:%s"+data.RouteChangesWebsocket, configServer.Address, configServer.Port)
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
		bytes, err := json.Marshal(&data.ChangePartial{
			DataId:          &dataId,
			DataVersion:     &version,
			DataType:        &dataType,
			DataServiceName: &dataService,
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
	var changes []*data.Change

	//upsert change (neither registration should see this change)
	dataId := r.generateId()
	dataVersion, dataType := rand.Intn(1000), r.generateId()
	dataServiceName, whenChanged := r.generateId(), time.Now().UnixNano()
	changedBy, dataAction := "test_change_crud", r.generateId()
	bytes, err := json.Marshal(&data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &dataServiceName,
		DataAction:      &dataAction,
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err := r.doRequest(data.RouteChanges, data.MethodChangeUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeCreated := &data.Change{}
	err = json.Unmarshal(bytes, changeCreated)
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		route := fmt.Sprintf(data.RouteChangesParamf, changeId)
		r.doRequest(route, data.MethodChangeDelete, bytes)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//create registration (1)
	registrationId1 := r.generateId()
	bytes, err = json.Marshal(&data.RequestRegister{
		RegistrationId: registrationId1,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err = r.doRequest(data.RouteChangesRegistration, data.MethodRegistrationUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	defer func(registrationId string) {
		route := fmt.Sprintf(data.RouteChangesRegistrationParamf, registrationId)
		r.doRequest(route, data.MethodRegistrationDelete, bytes)
	}(registrationId1)

	//validate that registration (1) doesn't include any changes
	route := fmt.Sprintf(data.RouteChangesRegistrationParamChangesf, registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest := &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 0)

	//upsert change (to be seen by registration (1))
	dataId = r.generateId()
	dataVersion, dataType = rand.Intn(1000), r.generateId()
	dataServiceName, whenChanged = r.generateId(), time.Now().UnixNano()
	dataAction, changedBy = r.generateId(), "test_change_crud"
	bytes, err = json.Marshal(&data.ChangePartial{
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
		DataId:          &dataId,
		DataServiceName: &dataServiceName,
		DataType:        &dataType,
		DataAction:      &dataAction,
		DataVersion:     &dataVersion,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err = r.doRequest(data.RouteChanges, data.MethodChangeUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeCreated = &data.Change{}
	err = json.Unmarshal(bytes, changeCreated)
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		route := fmt.Sprintf(data.RouteChangesParamf, changeId)
		r.doRequest(route, data.MethodChangeDelete, bytes)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (1) sees the second change, but not the first
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf, registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 1)
	assert.Contains(t, changeDigest.Changes, changes[1])
	assert.NotContains(t, changeDigest.Changes, changes[0])

	//create registration
	registrationId2 := r.generateId()
	bytes, err = json.Marshal(&data.RequestRegister{
		RegistrationId: registrationId2,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err = r.doRequest(data.RouteChangesRegistration, data.MethodRegistrationUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	defer func(registrationId string) {
		route := fmt.Sprintf(data.RouteChangesRegistrationParamf, registrationId)
		r.doRequest(route, data.MethodRegistrationDelete, bytes)
	}(registrationId2)

	//upsert change (to be seen by registration (1) and (2))
	dataId = r.generateId()
	dataVersion, dataType = rand.Intn(1000), r.generateId()
	dataServiceName, whenChanged = r.generateId(), time.Now().UnixNano()
	dataAction, changedBy = r.generateId(), "test_change_crud"
	bytes, err = json.Marshal(&data.ChangePartial{
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
		DataId:          &dataId,
		DataServiceName: &dataServiceName,
		DataType:        &dataType,
		DataAction:      &dataAction,
		DataVersion:     &dataVersion,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err = r.doRequest(data.RouteChanges, data.MethodChangeUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeCreated = &data.Change{}
	err = json.Unmarshal(bytes, changeCreated)
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		route := fmt.Sprintf(data.RouteChangesParamf, changeId)
		r.doRequest(route, data.MethodChangeDelete, bytes)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (2) sees the third change, but not the first or second
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf, registrationId2)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 1)
	assert.Contains(t, changeDigest.Changes, changes[2])
	assert.NotContains(t, changeDigest.Changes, changes[0])
	assert.NotContains(t, changeDigest.Changes, changes[1])

	//acknowledge the initial change for registration 1
	bytes, err = json.Marshal(&data.RequestAcknowledge{
		ChangeIds: []string{changes[0].Id},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	route = fmt.Sprintf(data.RouteChangesRegistrationServiceIdAcknowledgef, registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodRegistrationChangeAcknowledge, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	//get changes for registration 1
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf, registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 2)
	//acknowledge the initial change for registration 1
	bytes, err = json.Marshal(&data.RequestAcknowledge{
		ChangeIds: []string{changes[0].Id},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	//get changes for registration 2
	route = fmt.Sprintf(data.RouteChangesRegistrationServiceIdAcknowledgef, registrationId2)
	bytes, statusCode, err = r.doRequest(route, data.MethodRegistrationChangeAcknowledge, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf, registrationId2)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 1)

	//validate that initial change has been removed
	//KIM: this change is removed because of the acknowledgement(s) earlier
	route = fmt.Sprintf(data.RouteChangesParamf, changes[0].Id)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, statusCode)

	//acknowledge the second change for both registrations and confirm actual changes
	bytes, err = json.Marshal(&data.RequestAcknowledge{
		ChangeIds: []string{changes[1].Id},
	})
	assert.NotEmpty(t, bytes)
	assert.Nil(t, err)
	route = fmt.Sprintf(data.RouteChangesRegistrationServiceIdAcknowledgef,
		registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodRegistrationChangeAcknowledge, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	//read changes for registration 1
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf,
		registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 1)
	assert.Contains(t, changeDigest.Changes, changes[2])

	//acknowledge the second change for registration 2
	bytes, err = json.Marshal(&data.RequestAcknowledge{
		ChangeIds: []string{changes[1].Id},
	})
	assert.NotEmpty(t, bytes)
	assert.Nil(t, err)
	route = fmt.Sprintf(data.RouteChangesRegistrationServiceIdAcknowledgef,
		registrationId2)
	bytes, statusCode, err = r.doRequest(route, data.MethodRegistrationChangeAcknowledge, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	// read the changes for registration 2
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf,
		registrationId2)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Len(t, changeDigest.Changes, 1)
	assert.Contains(t, changeDigest.Changes, changes[2])

	//validate that second change has been removed
	route = fmt.Sprintf(data.RouteChangesParamf, changes[1].Id)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, statusCode)

	//delete the initial registration,
	route = fmt.Sprintf(data.RouteChangesRegistrationParamf, registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodRegistrationDelete, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	// re-create the registration
	bytes, err = json.Marshal(&data.RequestRegister{
		RegistrationId: registrationId1,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	bytes, statusCode, err = r.doRequest(data.RouteChangesRegistration, data.MethodRegistrationUpsert, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Empty(t, bytes)
	// read the registration changes
	route = fmt.Sprintf(data.RouteChangesRegistrationParamChangesf,
		registrationId1)
	bytes, statusCode, err = r.doRequest(route, data.MethodChangeRead, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.NotEmpty(t, bytes)
	changeDigest = &data.ChangeDigest{}
	err = json.Unmarshal(bytes, changeDigest)
	assert.Nil(t, err)
	assert.Empty(t, changeDigest.Changes)
}

func TestChangesRestService(t *testing.T) {
	r := newRestServerTest()

	r.initialize(t)
	defer r.shutdown(t)

	t.Run("Change Operations", r.testChangeOperations)
	t.Run("Change Streaming", r.testChangeStreaming)
	t.Run("Change Registration", r.testChangeRegistration)
}
