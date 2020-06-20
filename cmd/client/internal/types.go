package client

import (
	"time"
)

// These variables are populated at build time
// REFERENCE: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
// to find where the variables are...
//  go tool nm ./app | grep app
var (
	Version   string
	GitCommit string
	GitBranch string
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

type CommandData struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	PauseTime  time.Time
}
