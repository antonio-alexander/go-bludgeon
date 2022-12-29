package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	file "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	service "github.com/antonio-alexander/go-bludgeon/timers/service/rest"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_logic.json"

var (
	configMetaMysql          = new(internal_mysql.Configuration)
	configMetaFile           = new(internal_file.Configuration)
	configLogic              = new(logic.Configuration)
	configServer             = new(internal_server.Configuration)
	configChangesClientRest  = new(changesclientrest.Configuration)
	configChangesClientKafka = new(changesclientkafka.Configuration)
	configKafkaClient        = new(internal_kafka.Configuration)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configKafkaClient.Default()
	configKafkaClient.FromEnv(envs)
	configChangesClientKafka.Default()
	configChangesClientRest.Default()
	configChangesClientRest.FromEnv(envs)
	configLogic.Default()
	configLogic.FromEnv(envs)
	configServer.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "7998"
	configServer.ShutdownTimeout = 10 * time.Second
}

type restServiceTest struct {
	server interface {
		internal.Initializer
		internal.Configurer
	}
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
	}
	changesHandler interface {
		internal.Initializer
		internal.Configurer
	}
	logic interface {
		internal.Initializer
		internal.Configurer
	}
	client *http.Client
}

func newRestServerTest(metaType internal_meta.Type, protocol string) *restServiceTest {
	var timerMeta interface {
		meta.Timer
		meta.TimeSlice
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}
	var changesClient interface {
		changesclient.Client
		changesclient.Handler
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}
	var changesHandler interface {
		changesclient.Handler
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}

	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "bludgeon_rest_server_test",
		Level:  internal_logger.Trace,
	})
	switch metaType {
	default:
		timerMeta = memory.New()
	case internal_meta.TypeMySQL:
		timerMeta = mysql.New()
	case internal_meta.TypeFile:
		timerMeta = file.New()
	}
	timerMeta.SetUtilities(logger)
	switch protocol {
	default: //rest
		c := changesclientrest.New()
		changesClient, changesHandler = c, c
	case "kafka":
		c, h := changesclientrest.New(), changesclientkafka.New()
		changesClient, changesHandler = c, h
	}
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	timerLogic := logic.New()
	timerLogic.SetParameters(timerMeta, changesClient, changesHandler)
	timerLogic.SetUtilities(logger)
	server := internal_server.New()
	timerService := service.New()
	server.SetUtilities(logger)
	server.SetParameters(timerService)
	timerService.SetUtilities(logger)
	timerService.SetParameters(timerLogic, server)
	return &restServiceTest{
		server:         server,
		meta:           timerMeta,
		changesClient:  changesClient,
		changesHandler: changesHandler,
		logic:          timerLogic,
		client:         new(http.Client),
	}
}

func (r *restServiceTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s%s", configServer.Address, configServer.Port, route)
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

func (r *restServiceTest) Initialize(t *testing.T, metaType internal_meta.Type, protocol string) {
	switch metaType {
	case internal_meta.TypeMySQL:
		err := r.meta.Configure(configMetaMysql)
		assert.Nil(t, err)
	case internal_meta.TypeFile:
		err := r.meta.Configure(configMetaFile)
		assert.Nil(t, err)
	}
	switch protocol {
	default: //rest
		err := r.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = r.changesClient.Initialize()
		assert.Nil(t, err)
	case "kafka":
		err := r.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = r.changesClient.Initialize()
		assert.Nil(t, err)
		err = r.changesHandler.Configure(configChangesClientKafka, configKafkaClient)
		assert.Nil(t, err)
		err = r.changesHandler.Initialize()
		assert.Nil(t, err)
	}
	err := r.meta.Initialize()
	assert.Nil(t, err)
	err = r.logic.Configure(configLogic)
	assert.Nil(t, err)
	err = r.logic.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Initialize()
	assert.Nil(t, err)
}

func (r *restServiceTest) Shutdown(t *testing.T) {
	r.server.Shutdown()
	r.meta.Shutdown()
}

func (r *restServiceTest) TestTimerOperations(t *testing.T) {
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

func testTimersRestService(t *testing.T, metaType internal_meta.Type, protocol string) {
	r := newRestServerTest(metaType, protocol)

	r.Initialize(t, metaType, protocol)
	defer r.Shutdown(t)

	t.Run("Timer Operations", r.TestTimerOperations)
}

func TestTimersRestServiceMemory(t *testing.T) {
	testTimersRestService(t, internal_meta.TypeMemory, "rest")
}

func TestTimersRestServiceFile(t *testing.T) {
	testTimersRestService(t, internal_meta.TypeFile, "rest")
}

func TestTimersRestServiceMySQL(t *testing.T) {
	testTimersRestService(t, internal_meta.TypeMySQL, "rest")
}
