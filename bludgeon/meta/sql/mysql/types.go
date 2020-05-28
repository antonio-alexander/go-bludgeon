package bludgeonmetamysql

import (
	"database/sql"
	"time"
)

//common constants
const (
	//DatabaseIsolation provides a constant that shows the database isolation level
	DatabaseIsolation = sql.LevelSerializable
	//LogAlias provides the alias when data is logged
	LogAlias = "Database"
	//DefaultDriver provides a default driver to use for configuration when no configuration is found
	DefaultDriver = "mysql"
	//DefaultUsername provides a default username to use for configuration when no configuration is found
	DefaultUsername = "bludgeon"
	//DefaultPassword provides a default pasword to use for configuration when no configuration is found
	DefaultPassword = "bludgeon"
	//DefaultDatabase provides a default database to use for configuration when no configuration is found
	DefaultDatabase = "bludgeon"
	//DefaultParseTime provides a default value for whether or not to parse time for configuration when no configuration is found
	DefaultParseTime = true
	//DefaultDatabasePath provides a default database filepath to use for configuration when no configuration is found
	DefaultDatabasePath = "bludgeon.db"
	//DefaultMysqlPort provides a default port for mysql databases to use for configuration when no configuration is found
	DefaultMysqlPort = "3306"
	//DefaultPostgresPort provides a default port for postgres databases to use for configuration when no configuration is found
	DefaultPostgresPort = "5432"
	//DefaultHostname provides a default hostname to connect to databases to use for configuration when no configuration is found
	DefaultHostname = "Localhost"
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
	//ErrDriverUnsupported provides a string to return as an error if a driver isn't supported
	ErrDriverUnsupported string = "Configured driver: %s, not supported"
	//ErrDatabaseNotNil provides a string to return as an error if you attempt to connect to an already initialized database
	ErrDatabaseNotNil string = "Internal database pointer is not nil, reconnect or close to connect"
	//ErrUpdateFailed provides a string to return as an error if an update fails and the result returns 0 rows affected
	ErrUpdateFailed string = "Update failed, no rows affected"
	//ErrDeleteFailed provides a string to return as an error if a delete fails and the result returns 0 rows affected
	ErrDeleteFailed string = "Delete failed, id not found"
	//ErrQueryFailed provides a string to return as an error in the event a query fails and no other error is returned
	ErrQueryFailed string = "Query: \"%s\", failed"
)

//Configuration is a struct that contains al lthe possible configuration for supported database drivers
type Configuration struct {
	Driver          string        `json:"Driver"`          //go sql driver to user
	DataSource      string        `json:"-"`               //data source
	FilePath        string        `json:"FilePath"`        //filepath for sqlite
	Hostname        string        `json:"Hostname"`        //hostame to user to access the database
	Port            string        `json:"Port"`            //port to connect to
	Username        string        `json:"Username"`        //username to authenticate with
	Password        string        `json:"Password"`        //password to authenticate with
	Database        string        `json:"Database"`        //database to connect to
	ParseTime       bool          `json:"ParseTime"`       //whether or not to parse time
	UseTransactions bool          `json:"UseTransactions"` //whether or not to use transactions
	Timeout         time.Duration `json:"Timeout"`         //how long to wait with configuration
}

//configuration constants
const (
	DefaultTimeout time.Duration = 5 * time.Second
)

//query constants
const (
	TableTimer            string = "timer"
	QueryTimerCreateTable string = `CREATE TABLE IF NOT EXISTS ` + TableTimer + ` (
		id BIGINT NOT NULL AUTO_INCREMENT,
		uuid TEXT,
		activesliceid BIGINT,
		activesliceuuid TEXT,    
		start BIGINT,
		finish BIGINT,
		elapsedtime BIGINT,
		INDEX(id, uuid(36)),
	
		PRIMARY KEY (id)
		-- FOREIGN KEY (employeeid)
			-- REFERENCES employee(id)
			-- ON UPDATE CASCADE ON DELETE RESTRICT
	)ENGINE=InnoDB;`
	QueryTimerUpsertf string = "INSERT into %s (activesliceid, uuid, activesliceuuid, start, finish, elapsedtime) VALUES(?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE activesliceid=?, uuid=?, activesliceuuid=?, start=?, finish=?, elapsedtime=?"
	QueryTimerDeletef string = "DELETE FROM %s where uuid = ?"
	QueryTimerSelectf string = "SELECT * from %s WHERE uuid = \"%s\""
	//
	TableTimeSlice            string = "timeslice"
	QueryTimeSliceCreateTable string = `CREATE TABLE IF NOT EXISTS ` + TableTimeSlice + ` (
		id BIGINT NOT NULL AUTO_INCREMENT,
		uuid TEXT,
		timerid BIGINT,
		timeruuid TEXT,    
		start BIGINT,
		finish BIGINT,
		elapsedtime BIGINT,
		INDEX(id, uuid(36)),
	
		PRIMARY KEY (id)
		-- FOREIGN KEY (employeeid)
			-- REFERENCES employee(id)
			-- ON UPDATE CASCADE ON DELETE RESTRICT
	)ENGINE=InnoDB;`
	QueryTimeSliceUpsertf string = "INSERT into %s (timerid, uuid, timeruuid, start, finish, elapsedtime) VALUES(?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE timerid=?, uuid=?, timeruuid=?, start=?, finish=?, elapsedtime=?"
	QueryTimeSliceDeletef string = "DELETE FROM %s where uuid = ?"
	QueryTimeSliceSelectf string = "SELECT * from %s WHERE uuid = \"%s\""
)

// //TableProject is the string defining the name of the project table
// TableProject string = "project"
// //TableClient is the string defining the name of the client table
// TableClient string = "client"
// //TableEmployee is the string defining the name of the employee table
// TableEmployee string = "employee"
