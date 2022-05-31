package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/logic"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"
	"github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/timers/service/rest"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/stretchr/testify/assert"
)

var config *server.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config = new(server.Configuration)
	config.Default()
	config.Address = "127.0.0.1"
	config.Port = "8081"
	config.FromEnv(pwd, envs)
}

type restServerTest struct {
	server interface {
		server.Owner
		server.Router
	}
	meta interface {
		meta.Serializer
		meta.Timer
		meta.TimeSlice
	}
	logic interface {
		logic.Logic
	}
	client *http.Client
}

func newRestServerTest() *restServerTest {
	logger := logger.New("bludgeon_rest_server_test")
	server := server.New(logger)
	employeeMeta := memory.New()
	employeeLogic := logic.New(logger, employeeMeta)
	rest.New(logger, server, employeeLogic)
	return &restServerTest{
		server: server,
		meta:   employeeMeta,
		logic:  employeeLogic,
		client: &http.Client{},
	}
}

func (r *restServerTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s%s", config.Address, config.Port, route)
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

func (r *restServerTest) Initialize(t *testing.T) {
	if err := r.server.Start(config); err != nil {
		t.Logf("error while starting server: %s", err)
	}
}

func (r *restServerTest) Shutdown(t *testing.T) {
	r.server.Stop()
}

func (r *restServerTest) TestTimerOperations(t *testing.T) {
	//create timer
	bytes, err := json.Marshal(&data.TimerPartial{})
	assert.Nil(t, err)
	bytes, statusCode, err := r.doRequest(data.RouteTimers, http.MethodPost, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timerCreated := &data.Timer{}
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, timerCreated)
	assert.Nil(t, err)
	assert.NotEmpty(t, timerCreated.ID)
	//REVIEW: should we assert more fields?

	//read timer
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDf, timerCreated.ID), http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	timerRead := &data.Timer{}
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, timerRead)
	assert.Nil(t, err)
	assert.Equal(t, timerCreated.ID, timerRead.ID)
	assert.Equal(t, int64(0), timerCreated.ElapsedTime)
	assert.Empty(t, timerRead.ActiveTimeSliceID)

	// read multiple timers
	bytes, statusCode, err = r.doRequest(data.RouteTimersSearch, http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	timers := []data.Timer{}
	err = json.Unmarshal(bytes, &timers)
	assert.Nil(t, err)
	assert.Contains(t, timers, *timerRead)

	//start
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDStartf, timerCreated.ID), http.MethodPut, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	timer := &data.Timer{}
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, timer)
	assert.Nil(t, err)
	assert.NotEmpty(t, timer.ActiveTimeSliceID)

	//wait a second so the elapsed time is greater
	time.Sleep(time.Second)

	//stop
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDStopf, timerCreated.ID), http.MethodPut, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	timer = &data.Timer{}
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, timer)
	assert.Nil(t, err)
	assert.Empty(t, timer.ActiveTimeSliceID)
	assert.Greater(t, timer.ElapsedTime, int64(time.Second))

	//start
	assert.Nil(t, err)
	_, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDStartf, timerCreated.ID), http.MethodPut, nil)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//submit
	bytes, err = json.Marshal(&data.Contract{
		Finish: time.Now().UnixNano(),
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDSubmitf, timerCreated.ID), http.MethodPut, bytes)
	assert.Equal(t, http.StatusOK, statusCode)
	timer = &data.Timer{}
	assert.Nil(t, err)
	err = json.Unmarshal(bytes, timer)
	assert.Nil(t, err)
	assert.Empty(t, timer.ActiveTimeSliceID)
	assert.Greater(t, timer.ElapsedTime, int64(2*time.Second))
	assert.True(t, timer.Completed)
	//TODO: add code to submit without a finish time?

	//delete
	assert.Nil(t, err)
	_, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteTimersIDf, timerCreated.ID), http.MethodDelete, nil)
	assert.Equal(t, http.StatusNoContent, statusCode)
	assert.Nil(t, err)
}

func TestTimersRestService(t *testing.T) {
	r := newRestServerTest()
	r.Initialize(t)
	t.Run("Timer Operations", r.TestTimerOperations)
	r.Shutdown(t)
}
