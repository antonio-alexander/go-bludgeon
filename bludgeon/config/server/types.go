package bludgeonconfigserver

import (
	"encoding/json"

	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/common"
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

//defaults
const (
	DefaultDatabaseAddress      string = "127.0.0.1"
	DefaultDatabasePort         string = "3306"
	DefaultBludgeonMeta         string = "json"
	DefaultBludgeonRestAddress  string = ""
	DefaultBludgeonRestPort     string = "8080"
	DefaultBludgeonMetaJSONFile string = "data/bludgeon.json"
)

type Configuration struct {
	Meta struct {
		Type   string                     `json:"Type"`
		Config map[string]json.RawMessage `json:"Config"`
	} `json:"Meta"`
	Server server.Configuration `json:"Server"`
	Rest   struct {
		Address string `json:"Address"`
		Port    string `json:"Port"`
	} `json:"Rest"`
}
