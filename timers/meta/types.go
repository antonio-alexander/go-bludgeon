package meta

import (
	"context"
	"strings"

	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// error constants
const (
	TimerNotFound       string = "timer not found"
	TimerNotUpdated     string = "timer not updated"
	TimerNotCreated     string = "timer not created, email address not provided"
	TimerConflictCreate string = "cannot create timer; email address in use"
	TimerConflictUpdate string = "cannot update timer; email address in use"
	TimeSliceNotFound   string = "time slice not found"
)

// error variables
var (
	ErrTimerNotFound       = errors.NewNotFound(errors.New(TimerNotFound))
	ErrTimerNotUpdated     = errors.NewNotUpdated(errors.New(TimerNotUpdated))
	ErrTimerNotCreated     = errors.NewNotCreated(errors.New(TimerNotCreated))
	ErrTimerConflictCreate = errors.NewConflict(errors.New(TimerConflictCreate))
	ErrTimerConflictUpdate = errors.NewConflict(errors.New(TimerConflictUpdate))
	ErrTimeSliceNotFound   = errors.NewNotFound(errors.New(TimeSliceNotFound))
)

// SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Timers     map[string]data.Timer     `json:"timers"`
	TimeSlices map[string]data.TimeSlice `json:"time_slices"`
}

type Type string

// task states
const (
	TypeInvalid Type = "invalid"
	TypeFile    Type = "file"
	TypeMySQL   Type = "mysql"
)

func (m Type) String() string {
	switch m {
	case TypeFile:
		return "file"
	case TypeMySQL:
		return "mysql"
	default:
		return "invalid"
	}
}

func AtoType(s string) Type {
	switch strings.ToLower(s) {
	case "file":
		return TypeFile
	case "mysql":
		return TypeMySQL
	default:
		return TypeInvalid
	}
}

type Serializer interface {
	//Serialize can be used to convert all available metadata
	// to a single pointer to be used to serialize to bytes
	Serialize() (*SerializedData, error)

	//Deserialize can be used to provide metadata as a single pointer
	// once it's been deserialized from bytes
	Deserialize(data *SerializedData) error
}

// Timer provides an interface that can be used to interact with timers
type Timer interface {
	//TimerCreate can be used to create a timer, although
	// all fields are available, the only fields that will
	// actually be set are: timer_id and comment
	TimerCreate(ctx context.Context, timer data.TimerPartial) (*data.Timer, error)

	//TimerRead can be used to read the current value of a given
	// timer, values such as start/finish and elapsed time are
	// "calculated" values rather than values that can be set
	TimerRead(ctx context.Context, id string) (*data.Timer, error)

	//TimerStart can be used to start a given timer or do nothing
	// if the timer is already started
	TimerStart(ctx context.Context, id string) (*data.Timer, error)

	//TimerStop can be used to stop a given timer or do nothing
	// if the timer is not started
	TimerStop(ctx context.Context, id string) (*data.Timer, error)

	//TimerUpdate can be used to update values a given timer
	// not associated with timer operations, values such as:
	// comment, archived and completed
	TimerUpdate(ctx context.Context, id string, timer data.TimerPartial) (*data.Timer, error)

	//TimerSubmit can be used to stop a timer and set completed to true
	TimerSubmit(ctx context.Context, id string, finishTime int64) (*data.Timer, error)

	//TimerDelete can be used to delete a timer if it exists
	TimerDelete(ctx context.Context, id string) error

	//TimersRead can be used to read one or more timers depending
	// on search values provided
	TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error)
}

// TimeSlice provides an interface that can be used to interact with time slices
type TimeSlice interface {
	//TimeSliceCreate can be used to create a single time
	// slice
	TimeSliceCreate(ctx context.Context, t data.TimeSlicePartial) (*data.TimeSlice, error)

	//TimeSliceRead can be used to read an existing time slice
	TimeSliceRead(ctx context.Context, id string) (*data.TimeSlice, error)

	//TimeSliceUpdate can be used to update an existing time slice
	TimeSliceUpdate(ctx context.Context, id string, t data.TimeSlicePartial) (*data.TimeSlice, error)

	//TimeSliceDelete can be used to delete an existing time slice
	TimeSliceDelete(ctx context.Context, id string) error

	//TimeSlicesRead can be used to read zero or more time slices depending on the
	// search criteria
	TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error)
}
