package logic

import (
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
)

//Logic defines functions that describe the business logic
// of the timers micro service
type Logic interface {
	meta.TimeSlice

	//KIM: most of these functions below are overriden (and masked)
	// by the underlying lgoic pointer
	// meta.Timer

	//TimerCreate can be used to create a timer, although
	// all fields are available, the only fields that will
	// actually be set are: timer_id and comment
	TimerCreate(timer data.TimerPartial) (*data.Timer, error)

	//TimerRead can be used to read the current value of a given
	// timer, values such as start/finish and elapsed time are
	// "calculated" values rather than values that can be set
	TimerRead(id string) (*data.Timer, error)

	//TimersRead can be used to read one or more timers depending
	// on search values provided
	TimersRead(search data.TimerSearch) ([]*data.Timer, error)
	//TimerStart can be used to start a given timer or do nothing
	// if the timer is already started
	TimerStart(id string) (*data.Timer, error)

	//TimerStop can be used to stop a given timer or do nothing
	// if the timer is not started
	TimerStop(id string) (*data.Timer, error)

	//TimerDelete can be used to delete a timer if it exists
	TimerDelete(id string) error

	//TimerSubmit can be used to stop a timer and set completed to true
	TimerSubmit(timerID string, finishTime *time.Time) (*data.Timer, error)

	//TimerUpdateCommnet will only update the comment for timer with
	// the provided id
	TimerUpdateComment(id, comment string) (*data.Timer, error)

	//TimerArchive will only update the archive for timer with
	// the provided id
	TimerArchive(id string, archive bool) (*data.Timer, error)
}
