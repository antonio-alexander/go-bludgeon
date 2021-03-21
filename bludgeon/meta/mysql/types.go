package bludgeonmetamysql

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
	ErrTimerNotFoundf     string = "Timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "TimeSlice with id, \"%s\", not found locally"
	ErrTimerIsArchivedf   string = "Timer with id, \"%s\", is archived"
	ErrNotImplemented     string = "Not implemented"
	ErrDatabaseNil        string = "Internal database pointer is nil"
	ErrDatabaseNotNil     string = "Internal database pointer is not nil, reconnect or close to connect"
	ErrUpdateFailed       string = "Update failed, no rows affected"
	ErrDeleteFailed       string = "Delete failed, id not found"
	ErrQueryFailed        string = "Query: \"%s\", failed"
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
