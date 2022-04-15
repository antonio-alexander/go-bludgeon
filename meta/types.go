package meta

import (
	"strings"

	"github.com/antonio-alexander/go-bludgeon/data"
)

//error constants
const (
	ErrTimerNotFoundf     string = "timer with id, \"%s\" not found"
	ErrEmployeeNotFoundf  string = "employee with id, \"%s\" not found"
	ErrTimeSliceNotFoundf string = "timeSlice with id, \"%s\" not found"
)

//SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Timers     map[string]data.Timer     `json:"timers"`
	TimeSlices map[string]data.TimeSlice `json:"time_slices"`
	Employees  map[string]data.Employee  `json:"employees"`
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

type Owner interface {
	Shutdown()
}

type Timer interface {
	//TimerCreate can be used to create a timer, although
	// all fields are available, the only fields that will
	// actually be set are: employee_id and comment
	TimerCreate(timer data.TimerPartial) (*data.Timer, error)

	//TimerRead can be used to read the current value of a given
	// timer, values such as start/finish and elapsed time are
	// "calculated" values rather than values that can be set
	TimerRead(id string) (*data.Timer, error)

	//TimerStart can be used to start a given timer or do nothing
	// if the timer is already started
	TimerStart(id string) (*data.Timer, error)

	//TimerStop can be used to stop a given timer or do nothing
	// if the timer is not started
	TimerStop(id string) (*data.Timer, error)

	//TimerUpdate can be used to update values a given timer
	// not associated with timer operations, values such as:
	// comment, archived and completed
	TimerUpdate(id string, timer data.TimerPartial) (*data.Timer, error)

	//TimerDelete can be used to delete a timer if it exists
	TimerDelete(id string) error

	//TimersRead can be used to read one or more timers depending
	// on search values provided
	TimersRead(search data.TimerSearch) ([]*data.Timer, error)
}

type TimeSlice interface {
	//TimeSliceCreate can be used to create a single time
	// slice
	TimeSliceCreate(t data.TimeSlicePartial) (*data.TimeSlice, error)

	//TimeSliceRead can be used to read an existing time slice
	TimeSliceRead(id string) (*data.TimeSlice, error)

	//TimeSliceUpdate can be used to update an existing time slice
	TimeSliceUpdate(id string, t data.TimeSlicePartial) (*data.TimeSlice, error)

	//TimeSliceDelete can be used to delete an existing time slice
	TimeSliceDelete(id string) error

	//TimeSlicesRead can be used to read zero or more time slices depending on the
	// search criteria
	TimeSlicesRead(search data.TimeSliceSearch) ([]*data.TimeSlice, error)
}

type Employee interface {
	//EmployeeCreate
	EmployeeCreate(e data.EmployeePartial) (*data.Employee, error)

	//EmployeeRead
	EmployeeRead(id string) (*data.Employee, error)

	//EmployeeUpdate
	EmployeeUpdate(id string, e data.EmployeePartial) (*data.Employee, error)

	//EmployeeDelete
	EmployeeDelete(id string) error

	//EmployeesRead
	EmployeesRead(search data.EmployeeSearch) ([]*data.Employee, error)
}