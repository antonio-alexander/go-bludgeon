package bludgeonclientcli

import (
	client "github.com/antonio-alexander/go-bludgeon/client"
	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
)

//arguments
const (
	ArgCommand          string = "command"
	ArgAddress          string = "address"
	ArgPort             string = "port"
	ArgObjectType       string = "type"
	ArgMetaType         string = "meta"
	ArgClientType       string = "remote"
	ArgTimer            string = "timer"
	ArgTimerID          string = "id"
	ArgTimerStart       string = "start"
	ArgTimerFinish      string = "finish"
	ArgTimerComment     string = "comment"
	DefaultCommand      string = ""
	DefaultAddress      string = "127.0.0.1"
	DefaultPort         string = "8080"
	DefaultObjectType   string = ""
	DefaultMetaType     string = ""
	DefaultClientType   string = ""
	DefaultTimerID      string = ""
	DefaultTimerStart   int64  = 0
	DefaultTimerFinish  int64  = 0
	DefaultTimerComment string = ""
	UsageCommand        string = "Command for operation to attempt"
	UsageAddress        string = "Address to connect to"
	UsagePort           string = "Port to connect to"
	UsageObjectType     string = "Type of object"
	UsageMetaType       string = "Type of meta"
	UsageClientType     string = "Type of remote"
	UsageTimerID        string = "ID for the timer"
	UsageTimerStart     string = "Timer start time"
	UsageTimerFinish    string = "Timer finish time"
	UsageTimerComment   string = "Comment for the timer"
)

type Options struct {
	Command    string          //command
	Address    string          //address
	Port       string          //port
	MetaType   meta.Type       //meta type
	ClientType client.Type     //client type
	ObjectType data.ObjectType //object type
	Timer      data.Timer      //timer object
	TimeSlice  data.TimeSlice  //time slice object
}
