package metajson

import (
	common "github.com/antonio-alexander/go-bludgeon/common"
)

//error constants
const (
	ErrTimerNotFoundf     string = "timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "timeSlice with id, \"%s\", not found locally"
)

//SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Timers     map[string]common.Timer     `json:"Timers"`
	TimeSlices map[string]common.TimeSlice `json:"TimeSlices"`
}
