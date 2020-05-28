package bludgeonmetamysql

import (
	"database/sql"
	"errors"
	"fmt"
)

//validateConfiguration is used to ensure that the values being configured make sense
// it's not necessarily to prevent a misconfiguration, but to use default values in the
// event a value doesn't exist
func validateConfiguration(configIn Configuration) (configOut Configuration, err error) {
	configOut = configIn
	switch configOut.Driver {
	case "mysql", "postgres":
		switch configOut.Driver {
		case "mysql":
			if configOut.Port == "" {
				configOut.Port = DefaultMysqlPort
			}
		case "postgres":
			if configOut.Port == "" {
				configOut.Port = DefaultPostgresPort
			}
		}
		if configOut.Database == "" {
			configOut.Database = DefaultDatabase
		}
		if configOut.Username == "" {
			configOut.Username = DefaultUsername
		}
		if configOut.Password == "" {
			configOut.Password = DefaultPassword
		}
		if configOut.Hostname == "" {
			configOut.Hostname = DefaultHostname
		}
	case "sqlite":
		if configOut.FilePath == "" {
			configOut.FilePath = DefaultDatabasePath
		}
	default:
		err = fmt.Errorf(ErrDriverUnsupported, configOut.Driver)
	}

	if configOut.Timeout <= 0 {
		configOut.Timeout = DefaultTimeout
	}

	return
}

//ConvertConfiguration will use a configuration and output a driver string and source
func convertConfiguration(config Configuration) (driver string, dataSource string, err error) {
	switch config.Driver {
	case "sqlite":
		//"sqlite3", "./foo.db"
		driver = "sqlite3"
		dataSource = config.FilePath
	case "mysql":
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		//user:password@tcp(localhost:5555)/dbname?charset=utf8
		driver = "mysql"
		dataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t", config.Username, config.Password, config.Hostname, config.Port, config.Database, config.ParseTime)
	case "postgres":
		//user=[username] password=[password] dbname=[dbname] sslmode=disable
		driver = "postgres"
		dataSource = fmt.Sprintf("user=[%s] password=[%s] dbname=[%s] sslmode=disable", config.Username, config.Password, config.Database)
	}

	return
}

//rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func rowsAffected(result sql.Result, errorString string) (err error) {
	var n int64

	//get the number of rows affected
	if n, err = result.RowsAffected(); err != nil {
		return
	}
	//return error string
	if n <= 0 {
		err = errors.New(errorString)
	}

	return
}

//lastInsertID will return the last id after an insert operation or an error if present
func lastInsertID(result sql.Result) (id int64, err error) {
	//get the last inserted id
	id, err = result.LastInsertId()

	return
}
