package internal

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	service "github.com/antonio-alexander/go-bludgeon/employees/service"

	service_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/pkg/errors"
)

const ErrMetaTypeEmpty string = "meta type empty"

const (
	DefaultConfigFile string       = "bludgeon_service_config.json"
	DefaultConfigPath string       = "config"
	DefaultMetaType   meta.Type    = meta.TypeFile
	DefaultServerType service.Type = service.TypeREST
)

const (
	EnvNameServiceType string = "BLUDGEON_SERVICE_TYPE"
	EnvNameMetaType    string = "BLUDGEON_META_TYPE"
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

func (c *Configuration) Read(configFile string) error {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}

func (c *Configuration) Write(configFile string) error {
	//REVIEW: do we need to ensure that the path exists?
	bytes, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return err
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
