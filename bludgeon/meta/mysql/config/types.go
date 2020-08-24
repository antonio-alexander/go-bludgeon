package bludgeonmetamysqlconfig

import (
	"time"
)

const (
	//ErrDriverUnsupported provides a string to return as an error if a driver isn't supported
	ErrDriverUnsupported string = "Configured driver: %s, not supported"
)

//environmental variables
const (
	EnvNameDatabaseAddress string = "DATABASE_ADDRESS"
	EnvNameDatabasePort    string = "DATABASE_PORT"
)

//defaults
const (
	DefaultDatabaseAddress string = "127.0.0.1"
	DefaultDatabasePort    string = "3306"
)

//defaults
const (
	DefaultDriver       string        = "mysql"
	DefaultUsername     string        = "bludgeon"
	DefaultPassword     string        = "bludgeon"
	DefaultDatabase     string        = "bludgeon"
	DefaultParseTime    bool          = true
	DefaultDatabasePath string        = "bludgeon.db"
	DefaultMysqlPort    string        = "3306"
	DefaultPostgresPort string        = "5432"
	DefaultHostname     string        = "Localhost"
	DefaultTimeout      time.Duration = 5 * time.Second
)

//Configuration is a struct that contains al lthe possible configuration for supported database drivers
type Configuration struct {
	Hostname        string        `json:"Hostname"`        //hostame to user to access the database
	Port            string        `json:"Port"`            //port to connect to
	Username        string        `json:"Username"`        //username to authenticate with
	Password        string        `json:"Password"`        //password to authenticate with
	Database        string        `json:"Database"`        //database to connect to
	UseTransactions bool          `json:"UseTransactions"` //whether or not to use transactions
	Timeout         time.Duration `json:"Timeout"`         //how long to wait with configuration
	Driver          string        `json:"Driver"`          //go sql driver to user
	DataSource      string        `json:"-"`               //data source
	FilePath        string        `json:"FilePath"`        //filepath for sqlite
	ParseTime       bool          `json:"ParseTime"`       //whether or not to parse time
}
