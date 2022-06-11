package memory

import "github.com/pkg/errors"

const (
	TimeSliceStartZero       string = "time slice start is zero"
	TimeSliceFinishLessStart string = "time slice finish less than or equal to start"
	UnfinishedTimeSlice      string = "timer id has unfnished time slice"
	TimeSliceOverlap         string = "time slice overlaps with another time slice"
	TimeSliceNoTimerID       string = "time slice has no timer id"
)

var (
	ErrTimeSliceStartZero       = errors.New(TimeSliceStartZero)
	ErrTimeSliceFinishLessStart = errors.New(TimeSliceFinishLessStart)
	ErrUnfinishedTimeSlice      = errors.New(UnfinishedTimeSlice)
	ErrTimeSliceOverlap         = errors.New(TimeSliceOverlap)
	ErrTimeSliceNoTimerID       = errors.New(TimeSliceNoTimerID)
)
