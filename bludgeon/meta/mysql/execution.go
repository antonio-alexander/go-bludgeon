package bludgeonmetamysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	config "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql/config"
)

func castConfiguration(element interface{}) (c config.Configuration, err error) {
	switch v := element.(type) {
	case json.RawMessage:
		err = json.Unmarshal(v, &c)
	case config.Configuration:
		c = v
	default:
		err = fmt.Errorf("Unsupported type: %t", element)
	}

	return
}

//ConvertConfiguration will use a configuration and output a driver string and source
func convertConfiguration(config config.Configuration) (driver string, dataSource string, err error) {
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
