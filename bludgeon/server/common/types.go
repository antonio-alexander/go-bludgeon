package bludgeonservercommon

import (
	"time"
)

//environmental variables
const (
	EnvNameDatabaseAddress      string = "DATABASE_ADDRESS"
	EnvNameDatabasePort         string = "DATABASE_PORT"
	EnvNameBludgeonMetaType     string = "BLUDGEON_META_TYPE"
	EnvNameBludgeonRestAddress  string = "BLUDGEON_REST_ADDRESS"
	EnvNameBludgeonRestPort     string = "BLUDGEON_REST_PORT"
	EnvNameBludgeonMetaJSONFile string = "BLUDGEON_META_JSON_FILE"
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

//defaults
const (
	DefaultDatabaseAddress      string = "127.0.0.1"
	DefaultDatabasePort         string = "3306"
	DefaultBludgeonMeta         string = "json"
	DefaultBludgeonRestAddress  string = ""
	DefaultBludgeonRestPort     string = "8080"
	DefaultBludgeonMetaJSONFile string = "data/bludgeon.json"
)
