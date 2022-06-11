package rest_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	"github.com/antonio-alexander/go-bludgeon/timers/client/rest"
	"github.com/antonio-alexander/go-bludgeon/timers/data"

	"github.com/stretchr/testify/assert"
)

var config *client.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config = new(client.Configuration)
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
	//create a timer
	timer, err := r.client.TimerCreate(data.TimerPartial{})
	assert.Nil(t, err)
	assert.NotEmpty(t, timer.ID)
	assert.Empty(t, timer.ActiveTimeSliceID)
	assert.False(t, timer.Completed)
	assert.False(t, timer.Archived)
	timerID := timer.ID
	//read the timer
	timerRead, err := r.client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timer, timerRead)
	//start the timer
	timerStarted, err := r.client.TimerStart(timerID)
	assert.Nil(t, err)
	assert.NotEmpty(t, timerStarted.ActiveTimeSliceID)
	assert.NotZero(t, timerStarted.Start)
	//wait for a second
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = r.client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, int64(time.Second), timerRead.ElapsedTime)
	//stop the timer
	timerStopped, err := r.client.TimerStop(timerID)
	assert.Nil(t, err)
	//read the timer
	timerRead, err = r.client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)
	//wait one second
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = r.client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)
	//submit the timer
	tNow := time.Now()
	timerSubmitted, err := r.client.TimerSubmit(timerID, &tNow)
	assert.Nil(t, err)
	//read the timer
	timerRead, err = r.client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerSubmitted, timerRead)
	//delete the timer
	err = r.client.TimerDelete(timerID)
	assert.Nil(t, err)
}

func TestTimersRestClient(t *testing.T) {
	r := newRestClientTest()
	err := r.Initialize()
	assert.Nil(t, err)
	t.Run("Test Employee Operations", r.TestTimers)
}
