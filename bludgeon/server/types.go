package bludgeonserver

import (
	"net/http"
	"time"

	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

//error constants
const (
	ErrStarted    string = "server started"
	ErrNotStarted string = "server not started"
)

//Configuration provides a struct to define the configurable elements of a server

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}

type Options struct {
	Configuration string //the filepath to the configuration
}

//endpoint constants
const (
	GET    = http.MethodGet
	PUT    = http.MethodPut
	POST   = http.MethodPost
	DELETE = http.MethodDelete
)

//arguments
const (
	ArgConfiguration string = "config"
)

//defaults
const ()

//usage
const (
	UsageConfiguration string = "The path to the configuration"
)

type Configuration struct {
	Meta struct {
		Type string `json:"Meta"`
		JSON struct {
			File string `json:"File"`
		} `json:"JSON"`
		MySQL mysql.Configuration `json:"Mysql"`
	}
	Remote struct {
		Type       string `json:"Type"`
		RestClient struct {
			Address string        `json:"Address"`
			Port    string        `json:"Port"`
			Timeout time.Duration `json:"Timeout"`
		} `json:"RestClient"`
	} `json:"Remote"`
	Server struct {
		TokenWait time.Duration //how long a token is valid (seconds)
		Rest      struct {
			Address string `json:"Address"`
			Port    string `json:"Port"`
		} `json:"Rest"`
	} `json:"Server"`
}
