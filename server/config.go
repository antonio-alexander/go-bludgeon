package server

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	common "github.com/antonio-alexander/go-bludgeon/common"
	rest "github.com/antonio-alexander/go-bludgeon/server/rest"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"

	"github.com/pkg/errors"
)

const (
	EnvNameRemoteType string = "BLUDGEON_REMOTE_TYPE"
	EnvNameMetaType   string = "BLUDGEON_META_TYPE"
)

const (
	DefaultConfigFile string          = "bludgeon_server_config.json"
	DefaultConfigPath string          = "config"
	DefaultMetaType   common.MetaType = common.MetaTypeJSON
)

const (
	ErrMetaTypeEmpty string = "meta type empty"
)

type Server struct {
	Rest *rest.Configuration `json:"rest"`
}

type Meta struct {
	MySQL *metamysql.Configuration `json:"mysql"`
	JSON  *metajson.Configuration  `json:"json"`
}

type Configuration struct {
	MetaType common.MetaType `json:"meta_type"`
	Server   Server          `json:"server"`
	Meta     Meta            `json:"meta"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Server: Server{
			Rest: &rest.Configuration{},
		},
		Meta: Meta{
			MySQL: &metamysql.Configuration{},
			JSON:  &metajson.Configuration{},
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
	c.Server.Rest.Default()
	c.MetaType = DefaultMetaType
	c.Meta.JSON.Default(pwd)
	c.Meta.MySQL.Default()
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if metaType, ok := envs[EnvNameMetaType]; ok {
		c.MetaType = common.AtoMetaType(metaType)
	}
	c.Server.Rest.FromEnv(pwd, envs)
	c.Meta.JSON.FromEnv(pwd, envs)
	c.Meta.MySQL.FromEnv(pwd, envs)
}

func (c *Configuration) Validate() (err error) {
	if c.MetaType == "" {
		err = errors.New(ErrMetaTypeEmpty)

		return
	}
	if err = c.Server.Rest.Validate(); err != nil {
		return
	}
	if err = c.Meta.JSON.Validate(); err != nil {
		return
	}
	if err = c.Meta.MySQL.Validate(); err != nil {
		return
	}

	return
}
