package bludgeonclient

import (
	"time"
)

//SerializedData
type SerializedData struct {
	//
}

//error constants
const (
	ErrStarted    string = "client started"
	ErrNotStarted string = "client not started"
)

type Cache struct {
	TimerID string `json:"TimerID"`
}

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}
