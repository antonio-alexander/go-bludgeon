package bludgeonclientcli

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

//arguments
const (
	ArgCommand      string = "command"
	ArgType         string = "type"
	ArgTimer        string = "timer"
	ArgTimerID      string = "id"
	ArgTimerStart   string = "start"
	ArgTimerFinish  string = "finish"
	ArgTimerComment string = "comment"
)

//defaults
const (
	DefaultCommand      string = ""
	DefaultType         string = ""
	DefaultTimerID      string = ""
	DefaultTimerStart   int64  = 0
	DefaultTimerFinish  int64  = 0
	DefaultTimerComment string = ""
)

//usage
const (
	UsageCommand      string = "Command for operation to attempt"
	UsageType         string = "Type of object"
	UsageTimerID      string = "ID for the timer"
	UsageTimerStart   string = "Timer start time"
	UsageTimerFinish  string = "Timer finish time"
	UsageTimerComment string = "Comment for the timer"
)

type Options struct {
	Command bludgeon.CommandClient //command
	Timer   bludgeon.Timer         //timer object
}
