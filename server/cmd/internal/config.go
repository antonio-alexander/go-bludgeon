package internal

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	meta "github.com/antonio-alexander/go-bludgeon/meta"
	server "github.com/antonio-alexander/go-bludgeon/server"

	meta_file "github.com/antonio-alexander/go-bludgeon/meta/file"
	meta_mysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
	server_rest "github.com/antonio-alexander/go-bludgeon/server/rest"

	"github.com/pkg/errors"
)

const (
	DefaultConfigFile string      = "bludgeon_server_config.json"
	DefaultConfigPath string      = "config"
	DefaultMetaType   meta.Type   = meta.TypeFile
	DefaultServerType server.Type = server.TypeREST
	EnvNameRemoteType string      = "BLUDGEON_REMOTE_TYPE"
	EnvNameMetaType   string      = "BLUDGEON_META_TYPE"
	ErrMetaTypeEmpty  string      = "meta type empty"
)

type ConfigurationMeta struct {
	Type  meta.Type                 `json:"type"`
	File  *meta_file.Configuration  `json:"file"`
	Mysql *meta_mysql.Configuration `json:"mysql"`
}

type ConfigurationServer struct {
	Type server.Type                `json:"type"`
	Rest *server_rest.Configuration `json:"rest"`
}

type Configuration struct {
	Meta   ConfigurationMeta   `json:"meta"`
	Server ConfigurationServer `json:"server"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Meta: ConfigurationMeta{
			Type:  "",
			File:  &meta_file.Configuration{},
			Mysql: &meta_mysql.Configuration{},
		},
		Server: ConfigurationServer{
			Type: "",
			Rest: &server_rest.Configuration{},
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
