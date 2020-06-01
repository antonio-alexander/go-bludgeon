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
	DefaultConfigurationFile = "/config/bludgeon_config.json"
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
	ErrTimerIsArchivedf    string = "Timer with id, \"%s\", is archived"
	ErrNoActiveTimeSlicef  string = "Timer with id, \"%s\", has no active slice"
	ErrTimerNotFoundf      string = "Timer with id, \"%s\" not found"
	ErrTimeSliceNotFoundf  string = "TimeSlice with id, \"%s\" not found"
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
	Start           int64  `json:"Start"`           //the start time in unix nano
	Finish          int64  `json:"Finish"`          //the finish time in unix nano
	ElapsedTime     int64  `json:"ElapsedTime"`     //how much time has elapsed
	Completed       bool   `json:"Completed"`       //this is set to true once the timer has been submitted
	Comment         string `json:"Comment"`         //a comment describing the timer
	Archived        bool   `json:"Archived"`        //whether or not the timer is archived
	//currently unused
	// Billed     bool   `json:"Billed"` //this is set once it has been billed so it "can't" be modified
	// EmployeeID int64  `json:"EmployeeID"`
}

func (t Timer) String() string {
	return fmt.Sprintf(" \"UUID\": %s \"Active Slice UUID\": %s\n \"Start\": %s\n \"Finish\": %s\n \"Elapsed Time\": %v\n \"Completed\": %t\n \"Comment\": %s",
		t.UUID, t.ActiveSliceUUID, time.Unix(0, t.Start).Format(fmtTimeLong), time.Unix(0, t.Finish).Format(fmtTimeLong), time.Duration(t.ElapsedTime)/time.Nanosecond, t.Completed, t.Comment)
}

//Token provides a way to store token information
type Token struct {
	Token string
	Time  int64
}

type CommandClient uint8

//command constants
const (
	CommandClientNull        CommandClient = iota
	CommandClientTimerCreate CommandClient = iota
	CommandClientTimerRead   CommandClient = iota
	CommandClientTimerDelete CommandClient = iota
	CommandClientTimerStart  CommandClient = iota
	CommandClientTimerStop   CommandClient = iota
	CommandClientTimerPause  CommandClient = iota
	CommandClientTimerSubmit CommandClient = iota
	CommandClientInvalid     CommandClient = iota
)

func (c CommandClient) String() string {
	switch c {
	case CommandClientTimerCreate:
		return "timercreate"
	case CommandClientTimerRead:
		return "timerread"
	case CommandClientTimerDelete:
		return "timerdelete"
	case CommandClientTimerStart:
		return "timerstart"
	case CommandClientTimerStop:
		return "timerstop"
	case CommandClientTimerPause:
		return "timerpause"
	case CommandClientTimerSubmit:
		return "timersubmit"
	default:
		return ""
	}
}

func AtoCommandClient(s string) CommandClient {
	switch strings.ToLower(s) {
	case "timercreate":
		return CommandClientTimerCreate
	case "timerread":
		return CommandClientTimerRead
	case "timerdelete":
		return CommandClientTimerDelete
	case "timerstart":
		return CommandClientTimerStart
	case "timerstop":
		return CommandClientTimerStop
	case "timerpause":
		return CommandClientTimerPause
	case "timersubmit":
		return CommandClientTimerSubmit
	default:
		return CommandClientInvalid
	}
}

func AtoObject(s string) interface{} {
	//parse the object type
	switch strings.ToLower(s) {
	case "t", "timer":
		return Timer{}
	default:
		return nil
	}
}

type CommandServer uint8

//command constants
const (
	CommandServerNull CommandServer = iota
	//
	CommandServerStop CommandServer = iota
	//timer
	CommandServerTimerCreate CommandServer = iota
	CommandServerTimerRead   CommandServer = iota
	CommandServerTimerUpdate CommandServer = iota
	CommandServerTimerDelete CommandServer = iota
	CommandServerTimerStart  CommandServer = iota
	CommandServerTimerPause  CommandServer = iota
	CommandServerTimerSubmit CommandServer = iota
	//time slice
	CommandServerTimeSliceCreate CommandServer = iota
	CommandServerTimeSliceRead   CommandServer = iota
	CommandServerTimeSliceUpdate CommandServer = iota
	CommandServerTimeSliceDelete CommandServer = iota
	//token
	CommandServerTokenAcquire CommandServer = iota
	CommandServerTokenRelease CommandServer = iota
	CommandServerTokenVerify  CommandServer = iota
	// CommandServerProjectCreate   CommandServer = iota
	// CommandServerProjectRead     CommandServer = iota
	// CommandServerProjectsRead    CommandServer = iota
	// CommandServerProjectUpdate   CommandServer = iota
	// CommandServerProjectDelete   CommandServer = iota
	// CommandServerEmployeeCreate  CommandServer = iota
	// CommandServerEmployeeRead    CommandServer = iota
	// CommandServerEmployeesRead   CommandServer = iota
	// CommandServerEmployeeUpdate  CommandServer = iota
	// CommandServerEmployeeDelete  CommandServer = iota
	// CommandServerClientCreate    CommandServer = iota
	// CommandServerClientRead      CommandServer = iota
	// CommandServerClientsRead     CommandServer = iota
	// CommandServerClientUpdate    CommandServer = iota
	// CommandServerClientDelete    CommandServer = iota
	CommandServerInvalid CommandServer = iota
)

func (c CommandServer) String() string {
	switch c {
	//TODO: generate stringers
	default:
		return "invalid"
	}
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

type CacheData struct {
	//REVIEW: this will probably just use a signal to send
	// a callback with a type and some common identiier, this
	// will be absolutely useless unless the client stays in
	// memory (e.g. like REST and not CLI)
	//command
	//interface
}
