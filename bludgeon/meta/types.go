package bludgeonmeta

import (
	common "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type Meta interface {
	//Serialize will attempt to commit current data
	Serialize() (err error)

	//Deserialize will attempt to read current data in-memory
	DeSerialize() (err error)
}

//MetaTimer
type MetaTimer interface {
	//MetaTimerWrite
	MetaTimerWrite(timerID string, timer common.Timer) (err error)

	//MetaTimerDelete
	MetaTimerDelete(timerID string) (err error)

	//MetaTimerRead
	MetaTimerRead(timerID string) (timer common.Timer, err error)
}

//MetaTimer
type MetaTimeSlice interface {
	//MetaTimerWrite
	MetaTimeSliceWrite(timeSliceID string, timeSlice common.TimeSlice) (err error)

	//MetaTimerDelete
	MetaTimeSliceDelete(timeSliceID string) (err error)

	//MetaTimerRead
	MetaTimeSliceRead(timeSliceID string) (timer common.TimeSlice, err error)
}
