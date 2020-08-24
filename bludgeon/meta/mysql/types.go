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
	//
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
	TableTimer            string = "timer"
	QueryTimerCreateTable string = `CREATE TABLE IF NOT EXISTS ` + TableTimer + ` (
		id BIGINT NOT NULL AUTO_INCREMENT,
		uuid TEXT(36),
		activesliceuuid TEXT(36),    
		start BIGINT,
		finish BIGINT,
		elapsedtime BIGINT,
		comment TEXT,
		INDEX(id),
		UNIQUE(uuid(36)),
		PRIMARY KEY (id)
		-- FOREIGN KEY (employeeid)
			-- REFERENCES employee(id)
			-- ON UPDATE CASCADE ON DELETE RESTRICT
	)ENGINE=InnoDB;`
	QueryTimerUpsert string = `INSERT into ` + TableTimer + ` (uuid, activesliceuuid, start, finish, elapsedtime, comment) VALUES(?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY 
		UPDATE uuid=?, activesliceuuid=?, start=?, finish=?, elapsedtime=?, comment=?`
	QueryTimerDeletef string = `DELETE FROM ` + TableTimer + ` where uuid = ?`
	QueryTimerSelectf string = `SELECT uuid, activesliceuuid, start, finish, elapsedtime from ` + TableTimer + ` WHERE uuid = ?`
	//
	TableTimeSlice            string = "timeslice"
	QueryTimeSliceCreateTable string = `CREATE TABLE IF NOT EXISTS ` + TableTimeSlice + ` (
		id BIGINT NOT NULL AUTO_INCREMENT,
		uuid TEXT(36),
		timeruuid TEXT(36),    
		start BIGINT,
		finish BIGINT,
		elapsedtime BIGINT,
		INDEX(id),
		UNIQUE(uuid(36)),
		PRIMARY KEY (id)
		-- FOREIGN KEY (timeruuid(36))
		--     REFERENCES timer(uuid)
		--     ON DELETE CASCADE
	)ENGINE=InnoDB;`
	QueryTimeSliceUpsert string = `INSERT into ` + TableTimeSlice + ` (uuid, timeruuid, start, finish, elapsedtime) VALUES(?, ?, ?, ?, ?)
	ON DUPLICATE KEY
		UPDATE uuid=?, start=?, finish=?, elapsedtime=?`
	QueryTimeSliceDeletef string = `DELETE FROM ` + TableTimeSlice + ` where uuid = ?`
	QueryTimeSliceSelectf string = `SELECT uuid, timeruuid, start, finish, elapsedtime from ` + TableTimeSlice + ` WHERE uuid = ?`
)

// //TableProject is the string defining the name of the project table
// TableProject string = "project"
// //TableClient is the string defining the name of the client table
// TableClient string = "client"
// //TableEmployee is the string defining the name of the employee table
// TableEmployee string = "employee"
