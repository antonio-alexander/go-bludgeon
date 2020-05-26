package bludgeonservercommon

import (
	"time"
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
