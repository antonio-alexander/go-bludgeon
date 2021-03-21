package bludgeon

//--------------------------------------------------------------------------------------------
// types.go will contain basic types used by the package in general these will be the type sued
// by other functions, it'll include both their types and their methods, the types created here
// can be used elsewhere, unlike other types
//--------------------------------------------------------------------------------------------

import (
	"fmt"
	"strings"
	"time"
)

//common constants
const (
	fmtTimeLong              = "Jan 2, 2006 at 3:04pm (MST)"
	fmtTimeShort             = "2006-Jan-02"
	DefaultFolder            = ".bludgeon"
	DefaultConfigurationFile = "config/bludgeon_config.json"
	DefaultCacheFile         = "data/bludgeon_cache.json"
)

//error constants
const (
	ErrBadProjectID        string = "ProjectID is invalid or missing"
	ErrBadEmployeeID       string = "EmployeeID is invalid or missing"
	ErrBadTaskID           string = "TaskID is invalid or missing"
	ErrBadClientID         string = "ClientID is invalid or missing"
	ErrBadTimerID          string = "TimerID is invalid or missing"
	ErrBadUnitID           string = "UnitID is invalid or missing"
	ErrBadEmployeeIDTaskID string = "EmployeeID and/or TaskID is invalid or mising"
	ErrBadClientIDUnitID   string = "ClientID and/or UnitID is invalid or missing"
	ErrTimerNotFoundf      string = "Timer with id, \"%s\" not found"
	ErrTimeSliceNotFoundf  string = "TimeSlice with id, \"%s\" not found"
	ErrTimerIsArchivedf    string = "Timer with id, \"%s\", is archived"
	ErrNoActiveTimeSlicef  string = "Timer with id, \"%s\", has no active slice"
)

//header constants
const (
	HeaderTimer   string = "TimerID\tTaskID\tEmployeeID\tStart\tFinish\tTimezone\tNote\n"
	HeaderTask    string = "TaskID\tProjectID\tState\tBudget\tDescription\n"
	HeaderProject string = "Project ID\tUnitID\tClientID\tDescription\n"
)

//TimeSlice is the basic unit of "time", the idea is that a task may span over multiple slices that
// aren't necessarily contiguous, but can be added together to get an altogether time.  This should
// reduce the overall error when you pause and restart timers (not time slices) from multiple locations
// time slices can be deleted/archived, but not "edited"
type TimeSlice struct {
	UUID        string `json:"UUID"`        //a unique identifier for the time slice
	TimerUUID   string `json:"TimerUUID"`   //the unique identifier referencing the timer it belongs to
	Start       int64  `json:"Start"`       //the start time in unix nano
	Finish      int64  `json:"Finish"`      //the finish time in unix nano
	ElapsedTime int64  `json:"ElapsedTime"` //the elapsed time in nanoeconds
	Archived    bool   `json:"Archived"`    //whether or not the slice is archived
	//REVIEW: if in the event a timer is "modified" post submission, there will need to be
	// someway to disassociate the time slice from the timer since the elapsed time would be
	// changed
}

//Timer is a high-level object that describes a single unit of time for a given "task". A timer may
// be started and paused many times, but only submitted once. Although there's an obvious desire to
// be able to edit a timer after its submitted, but before it's billed/invoiced. Prior to submission
// elapsed time should "always" be the sum of all the associated active timers and won't necessarily
// be the difference between finish and start (but may be if edited post submission).
type Timer struct {
	UUID            string `json:"UUID"`            //unique id to be shared across network
	ActiveSliceUUID string `json:"ActiveSliceUUID"` //unique id to be shared across network
	Comment         string `json:"Comment"`         //a comment describing the timer
	Start           int64  `json:"Start"`           //the start time in unix nano
	Finish          int64  `json:"Finish"`          //the finish time in unix nano
	ElapsedTime     int64  `json:"ElapsedTime"`     //how much time has elapsed
	Completed       bool   `json:"Completed"`       //this is set to true once the timer has been submitted
	Archived        bool   `json:"Archived"`        //whether or not the timer is archived
	Billed          bool   `json:"Billed"`          //this is set once it has been billed so it "can't" be modified
	EmployeeID      int64  `json:"EmployeeID"`
}

func (t Timer) String() string {
	return fmt.Sprintf(" \"ID\": %s\n \"Active Slice UUID\": %s\n \"Start\": %s\n \"Finish\": %s\n \"Elapsed Time\": %v\n \"Completed\": %t\n \"Comment\": %s",
		t.UUID, t.ActiveSliceUUID, time.Unix(0, t.Start).Format(fmtTimeLong), time.Unix(0, t.Finish).Format(fmtTimeLong), time.Duration(t.ElapsedTime)*time.Nanosecond, t.Completed, t.Comment)
}

//Token provides a way to store token information
type Token struct {
	Token string
	Time  int64
}

type CacheData struct {
	//REVIEW: this will probably just use a signal to send
	// a callback with a type and some common identiier, this
	// will be absolutely useless unless the client stays in
	// memory (e.g. like REST and not CLI)
	//command
	//interface
}

//TaskState provides a type to define the state of a task
type TaskState int

//task states
const (
	TaskStateNull        TaskState = iota
	TaskStateUnbilled    TaskState = iota //Unbilled is the default state, time has been recorded
	TaskStateInvoiced    TaskState = iota //Invoiced is when the task has been invoiced
	TaskStatePaid        TaskState = iota //Paid is when the task has been paid
	TaskStateNonBillable TaskState = iota //NonBillable is when the task is non-billable
	TaskStateInvalid     TaskState = iota
)

