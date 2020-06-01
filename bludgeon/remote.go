package bludgeon

import "time"

//REVIEW: would it make more sense to call this api instead of
// remote??

//Remote
type Remote interface {
	RemoteTimer
	RemoteTimeSlice
}

//RemoteTimer
type RemoteTimer interface {
	//TimerCreate
	TimerCreate() (timer Timer, err error)

	//TimerRead
	TimerRead(id string) (timer Timer, err error)

	//TimerUpdate
	TimerUpdate(timer Timer) (err error)

	//TimerDelete
	TimerDelete(id string) (err error)

	//TimerStart
	TimerStart(timerID string, startTime time.Time) (err error)

	//TimerPause
	TimerPause(timerID string, pauseTime time.Time) (err error)

	//TimerSubmit
	TimerSubmit(timerID string, finishTime time.Time) (err error)
}

//RemoteTimeSlice
type RemoteTimeSlice interface {
	//TimeSliceCreate
	TimeSliceCreate(id string) (timeSlice TimeSlice, err error)

	//TimeSliceRead
	TimeSliceRead(id string) (timeSlice TimeSlice, err error)

	//TimeSliceUpdate
	TimeSliceUpdate(timeSlice TimeSlice) (err error)

	//TimeSliceDelete
	TimeSliceDelete(id string) (err error)
}
