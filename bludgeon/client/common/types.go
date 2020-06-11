package bludgeonclient

import (
	"time"
)

//Configuration
type Configuration struct {
	Meta struct {
		Type   string                 `json:"Meta"`
		Config map[string]interface{} `json:"Config"`
	}
	Remote struct {
		Type   string                 `json:"Type"`
		Config map[string]interface{} `json:"Config"`
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

type Cache struct {
	TimerID string `json:"TimerID"`
}

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}
