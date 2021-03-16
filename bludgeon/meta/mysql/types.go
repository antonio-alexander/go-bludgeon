package bludgeonmetamysql

import (
	"database/sql"
)

//common constants
const (
	//DatabaseIsolation provides a constant that shows the database isolation level
	DatabaseIsolation = sql.LevelSerializable
	//LogAlias provides the alias when data is logged
	LogAlias = "Database"
)

//error constants
const (
	ErrTimerNotFoundf     string = "Timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "TimeSlice with id, \"%s\", not found locally"
	ErrTimerIsArchivedf   string = "Timer with id, \"%s\", is archived"
	ErrNotImplemented     string = "Not implemented"
	//ErrDatabaseNil provides a string to return as an error if the database pointer is nil
	ErrDatabaseNil string = "Internal database pointer is nil"
	//ErrDatabaseNotNil provides a string to return as an error if you attempt to connect to an already initialized database
	ErrDatabaseNotNil string = "Internal database pointer is not nil, reconnect or close to connect"
	//ErrUpdateFailed provides a string to return as an error if an update fails and the result returns 0 rows affected
	ErrUpdateFailed string = "Update failed, no rows affected"
	//ErrDeleteFailed provides a string to return as an error if a delete fails and the result returns 0 rows affected
	ErrDeleteFailed string = "Delete failed, id not found"
	//ErrQueryFailed provides a string to return as an error in the event a query fails and no other error is returned
	ErrQueryFailed string = "Query: \"%s\", failed"
)

//query constants
const (
	TableTimer       string = "timer"
	TableSlice       string = "slice"
	TableTimerSlice  string = "timer_slice"
	TableActiveSlice string = "timer_slice_active"
	TableProject     string = "project"
	TableClient      string = "client"
	TableEmployee    string = "employee"
)
