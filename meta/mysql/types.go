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
	TableTimer            string = "timer"
	TableSlice            string = "slice"
	TableTimerSliceActive string = "timer_slice_active"
	// TableProject          string = "project"
	// TableClient           string = "client"
	// TableEmployee         string = "employee"
)
