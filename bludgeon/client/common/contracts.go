package bludgeonclient

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

//route constants
const (
	//
	RouteClient      = "/api/client"
	RouteClientStop  = RouteClient + "/stop"
	RouteClientStart = RouteClient + "/start"
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
	// //
	// RouteTaskCreate = "/api/task/create"
	// RouteTaskRead   = "/api/task/read"
	// RouteTaskUpdate = "/api/task/update"
	// RouteTaskDelete = "/api/task/delete"
	// RouteTasksRead  = "/api/tasks/read"
	// //
	// RouteProjectCreate = "/api/project/create"
	// RouteProjectRead   = "/api/project/read"
	// RouteProjectUpdate = "/api/project/update"
	// RouteProjectDelete = "/api/project/delete"
	// RouteProjectsRead  = "/api/projects/read"
	// //
	// RouteEmployeeCreate = "/api/employee/create"
	// RouteEmployeeRead   = "/api/employee/read"
	// RouteEmployeeUpdate = "/api/employee/update"
	// RouteEmployeeDelete = "/api/employee/delete"
	// RouteEmployeesRead  = "/api/employees/read"
	// //
	// RouteClientCreate = "/api/client/create"
	// RouteClientRead   = "/api/client/read"
	// RouteClientUpdate = "/api/client/update"
	// RouteClientDelete = "/api/client/delete"
	// RouteClientsRead  = "/api/clients/read"
	// //
	// RouteTokenAcquire = "/api/token/acquire"
	// RouteTokenRelease = "/api/token/release"
	// RouteTokenVerify  = "/api/token/verify"
	// //
	// RouteDebugEnable  = "/api/debug/enable"
	// RouteDebugDisable = "/api/debug/disable"
	// //
	// RouteSyncTimersProject  = "/api/sync/timers_project"
	// RouteSyncTimersEmployee = "/api/sync/timers_employee"
	// RouteSyncTasksProject   = "/api/sync/tasks_project"
	// RouteSyncTasksEmployee  = "/api/sync/tasks_employee"
)

type ContractClientIn struct {
	ID         string             `json:"ID,omit_empty"`
	StartTime  int64              `json:"StartTime,omit_empty"`
	PauseTime  int64              `json:"PauseTime,omit_empty"`
	FinishTime int64              `json:"FinishTime,omit_empty"`
	Timer      bludgeon.Timer     `json:"Timer,omit_empty"`
	TimeSlice  bludgeon.TimeSlice `json:"TimeSlice,omit_empty"`
}
