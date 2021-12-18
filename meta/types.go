package meta

import (
	"strings"

	"github.com/antonio-alexander/go-bludgeon/data"
)

//error constants
const (
	ErrTimerNotFoundf     string = "timer with id, \"%s\" not found"
	ErrTimeSliceNotFoundf string = "timeSlice with id, \"%s\" not found"
)

//SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Timers     map[string]data.Timer     `json:"Timers"`
	TimeSlices map[string]data.TimeSlice `json:"TimeSlices"`
}

type Type string

//task states
const (
	TypeInvalid Type = "invalid"
	TypeFile    Type = "file"
	TypeMySQL   Type = "mysql"
)

func (m Type) String() string {
	switch m {
	case TypeFile:
		return "json"
	case TypeMySQL:
		return "mysql"
	default:
		return "invalid"
	}
}

func AtoType(s string) Type {
	switch strings.ToLower(s) {
	case "json":
		return TypeFile
	case "mysql":
		return TypeMySQL
	default:
		return TypeInvalid
	}
}

type Serializer interface {
	SerializedDataRead() SerializedData
	SerializedDataWrite(SerializedData)
}

type Owner interface {
	Shutdown() (err error)
}

//MetaTimer
type Timer interface {
	//MetaTimerWrite
	TimerWrite(timerID string, timer data.Timer) (err error)

	//MetaTimerDelete
	TimerDelete(timerID string) (err error)

	//MetaTimerRead
	TimerRead(timerID string) (timer data.Timer, err error)
}

//MetaTimeSlice
type TimeSlice interface {
	//MetaTimerWrite
	TimeSliceWrite(timeSliceID string, timeSlice data.TimeSlice) (err error)

	//MetaTimerDelete
	TimeSliceDelete(timeSliceID string) (err error)

	//MetaTimerRead
	TimeSliceRead(timeSliceID string) (timeSlice data.TimeSlice, err error)
}
