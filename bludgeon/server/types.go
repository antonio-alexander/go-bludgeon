package bludgeonserver

import "time"

//error constants
const (
	ErrStarted    string = "server started"
	ErrNotStarted string = "server not started"
)

//Configuration provides a struct to define the configurable elements of a server
type Configuration struct {
	TokenWait int64 //how long a token is valid (seconds)
}

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}
