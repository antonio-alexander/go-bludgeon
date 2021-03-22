package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
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
	EnvNameBludgeonAddress    string = "BLUDGEON_ADDRESS"
	EnvNameBludgeonPort       string = "BLUDGEON_PORT"
	EnvNameBluderonTimeout    string = "BLUDGEON_TIMEOUT"
)

//defaults
const (
	DefaultMetaType   common.MetaType   = common.MetaTypeJSON
	DefaultRemoteType common.RemoteType = common.RemoteTypeRest
	DefaultAddress    string            = "127.0.0.1"
	DefaultPort       string            = "8080"
	DefaultTimeout    time.Duration     = 10 * time.Second
)

type Configuration struct {
	MetaType   common.MetaType   `json:"MetaType"`
	RemoteType common.RemoteType `json:"RemoteType"`
	Address    string            `json:"Address"`
	Port       string            `json:"Port"`
	Timeout    time.Duration     `json:"Timeout"`
	Meta       map[common.MetaType]json.RawMessage
	Remote     map[common.RemoteType]json.RawMessage
}

func (c *Configuration) Default(pwd string) {
	c.MetaType = common.MetaType(DefaultMetaType)
	c.RemoteType = common.RemoteType(DefaultRemoteType)
	c.Address = DefaultAddress
	c.Port = DefaultPort
	c.Timeout = DefaultTimeout
	c.Meta = make(map[common.MetaType]json.RawMessage)
	c.Remote = make(map[common.RemoteType]json.RawMessage)
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) (err error) {
	if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
		c.MetaType = common.AtoMetaType(metaType)
	}
	if remoteType, ok := envs[EnvNameBludgeonRemoteType]; ok {
		c.RemoteType = common.AtoRemoteType(remoteType)
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
