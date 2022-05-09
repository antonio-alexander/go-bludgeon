package metamysql

import (
	"database/sql"
)

//common constants
const (
	DatabaseIsolation = sql.LevelSerializable
	LogAlias          = "Database"
)

//error constants
const (
	ErrEmployeeNotFoundf  string = "employee with id, \"%s\", not found locally"
	ErrTimerNotFoundf     string = "timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "timeSlice with id, \"%s\", not found locally"
	ErrTimerIsArchivedf   string = "timer with id, \"%s\", is archived"
	ErrNotImplemented     string = "not implemented"
	ErrDatabaseNil        string = "internal database pointer is nil"
	ErrDatabaseNotNil     string = "internal database pointer is not nil, reconnect or close to connect"
	ErrUpdateFailed       string = "update failed, no rows affected"
	ErrDeleteFailed       string = "delete failed, id not found"
	ErrQueryFailed        string = "query: \"%s\", failed"
	ErrStarted            string = "already started"
	ErrNotStarted         string = "not started"
)

//query constants
const (
	tableEmployees    string = "employees"
	tableTimers       string = "timers"
	tableTimeSlices   string = "time_slices"
	tableTimersV1     string = "timers_v1"
	tableTimeSlicesV1 string = "time_slices_v1"
	tableEmployeesV1  string = "employees_v1"
)

type Owner interface {
	Initialize(config *Configuration) (err error)
}
