package bludgeonclientcli

import (
	common "github.com/antonio-alexander/go-bludgeon/internal/common"
)

//arguments
const (
	ArgCommand      string = "command"
	ArgAddress      string = "address"
	ArgPort         string = "port"
	ArgObjectType   string = "type"
	ArgMetaType     string = "meta"
	ArgRemoteType   string = "remote"
	ArgTimer        string = "timer"
	ArgTimerID      string = "id"
	ArgTimerStart   string = "start"
	ArgTimerFinish  string = "finish"
	ArgTimerComment string = "comment"
)

//defaults
const (
	DefaultCommand      string = ""
	DefaultAddress      string = "127.0.0.1"
	DefaultPort         string = "8080"
	DefaultObjectType   string = ""
	DefaultMetaType     string = ""
	DefaultRemoteType   string = ""
	DefaultTimerID      string = ""
	DefaultTimerStart   int64  = 0
	DefaultTimerFinish  int64  = 0
	DefaultTimerComment string = ""
)

//usage
const (
	UsageCommand      string = "Command for operation to attempt"
	UsageAddress      string = "Address to connect to"
	UsagePort         string = "Port to connect to"
	UsageObjectType   string = "Type of object"
	UsageMetaType     string = "Type of meta"
	UsageRemoteType   string = "Type of remote"
	UsageTimerID      string = "ID for the timer"
	UsageTimerStart   string = "Timer start time"
	UsageTimerFinish  string = "Timer finish time"
	UsageTimerComment string = "Comment for the timer"
)

type Options struct {
	Command    string            //command
	Address    string            //address
	Port       string            //port
	MetaType   common.MetaType   //meta type
	RemoteType common.RemoteType //remote type
	ObjectType common.ObjectType //object type
	Timer      common.Timer      //timer object
	TimeSlice  common.TimeSlice  //time slice object
}
