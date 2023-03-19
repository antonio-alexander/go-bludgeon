package rest_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	restclient "github.com/antonio-alexander/go-bludgeon/timers/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/stretchr/testify/assert"
)

var config = new(restclient.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.Port = "8080"
}

type restClientTest struct {
	client interface {
		client.Client
		internal.Parameterizer
		internal.Configurer
		internal.Initializer
	}
}

func newRestClientTest() *restClientTest {
	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "bludgeon_rest_server_test",
		Level:  internal_logger.Trace,
	})
	client := restclient.New()
	client.SetUtilities(logger)
	return &restClientTest{
		client: client,
	}
}

func (r *restClientTest) Initialize(t *testing.T) {
	err := r.client.Configure(config)
	assert.Nil(t, err)
	err = r.client.Initialize()
	assert.Nil(t, err)
}

func (r *restClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
}

func (r *restClientTest) TestTimers(t *testing.T) {
	ctx := context.TODO()

	//create a timer
	timer, err := r.client.TimerCreate(ctx, data.TimerPartial{})
	assert.Nil(t, err)
	assert.NotEmpty(t, timer.ID)
	assert.Empty(t, timer.ActiveTimeSliceID)
	assert.False(t, timer.Completed)
	assert.False(t, timer.Archived)
	timerID := timer.ID
	//read the timer
	timerRead, err := r.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timer, timerRead)
	//start the timer
	timerStarted, err := r.client.TimerStart(ctx, timerID)
	assert.Nil(t, err)
	assert.NotEmpty(t, timerStarted.ActiveTimeSliceID)
	assert.NotZero(t, timerStarted.Start)
	//wait for a second
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = r.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, timerRead.ElapsedTime, int64(time.Second))
	//stop the timer
	timerStopped, err := r.client.TimerStop(ctx, timerID)
	assert.Nil(t, err)
	//read the timer
	timerRead, err = r.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)
	//wait one second
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = r.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)
	//submit the timer
	tNow := time.Now()
	timerSubmitted, err := r.client.TimerSubmit(ctx, timerID, tNow.UnixNano())
	assert.Nil(t, err)
	//read the timer
	timerRead, err = r.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerSubmitted, timerRead)
	//delete the timer
	err = r.client.TimerDelete(ctx, timerID)
	assert.Nil(t, err)
}

func TestTimersRestClient(t *testing.T) {
	r := newRestClientTest()
	r.Initialize(t)
	defer r.Shutdown(t)

	t.Run("Test Timer Operations", r.TestTimers)
}
