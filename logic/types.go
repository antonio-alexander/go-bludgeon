package logic

import (
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
)

type Functional interface {
	Start() (err error)
	Stop() (err error)
}

type Logic interface {
	TimerCreate() (timer common.Timer, err error)
	TimerRead(id string) (timer common.Timer, err error)
	TimerUpdate(t common.Timer) (timer common.Timer, err error)
	TimerDelete(timerID string) (err error)
	TimerStart(id string, startTime time.Time) (timer common.Timer, err error)
	TimerPause(timerID string, pauseTime time.Time) (timer common.Timer, err error)
	TimerSubmit(timerID string, submitTime time.Time) (timer common.Timer, err error)
	TimeSliceRead(timeSliceID string) (timeSlice common.TimeSlice, err error)
}
