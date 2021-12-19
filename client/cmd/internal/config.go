package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/client"
	rest "github.com/antonio-alexander/go-bludgeon/client/rest"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/meta/file"
)

//data constants
const (
	fileMode                  os.FileMode   = 0644
	ErrUnsupportedConfigf     string        = "unsupported Type: %T"
	ErrMetaNotFoundf          string        = "meta not found: %s"
	ErrRemoteNotFoundf        string        = "meta not found: %s"
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
	Type meta.Type
	File *metafile.Configuration
}

type Client struct {
	Type client.Type
	Rest *rest.Configuration
}

type Configuration struct {
	Client Client
	Meta   Meta
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Client: Client{
			Type: "",
			Rest: &rest.Configuration{},
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
	if clientType, ok := envs[EnvNameBludgeonClientType]; ok {
		c.Client.Type = client.AtoType(clientType)
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

func (c *Configuration) Read(configPath, pwd string, envs map[string]string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		c.Default(pwd)
		return c.FromEnv(pwd, envs)
	}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return err
	}
	return c.Validate()
}

func (c *Configuration) Write(configPath string) error {
	bytes, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, bytes, fileMode)
}
