package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logic "github.com/antonio-alexander/go-bludgeon/internal/logic"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	server "github.com/antonio-alexander/go-bludgeon/server"

	meta_memory "github.com/antonio-alexander/go-bludgeon/meta/memory"
	server_rest "github.com/antonio-alexander/go-bludgeon/server/rest"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	internal_logic "github.com/antonio-alexander/go-bludgeon/internal/logic/simple"

	"github.com/stretchr/testify/assert"
)

var (
	address string        = "localhost"
	port    string        = "8080"
	timeout time.Duration = 15 * time.Second
)

type restServerTest struct {
	server interface {
		server_rest.Owner
		server.Owner
		logic.Logic
	}
	meta interface {
		meta.Owner
		meta.Serializer
		meta.TimeSlice
		meta.Timer
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
	if _, ok := envs["BLUDGEON_REST_TIMEOUT"]; ok {
		if i, err := strconv.Atoi(envs["BLUDGEON_REST_TIMEOUT"]); err != nil {
			timeout = time.Duration(i) * time.Second
		}
	}
}

func new() *restServerTest {
	meta := meta_memory.New()
	logger := internal_logger.New("bludgeon_rest_server_test")
	logic := internal_logic.New(logger, meta)
	return &restServerTest{
		server: server_rest.New(logger, logic),
		meta:   meta,
		client: &http.Client{},
	}
}

func (r *restServerTest) Initialize(t *testing.T) {
	if err := r.server.Start(&server_rest.Configuration{
		Address: address,
		Port:    port,
		Timeout: timeout,
	}); err != nil {
		t.Logf("error while starting server: %s", err)
	}
}

func (r *restServerTest) Shutdown(t *testing.T) {
	if err := r.meta.Shutdown(); err != nil {
		t.Logf("error while shutting down meta: %s", err)
	}
	r.server.Stop()
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
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return bytes, response.StatusCode, err
}

func (r *restServerTest) testTimerCRUD(t *testing.T) {
	//read timer that doesn't exist
	bytes, err := json.Marshal(&data.Contract{
		ID: "", //KIM: this should never be a valid id
	})
	assert.Nil(t, err)
	bytes, statusCode, err := r.doRequest(data.RouteTimerRead, http.MethodPost, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.NotEmpty(t, string(bytes))

	//create timer
	bytes, statusCode, err = r.doRequest(data.RouteTimerCreate, http.MethodPost, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	timerCreated := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerCreated)
		assert.Nil(t, err)
		assert.NotEqual(t, "", timerCreated.UUID)
		//REVIEW: should we assert more fields?
	}

	//read timer
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerRead, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timerRead := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerRead)
		assert.Nil(t, err)
		assert.Equal(t, timerCreated, timerRead)
	}

	//update timer comment
	comment := "Hello, World!"
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
		Timer: data.Timer{
			UUID:    timerCreated.UUID,
			Comment: comment,
		},
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerUpdate, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timerUpdated := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerUpdated)
		assert.Nil(t, err)
		assert.Equal(t, &data.Timer{
			UUID:    timerCreated.UUID,
			Comment: comment,
		}, timerUpdated)
	}

	//read timer
	bytes, err = json.Marshal(&data.Contract{
		ID: timerUpdated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerRead, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timerRead = &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerRead)
		assert.Nil(t, err)
		assert.Equal(t, timerUpdated, timerRead)
	}

	//delete timer
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	_, statusCode, err = r.doRequest(data.RouteTimerDelete, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)

	//read timer
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerRead, http.MethodPost, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.NotEmpty(t, string(bytes))

	//delete timer again
	//REVIEW: confirm why this doesn't return an error
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerDelete, http.MethodPost, bytes)
	assert.Nil(t, err)
	// assert.Equal(t, http.StatusInternalServerError, statusCode)
	// assert.NotEmpty(t, string(bytes))
}

func (r *restServerTest) testTimerOperations(t *testing.T) {
	//create timer
	bytes, statusCode, err := r.doRequest(data.RouteTimerCreate, http.MethodPost, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	timerCreated := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerCreated)
		assert.Nil(t, err)
		assert.NotEqual(t, "", timerCreated.UUID)
		//REVIEW: should we assert more fields?
	}

	//read timer
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerRead, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timerRead := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timerRead)
		assert.Nil(t, err)
		assert.Equal(t, timerCreated.UUID, timerRead.UUID)
		assert.Equal(t, int64(0), timerCreated.ElapsedTime)
		assert.Empty(t, timerRead.ActiveSliceUUID)
	}

	//start
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerStart, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timer := &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timer)
		assert.Nil(t, err)
		assert.NotEmpty(t, timer.ActiveSliceUUID)
	}

	//wait a second so the elapsed time is greater
	time.Sleep(time.Second)

	//pause
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerPause, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timer = &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timer)
		assert.Nil(t, err)
		assert.Empty(t, timer.ActiveSliceUUID)
		assert.Greater(t, timer.ElapsedTime, int64(time.Second))
	}

	//start
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	_, statusCode, err = r.doRequest(data.RouteTimerStart, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//submit
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(data.RouteTimerSubmit, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timer = &data.Timer{}
	if assert.Nil(t, err) {
		err = json.Unmarshal(bytes, timer)
		assert.Nil(t, err)
		assert.Empty(t, timer.ActiveSliceUUID)
		assert.Greater(t, timer.ElapsedTime, int64(2*time.Second))
		assert.True(t, timer.Completed)
	}
	//TODO: add code to submit without a finish time?

	//delete
	bytes, err = json.Marshal(&data.Contract{
		ID: timerCreated.UUID,
	})
	assert.Nil(t, err)
	_, statusCode, err = r.doRequest(data.RouteTimerDelete, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

func TestRestServer(t *testing.T) {
	r := new()
	r.Initialize(t)
	t.Run("Test Timer CRUD", r.testTimerCRUD)
	t.Run("Test Timer Operations", r.testTimerOperations)
	time.Sleep(time.Second)
	r.Shutdown(t)
}
