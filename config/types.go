package bludgeonconfig

import (
	"encoding/json"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/config"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/config"
)

//common constants
const (
	fileMode os.FileMode = 0644
)

//error constants
const (
	ErrUnsupportedConfigf string = "Unsupported Type: %T"
	ErrMetaNotFoundf      string = "Meta not found: %s"
	ErrRemoteNotFoundf    string = "Meta not found: %s"
)

//environmental variables
const (
	EnvNameBludgeonMetaType   string = "BLUDGEON_META_TYPE"
	EnvNameBludgeonRemoteType string = "BLUDGEON_REMOTE_TYPE"
)

//defaults
const (
	DefaultMetaType   bludgeon.MetaType   = bludgeon.MetaTypeJSON
	DefaultRemoteType bludgeon.RemoteType = bludgeon.RemoteTypeRest
)

type Server struct {
	MetaType bludgeon.MetaType                     `json:"MetaType"`
	Meta     map[bludgeon.MetaType]json.RawMessage `json:"Meta"`
	Rest     rest.Configuration                    `json:"Rest"`
	Server   server.Configuration                  `json:"Server"`
}
