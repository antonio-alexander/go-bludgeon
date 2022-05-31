package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestTimerCRUD(m interface {
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create timer
		comment := randomString(25)
		timer, err := m.TimerCreate(data.TimerPartial{
			Comment: &comment,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, timer.ID)
		assert.Equal(t, timer.Comment, comment)
		//read
		timerRead, err := m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timer, timerRead)
		timers, err := m.TimersRead(data.TimerSearch{
			Completed: new(bool),
		})
		assert.Nil(t, err)
		assert.Contains(t, timers, timer)
		//update
		updatedComment := randomString(25)
		timerUpdated, err := m.TimerUpdate(timer.ID, data.TimerPartial{
			Comment: &updatedComment,
		})
		assert.Nil(t, err)
		assert.Equal(t, updatedComment, timerUpdated.Comment)
		timer.Comment = timerUpdated.Comment
		timer.LastUpdated = timerUpdated.LastUpdated
		timer.LastUpdatedBy = timerUpdated.LastUpdatedBy
		timer.Version = timerUpdated.Version
		assert.Equal(t, timer, timerUpdated)
		//delete
		err = m.TimerDelete(timer.ID)
		assert.Nil(t, err)
		err = m.TimerDelete(timer.ID)
		assert.NotNil(t, err)
		//read
		timerRead, err = m.TimerRead(timer.ID)
		assert.NotNil(t, err)
		assert.Nil(t, timerRead)
	}
}

func TestTimersRead(m interface {
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//TODO: create test
	}
}

func TestTimerLogic(m interface {
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create timer
		comment := randomString(25)
		timer, err := m.TimerCreate(data.TimerPartial{
			Comment: &comment,
			// EmployeeID: &employee.ID,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, timer.ID)
		assert.Equal(t, timer.Comment, comment)
		// assert.Equal(t, timer.EmployeeID, employee.ID)
		//start
		timerStarted, err := m.TimerStart(timer.ID)
		assert.Nil(t, err)
		assert.Greater(t, timerStarted.Start, int64(0))
		assert.Zero(t, timerStarted.Finish)
		time.Sleep(time.Second)
		//read
		timerRead, err := m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, timerRead.ElapsedTime, int64(1))
		//stop
		timerStopped, err := m.TimerStop(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timerStarted.Start, timerStopped.Start)
		//read
		timerRead, err = m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timerStopped, timerRead)
	}
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
