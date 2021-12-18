package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/client"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/meta/file"
)

//data constants
const (
	fileMode                  os.FileMode   = 0644
	ErrUnsupportedConfigf     string        = "Unsupported Type: %T"
	ErrMetaNotFoundf          string        = "Meta not found: %s"
	ErrRemoteNotFoundf        string        = "Meta not found: %s"
	EnvNameBludgeonMetaType   string        = "BLUDGEON_META_TYPE"
	EnvNameBludgeonClientType string        = "BLUDGEON_REMOTE_TYPE"
	EnvNameBludgeonAddress    string        = "BLUDGEON_ADDRESS"
	EnvNameBludgeonPort       string        = "BLUDGEON_PORT"
	EnvNameBluderonTimeout    string        = "BLUDGEON_TIMEOUT"
	DefaultMetaType           meta.Type     = meta.TypeFile
	DefaultClientType         client.Type   = client.TypeRest
	DefaultAddress            string        = "127.0.0.1"
	DefaultPort               string        = "8080"
	DefaultTimeout            time.Duration = 10 * time.Second
)

type Meta struct {
	Type meta.Type `json:"MetaType"`
	File *metafile.Configuration
}

type Client struct {
	Type client.Type
	Rest *ClientRest
}

type ClientRest struct {
	Address string
	Port    string
	Timeout time.Duration
}

type Configuration struct {
	Client Client
	Meta   Meta
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Client: Client{
			Type: "",
			Rest: &ClientRest{},
		},
		Meta: Meta{
			Type: "",
			File: &metafile.Configuration{},
		},
	}
}

func (c *Configuration) Default(pwd string) {
	c.Meta.Type = DefaultMetaType
	c.Client.Type = DefaultClientType
	c.Client.Rest.Address = DefaultAddress
	c.Client.Rest.Port = DefaultPort
	c.Client.Rest.Timeout = DefaultTimeout
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) (err error) {
	if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
		c.Meta.Type = meta.AtoType(metaType)
	}
	if remoteType, ok := envs[EnvNameBludgeonClientType]; ok {
		c.Client.Type = client.AtoType(remoteType)
	}
	//TODO: add code to get address from env
	//TODO: add code to get port from env
	//TODO: add code to get timeout from env
	err = c.Validate()

	return
}

func (c *Configuration) Validate() (err error) {
	return
}

func (c *Configuration) Read(configPath, pwd string, envs map[string]string) (err error) {
	var bytes []byte
	var exists bool

	//check if configPath exists, then swtich on the input
	// maintain the pointer for configuration, if config file
	// exists, read from it, otherwise, populate the defaults
	// then attempt to read from the environment and then
	// validate the configuration. In the event the provided
	// configuration is not a supported type, output an error
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		exists = true
	}
	if !exists {
		c.Default(pwd)
		if err = c.FromEnv(pwd, envs); err != nil {
			return
		}
	} else {
		if bytes, err = ioutil.ReadFile(configPath); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &c); err != nil {
			return
		}
	}
	err = c.Validate()

	return
}

func (c *Configuration) Write(configPath string) (err error) {
	var bytes []byte

	if bytes, err = json.MarshalIndent(&c, "", "    "); err != nil {
		return
	}
	err = ioutil.WriteFile(configPath, bytes, fileMode)

	return
}
