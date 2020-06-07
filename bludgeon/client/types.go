package bludgeonclient

import (
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

// "text/tabwriter"

//error constants
const (
	ErrStarted    string = "client started"
	ErrNotStarted string = "client not started"
)

//SerializedData
type SerializedData struct {
	//
}

//common constants
const (
	SQL_DRIVER = "sqlite3"
	HELP       = "Help Goes Here!"
	// tabwriterFlag = tabwriter.Debug //tabwriter.AlignRight | tabwriter.Debug
)

//Configuration
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
	Client struct {
		ServerAddress string
		ServerPort    string
		ClientAddress string
		ClientPort    string
		Task          int64
		Employee      int64
	} `json:"Client"`
}

//arguments
const (
	ArgCommand       string = "command"
	ArgType          string = "type"
	ArgTimer         string = "timer"
	ArgConfiguration string = "config"
	//
	ArgTimerID      string = "id"
	ArgTimerStart   string = "start"
	ArgTimerFinish  string = "finish"
	ArgTimerComment string = "comment"
)

//defaults
const (
	DefaultCommand string = ""
	DefaultType    string = ""
	//
	DefaultTimerID      string = ""
	DefaultTimerStart   int64  = 0
	DefaultTimerFinish  int64  = 0
	DefaultTimerComment string = ""
)

//usage
const (
	UsageCommand       string = "Command for operation to attempt"
	UsageType          string = "Type of object"
	UsageConfiguration string = "The path to the configuration"
	//
	UsageTimerID      string = "ID for the timer"
	UsageTimerStart   string = "Timer start time"
	UsageTimerFinish  string = "Timer finish time"
	UsageTimerComment string = "Comment for the timer"
)

type Options struct {
	Command       bludgeon.CommandClient //command
	Configuration string                 //the filepath to the configuration
	Timer         bludgeon.Timer         //timer object
}

type Cache struct {
	TimerID string `json:"TimerID"`
}
