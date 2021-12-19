package internal

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	meta "github.com/antonio-alexander/go-bludgeon/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/meta/file"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
	server "github.com/antonio-alexander/go-bludgeon/server"
	serverrest "github.com/antonio-alexander/go-bludgeon/server/rest"

	"github.com/pkg/errors"
)

const (
	EnvNameRemoteType string      = "BLUDGEON_REMOTE_TYPE"
	EnvNameMetaType   string      = "BLUDGEON_META_TYPE"
	DefaultConfigFile string      = "bludgeon_server_config.json"
	DefaultConfigPath string      = "config"
	DefaultMetaType   meta.Type   = meta.TypeFile
	ErrMetaTypeEmpty  string      = "meta type empty"
	DefaultServerType server.Type = server.TypeREST
)

type ConfigurationMeta struct {
	Type  meta.Type
	File  *metafile.Configuration
	Mysql *metamysql.Configuration
}

type ConfigurationServer struct {
	Type server.Type
	Rest *serverrest.Configuration
}

type Configuration struct {
	Meta   ConfigurationMeta
	Server ConfigurationServer
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Meta: ConfigurationMeta{
			Type:  "",
			File:  &metafile.Configuration{},
			Mysql: &metamysql.Configuration{},
		},
		Server: ConfigurationServer{
			Type: "",
			Rest: &serverrest.Configuration{},
		},
	}
}

func (c *Configuration) Read(configFile string) (err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadFile(configFile); err != nil {
		return
	}
	err = json.Unmarshal(bytes, c)

	return
}

func (c *Configuration) Write(configFile string) (err error) {
	var bytes []byte

	//REVIEW: do we need to ensure that the path exists?
	if bytes, err = json.MarshalIndent(&c, "", "  "); err != nil {
		return
	}
	err = ioutil.WriteFile(configFile, bytes, fs.FileMode(0644))

	return
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
