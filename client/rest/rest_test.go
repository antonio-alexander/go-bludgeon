package rest_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/client/rest"
	"github.com/stretchr/testify/assert"
)

var (
	validConfig *rest.Configuration
	pwd         string
	envs        map[string]string
)

func init() {
	pwd, _ = os.Getwd()
	envs = make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	v := &rest.Configuration{}
	v.Default()
	v.FromEnv(pwd, envs)
	validConfig = v
}

func TestTimer(t *testing.T) {
	var timerID string

	client := rest.New(*validConfig)
	//create a timer
	timer, err := client.TimerCreate()
	assert.Nil(t, err)
	assert.Condition(t, func() bool {
		if timer.UUID == "" {
			return false
		}
		if timer.ActiveSliceUUID != "" {
			return false
		}
		if timer.Completed {
			return false
		}
		return true
	})
	timerID = timer.UUID
	//read the timer
	timerRead, err := client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timer, timerRead)
	//start the timer
	tStart := time.Now()
	timerStart, err := client.TimerStart(timerID, tStart)
	assert.Nil(t, err)
	assert.Condition(t, func() bool {
		if timerStart.ActiveSliceUUID == "" {
			return false
		}
		if timerStart.Start != tStart.UnixNano() {
			return false
		}
		return true
	})
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Condition(t, func() bool {
		return timerRead.ElapsedTime >= int64(time.Second)
	})
	//pause the timer
	timerPaused, err := client.TimerPause(timerID, time.Now())
	assert.Nil(t, err)
	//read the timer
	timerRead, err = client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerPaused, timerRead)
	time.Sleep(time.Second)
	//read the timer
	timerRead, err = client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerPaused, timerRead)
	//submit the timer
	timerSubmitted, err := client.TimerSubmit(timerID, time.Now())
	assert.Nil(t, err)
	//read the timer
	timerRead, err = client.TimerRead(timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerSubmitted, timerRead)
}
