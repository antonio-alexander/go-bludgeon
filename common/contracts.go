package common

//route constants
const (
	//
	RouteStop = "/api/stop"
	//
	RouteAdmin       = "/api/admin"
	RouteAdminConfig = RouteAdmin + "/config"
	RouteAdminStatus = RouteAdmin + "/status"
	//
	RouteTimer       = "/api/timer"
	RouteTimerCreate = RouteTimer + "/create"
	RouteTimerRead   = RouteTimer + "/read"
	RouteTimerUpdate = RouteTimer + "/update"
	RouteTimerDelete = RouteTimer + "/delete"
	RouteTimersRead  = RouteTimer + "/read"
	RouteTimerStart  = RouteTimer + "/start"
	RouteTimerPause  = RouteTimer + "/pause"
	RouteTimerSubmit = RouteTimer + "/submit"
	//
	RouteTimeSlice     = "/api/timeslice"
	RouteTimeSliceRead = RouteTimeSlice + "/read"
	//
	RouteTaskCreate = "/api/task/create"
	RouteTaskRead   = "/api/task/read"
	RouteTaskUpdate = "/api/task/update"
	RouteTaskDelete = "/api/task/delete"
	RouteTasksRead  = "/api/tasks/read"
	//
	RouteProjectCreate = "/api/project/create"
	RouteProjectRead   = "/api/project/read"
	RouteProjectUpdate = "/api/project/update"
	RouteProjectDelete = "/api/project/delete"
	RouteProjectsRead  = "/api/projects/read"
	//
	RouteEmployeeCreate = "/api/employee/create"
	RouteEmployeeRead   = "/api/employee/read"
	RouteEmployeeUpdate = "/api/employee/update"
	RouteEmployeeDelete = "/api/employee/delete"
	RouteEmployeesRead  = "/api/employees/read"
	//
	RouteClientCreate = "/api/client/create"
	RouteClientRead   = "/api/client/read"
	RouteClientUpdate = "/api/client/update"
	RouteClientDelete = "/api/client/delete"
	RouteClientsRead  = "/api/clients/read"
	//
	RouteTokenAcquire = "/api/token/acquire"
	RouteTokenRelease = "/api/token/release"
	RouteTokenVerify  = "/api/token/verify"
	//
	RouteDebug        = "/api/debug"
	RouteDebugEnable  = RouteDebug + "enable"
	RouteDebugDisable = RouteDebug + "disable"
	//
	RouteSync               = "/api/sync"
	RouteSyncTimersProject  = RouteSync + "/timers_project"
	RouteSyncTimersEmployee = RouteSync + "/timers_employee"
	RouteSyncTasksProject   = RouteSync + "/tasks_project"
	RouteSyncTasksEmployee  = RouteSync + "/tasks_employee"
)

//Contract
type Contract struct {
	ID         string    `json:"ID,omitempty"`
	StartTime  int64     `json:"StartTime,omitempty"`
	PauseTime  int64     `json:"PauseTime,omitempty"`
	FinishTime int64     `json:"FinishTime,omitempty"`
	Timer      Timer     `json:"Timer,omitempty"`
	TimeSlice  TimeSlice `json:"TimeSlice,omitempty"`
}
