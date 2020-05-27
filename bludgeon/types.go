package bludgeon

//--------------------------------------------------------------------------------------------
// types.go will contain basic types used by the package in general these will be the type sued
// by other functions, it'll include both their types and their methods, the types created here
// can be used elsewhere, unlike other types
//--------------------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"strings"
	"time"
)

//common constants
const (
	fmtTimeLong  = "Jan 2, 2006 at 3:04pm (MST)"
	fmtTimeShort = "2006-Jan-02"
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
)

//route constants
const (
	RouteTimerCreate        = "/api/timer/create" //RouteTimerCreate is a route used to create a timer
	RouteTimerRead          = "/api/timer/read"
	RouteTimerUpdate        = "/api/timer/update"
	RouteTimerDelete        = "/api/timer/delete"
	RouteTimersRead         = "/api/timers/read"
	RouteTaskCreate         = "/api/task/create"
	RouteTaskRead           = "/api/task/read"
	RouteTaskUpdate         = "/api/task/update"
	RouteTaskDelete         = "/api/task/delete"
	RouteTasksRead          = "/api/tasks/read"
	RouteProjectCreate      = "/api/project/create"
	RouteProjectRead        = "/api/project/read"
	RouteProjectUpdate      = "/api/project/update"
	RouteProjectDelete      = "/api/project/delete"
	RouteProjectsRead       = "/api/projects/read"
	RouteEmployeeCreate     = "/api/employee/create"
	RouteEmployeeRead       = "/api/employee/read"
	RouteEmployeeUpdate     = "/api/employee/update"
	RouteEmployeeDelete     = "/api/employee/delete"
	RouteEmployeesRead      = "/api/employees/read"
	RouteClientCreate       = "/api/client/create"
	RouteClientRead         = "/api/client/read"
	RouteClientUpdate       = "/api/client/update"
	RouteClientDelete       = "/api/client/delete"
	RouteClientsRead        = "/api/clients/read"
	RouteServerStop         = "/api/server/stop"
	RouteTokenAcquire       = "/api/token/acquire"
	RoutetokenRelease       = "/api/token/release"
	RouteTokenVerify        = "/api/token/verify"
	RouteDebugEnable        = "/api/debug/enable"
	RouteDebugDisable       = "/api/debug/disable"
	RouteServerStart        = "/api/server/start"
	RouteAdminConfig        = "/api/admin/config"
	RouteAdminStatus        = "/api/admin/status"
	RouteSyncTimersProject  = "/api/sync/timers_project"
	RouteSyncTimersEmployee = "/api/sync/timers_employee"
	RouteSyncTasksProject   = "/api/sync/tasks_project"
	RouteSyncTasksEmployee  = "/api/sync/tasks_employee"
)

//header constants
const (
	HeaderTimer   string = "TimerID\tTaskID\tEmployeeID\tStart\tFinish\tTimezone\tNote\n"
	HeaderTask    string = "TaskID\tProjectID\tState\tBudget\tDescription\n"
	HeaderProject string = "Project ID\tUnitID\tClientID\tDescription\n"
)

//Token provides a way to store token information
type Token struct {
	Token string
	Time  int64
}

//TaskState provides a type to define the state of a task
type TaskState int

//task states
const (
	TaskStateUnbilled    TaskState = 0 //Unbilled is the default state, time has been recorded
	TaskStateInvoiced    TaskState = 1 //Invoiced is when the task has been invoiced
	TaskStatePaid        TaskState = 2 //Paid is when the task has been paid
	TaskStateNonBillable TaskState = 3 //NonBillable is when the task is non-billable
)

//TimeSlice is the basic unit of "time", the idea is that a task may span over multiple slices that
// aren't necessarily contiguous, but can be added together to get an altogether time.  This should
// reduce the overall error when you pause and restart timers (not time slices) from multiple locations
// time slices can be deleted/archived, but not "edited"
type TimeSlice struct {
	ID          int64
	TimerID     int64
	UUID        string `json:"UUID"`
	TimerUUID   string `json:"TimerUUID"`
	Start       int64  `json:"Start"`
	Finish      int64  `json:"Finish"`
	ElapsedTime int64  `json:"ElapsedTime"`
	//not currently used
	Archived bool `json:"Archived"`
}

//Timer
type Timer struct {
	ID              int64  //primary key for database, not shared across network
	ActiveSliceID   int64  //foreign key for active slice, not shared across network
	UUID            string `json:"UUID"`            //unique id to be shared across network
	ActiveSliceUUID string `json:"ActiveSliceUUID"` //unique id to be shared across network
	Start           int64  `json:"Start"`
	Finish          int64  `json:"Finish"`
	ElapsedTime     int64  `json:"ElapsedTime"`
	Completed       bool   `json:"Completed"` //this is set to true once the timer has been submitted
	//currently unused
	Billed     bool   `json:"Billed"` //this is set once it has been billed so it "can't" be modified
	EmployeeID int64  `json:"EmployeeID"`
	Comment    string `json:"Comment"`
	Archived   bool   `json:"Archived"`
}

func (t Timer) String() string {
	return fmt.Sprintf(" \"UUID\": %s\n \"Active Slice UUID\": %s\n \"Start\": %s\n \"Finish\": %s\n \"Elapsed Time\": %v\n \"Completed\": %t\n",
		t.UUID, t.ActiveSliceUUID, time.Unix(0, t.Start).Format(fmtTimeLong), time.Unix(0, t.Finish).Format(fmtTimeLong), time.Duration(t.ElapsedTime)/time.Nanosecond, t.Completed)
}

type Task struct {
	ID          int64     `json:"TaskID,omit_empty"`
	ProjectID   int64     `json:"ProjectID,omit_empty"`
	Description string    `json:"Description,omit_empty"`
	State       TaskState `json:"State,omit_empty"`
	Budget      int64     `json:"Budget,omit_empty"`
}

func (t *Task) Print(w io.Writer) error {
	var state string

	switch t.State {
	case TaskStateUnbilled:
		state = "Unbilled"
	case TaskStateInvoiced:
		state = "Invoiced"
	case TaskStatePaid:
		state = "Paid"
	case TaskStateNonBillable:
		state = "NonBillable"
	}

	_, err := fmt.Fprintf(w, "%d\t%d\t%s\t%d\t%s\n",
		t.ID, t.ProjectID, state, t.Budget, t.Description)

	return err
}

type Project struct {
	ID          int64  `json:"ProjectID,omit_empty"`
	ClientID    int64  `json:"ClientId,omit_empty"`
	Description string `json:"Descrition,omit_empty"`
}

func (p *Project) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d\t%d\t%s\n",
		p.ID, p.ClientID, p.Description)

	return err
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
	CommandServerNull    CommandServer = iota
	CommandServerInvalid CommandServer = iota
)
