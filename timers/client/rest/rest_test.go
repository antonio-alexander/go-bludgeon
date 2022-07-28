package rest_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	rest "github.com/antonio-alexander/go-bludgeon/timers/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"

	"github.com/stretchr/testify/assert"
)

var config *restclient.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config = new(restclient.Configuration)
	config.Default()
	config.FromEnv(pwd, envs)
}

type restClientTest struct {
	client rest.Client
}

func newRestClientTest() *restClientTest {
	logger := logger.New("bludgeon_rest_server_test")
	client := rest.New(logger, config)
	return &restClientTest{
		client: client,
	}
}

func (r *restClientTest) Initialize() error {
	return r.client.Initialize(config)
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
	assert.GreaterOrEqual(t, int64(time.Second), timerRead.ElapsedTime)
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
	timerSubmitted, err := r.client.TimerSubmit(ctx, timerID, &tNow)
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
	err := r.Initialize()
	assert.Nil(t, err)
	t.Run("Test Timer Operations", r.TestTimers)
}
