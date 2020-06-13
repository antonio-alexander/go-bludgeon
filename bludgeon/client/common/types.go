package bludgeonclient

import (
	"time"
)

//Configuration
type Configuration struct {
	ServerAddress string `json:"ServerAddress"`
	ServerPort    string `json:"ServerPort"`
	ClientAddress string `json:"ClientAddress"`
	ClientPort    string `json:"ClientPort"`
	// Task          int64  `json:"Task"`
	// Employee      int64  `json:"Employee"`
}

type Cache struct {
	TimerID string `json:"TimerID"`
}

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}
