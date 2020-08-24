package bludgeonclientcli

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
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
	Command    string              //command
	Address    string              //address
	Port       string              //port
	MetaType   bludgeon.MetaType   //meta type
	RemoteType bludgeon.RemoteType //remote type
	ObjectType bludgeon.ObjectType //object type
	Timer      bludgeon.Timer      //timer object
	TimeSlice  bludgeon.TimeSlice  //time slice object
}