func (t TaskState) String() string {
	switch t {
	case TaskStateUnbilled:
		return "unbilled"
	case TaskStateInvoiced:
		return "invoiced"
	case TaskStatePaid:
		return "paid"
	case TaskStateNonBillable:
		return "billable"
	default:
		return "invalid"
	}
}

type Task struct {
	ID          int64     `json:"TaskID,omit_empty"`
	ProjectID   int64     `json:"ProjectID,omit_empty"`
	Description string    `json:"Description,omit_empty"`
	State       TaskState `json:"State,omit_empty"`
	Budget      int64     `json:"Budget,omit_empty"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%d\t%d\t%s\t%d\t%s\n",
		t.ID, t.ProjectID, t.State, t.Budget, t.Description)
}

type Project struct {
	ID          int64  `json:"ProjectID,omit_empty"`
	ClientID    int64  `json:"ClientId,omit_empty"`
	Description string `json:"Descrition,omit_empty"`
}

func (p *Project) String() string {
	return fmt.Sprintf("%d\t%d\t%s\n",
		p.ID, p.ClientID, p.Description)

}

type Employee struct {
	ID        int64  `json:"EmployeeID,omit_empty"`
	FirstName string `json:"FirstName,omit_empty"`
	LastName  string `json:"LastName,omit_empty"`
}

type Client struct {
	ID   int64  `json:"ClientID,omit_empty"`
	Name string `json:"Name,omit_empty"`
}

type Options struct {
	EmployeeID int64  `json:"EmployeeID,omit_empty"`
	ClientID   int64  `json:"ClientID,omit_empty"`
	TimerID    int64  `json:"TimerID,omit_empty"`
	ProjectID  int64  `json:"ProjectID,omit_empty"`
	Token      string `json:"TokenID,omit_empty"`
}

type Logger interface {
	//Println
	Println(v ...interface{})

	//Printf
	Printf(format string, v ...interface{})

	//Print
	Print(v ...interface{})

	//Error
	Error(err error)

	//Errorf
	Errorf(format string, v ...interface{})
}

type FunctionalOwner interface {
	//Initialize
	Initialize(config interface{}) (err error)

	//Shutdown
	Shutdown() (err error)
}

type FunctionalManage interface {
	//
	Stop() (err error)
}

//FunctionalTimer
type FunctionalTimer interface {
	//TimerCreate
	TimerCreate() (timer Timer, err error)

	//TimerRead
	TimerRead(id string) (timer Timer, err error)

	//TimerUpdate
	TimerUpdate(timerIn Timer) (timerOut Timer, err error)

	//TimerDelete
	TimerDelete(id string) (err error)

	//TimerStart
	TimerStart(timerID string, startTime time.Time) (timer Timer, err error)

	//TimerPause
	TimerPause(timerID string, pauseTime time.Time) (timer Timer, err error)

	//TimerSubmit
	TimerSubmit(timerID string, finishTime time.Time) (timer Timer, err error)
}

//FunctionalTimeSlice
type FunctionalTimeSlice interface {
	//TimeSliceRead
	TimeSliceRead(id string) (timeSlice TimeSlice, err error)
}

func AtoMetaType(s string) MetaType {
	switch strings.ToLower(s) {
	case "json":
		return MetaTypeJSON
	case "mysql":
		return MetaTypeMySQL
	default:
		return MetaTypeInvalid
	}
}

type RemoteType string

const (
	RemoteTypeInvalid RemoteType = "invalid"
	RemoteTypeRest    RemoteType = "rest"
)

func AtoRemoteType(s string) RemoteType {
	switch strings.ToLower(s) {
	case "rest":
		return RemoteTypeRest
	default:
		return RemoteTypeInvalid
	}
}

type ObjectType string

const (
	ObjectTypeInvalid   ObjectType = "invalid"
	ObjectTypeTimer     ObjectType = "timer"
	ObjectTypeTimeSlice ObjectType = "timeslice"
)

func AtoObjectType(s string) ObjectType {
	switch strings.ToLower(s) {
	case "timer":
		return ObjectTypeTimer
	case "timeslice":
		return ObjectTypeTimeSlice
	default:
		return ObjectTypeInvalid
	}
}

type MetaType string

//task states
const (
	MetaTypeInvalid MetaType = "invalid"
	MetaTypeJSON    MetaType = "json"
	MetaTypeMySQL   MetaType = "mysql"
)

func (m MetaType) String() string {
	switch m {
	case MetaTypeJSON:
		return "json"
	case MetaTypeMySQL:
		return "mysql"
	default:
		return "invalid"
	}
}

type MetaOwner interface {
	//
	Initialize(config interface{}) (err error)

	//
	Shutdown() (err error)
}

type MetaSerialize interface {
	//Serialize will attempt to commit current data
	Serialize() (err error)

	//Deserialize will attempt to read current data in-memory
	DeSerialize() (err error)
}

//MetaTimer
type MetaTimer interface {
	//MetaTimerWrite
	TimerWrite(timerID string, timer Timer) (err error)

	//MetaTimerDelete
	TimerDelete(timerID string) (err error)

	//MetaTimerRead
	TimerRead(timerID string) (timer Timer, err error)
}

//MetaTimeSlice
type MetaTimeSlice interface {
	//MetaTimerWrite
	TimeSliceWrite(timeSliceID string, timeSlice TimeSlice) (err error)

	//MetaTimerDelete
	TimeSliceDelete(timeSliceID string) (err error)

	//MetaTimerRead
	TimeSliceRead(timeSliceID string) (timeSlice TimeSlice, err error)
}
