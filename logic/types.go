package logic

import (
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
)

type Functional interface {
	Start() (err error)
	Stop() (err error)
}

type Logic interface {
	TimerCreate() (timer data.Timer, err error)
	TimerRead(id string) (timer data.Timer, err error)
	TimerUpdate(t data.Timer) (timer data.Timer, err error)
	TimerDelete(timerID string) (err error)
	TimerStart(id string, startTime time.Time) (timer data.Timer, err error)
	TimerPause(timerID string, pauseTime time.Time) (timer data.Timer, err error)
	TimerSubmit(timerID string, submitTime time.Time) (timer data.Timer, err error)
	TimeSliceRead(timeSliceID string) (timeSlice data.TimeSlice, err error)
}
