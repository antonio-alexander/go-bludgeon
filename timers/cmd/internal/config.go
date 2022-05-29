package internal

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	service "github.com/antonio-alexander/go-bludgeon/timers/service"

	service_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/pkg/errors"
)

const (
	DefaultConfigFile string       = "bludgeon_service_config.json"
	DefaultConfigPath string       = "config"
	DefaultMetaType   meta.Type    = meta.TypeFile
	DefaultServerType service.Type = service.TypeREST
	EnvNameRemoteType string       = "BLUDGEON_REMOTE_TYPE"
	EnvNameMetaType   string       = "BLUDGEON_META_TYPE"
	ErrMetaTypeEmpty  string       = "meta type empty"
)

type ConfigurationMeta struct {
	Type  meta.Type                     `json:"type"`
	File  *internal_file.Configuration  `json:"file"`
	Mysql *internal_mysql.Configuration `json:"mysql"`
}

type ConfigurationServer struct {
	Type service.Type                `json:"type"`
	Rest *service_rest.Configuration `json:"rest"`
}

type Configuration struct {
	Meta   ConfigurationMeta   `json:"meta"`
	Server ConfigurationServer `json:"service"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Meta: ConfigurationMeta{
			Type:  "",
			File:  &internal_file.Configuration{},
			Mysql: &internal_mysql.Configuration{},
		},
		Server: ConfigurationServer{
			Type: "",
			Rest: &service_rest.Configuration{},
		},
	}
}

func (c *Configuration) Read(configFile string) (err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadFile(configFile); err != nil {
		return
	}
	return json.Unmarshal(bytes, c)
}

func (c *Configuration) Write(configFile string) (err error) {
	var bytes []byte

	//REVIEW: do we need to ensure that the path exists?
	if bytes, err = json.MarshalIndent(&c, "", "  "); err != nil {
		return
	}
	return ioutil.WriteFile(configFile, bytes, fs.FileMode(0644))
}

func (c *Configuration) Default(pwd string) {
	c.Meta.Type = DefaultMetaType
	c.Meta.File.Default(pwd)
	c.Meta.Mysql.Default()
	c.Server.Type = DefaultServerType
	c.Server.Rest.Default()
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if metaType, ok := envs[EnvNameMetaType]; ok {
		c.Meta.Type = meta.AtoType(metaType)
	}
	c.Server.Rest.FromEnv(pwd, envs)
	c.Meta.File.FromEnv(pwd, envs)
	c.Meta.Mysql.FromEnv(pwd, envs)
}

func (c *Configuration) Validate() (err error) {
	if c.Meta.Type == "" {
		return errors.New(ErrMetaTypeEmpty)
	}
	if err = c.Server.Rest.Validate(); err != nil {
		return
	}
	if err = c.Meta.File.Validate(); err != nil {
		return
	}
	if err = c.Meta.Mysql.Validate(); err != nil {
		return
	}

	return
}
