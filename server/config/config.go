package config

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/pkg/errors"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
)

const (
	EnvNameRemoteType string = "BLUDGEON_REMOTE_TYPE"
	EnvNameMetaType   string = "BLUDGEON_META_TYPE"
)

const (
	DefaultConfigFile string            = "bludgeon_server_config.json"
	DefaultConfigPath string            = "config"
	DefaultMetaType   common.MetaType   = common.MetaTypeJSON
	DefaultRemoteType common.RemoteType = common.RemoteTypeRest
)

const (
	ErrRemoteTypeEmpty string = "remote type empty"
	ErrMetaTypeEmpty   string = "meta type empty"
)

type Configuration struct {
	RemoteType common.RemoteType
	RemoteRest *Rest
	MetaType   common.MetaType
	MetaMySQL  *metamysql.Configuration
	MetaJSON   *metajson.Configuration
}

func New() *Configuration {
	return &Configuration{
		RemoteRest: &Rest{},
		MetaMySQL:  &metamysql.Configuration{},
		MetaJSON:   &metajson.Configuration{},
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
	c.RemoteType = DefaultRemoteType
	c.RemoteRest.Default()
	c.MetaType = DefaultMetaType
	c.MetaJSON.Default(pwd)
	c.MetaMySQL.Default()
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if remoteType, ok := envs[EnvNameRemoteType]; ok {
		c.RemoteType = common.AtoRemoteType(remoteType)
	}
	if metaType, ok := envs[EnvNameMetaType]; ok {
		c.MetaType = common.AtoMetaType(metaType)
	}
	c.RemoteRest.FromEnv(pwd, envs)
	c.MetaJSON.FromEnv(pwd, envs)
	c.MetaMySQL.FromEnv(pwd, envs)
}

func (c *Configuration) Validate() (err error) {
	if c.MetaType == "" {
		err = errors.New(ErrMetaTypeEmpty)

		return
	}
	if c.RemoteType == "" {
		err = errors.New(ErrRemoteTypeEmpty)

		return
	}
	if err = c.RemoteRest.Validate(); err != nil {
		return
	}
	if err = c.MetaJSON.Validate(); err != nil {
		return
	}
	if err = c.MetaMySQL.Validate(); err != nil {
		return
	}

	return
}
