package config

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"strconv"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	meta "github.com/antonio-alexander/go-bludgeon/internal/meta"

	service_grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	service_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	meta_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	meta_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/pkg/errors"
)

type ConfigurationMeta struct {
	Type  meta.Type                 `json:"type"`
	File  *meta_file.Configuration  `json:"file"`
	Mysql *meta_mysql.Configuration `json:"mysql"`
}

type ConfigurationServer struct {
	RestEnabled bool                        `json:"rest_enabled"`
	Rest        *service_rest.Configuration `json:"rest"`
	GrpcEnabled bool                        `json:"grpc_enabled"`
	Grpc        *service_grpc.Configuration `json:"grpc"`
}

type Configuration struct {
	Logger *logger.Configuration `json:"logger"`
	Meta   ConfigurationMeta     `json:"meta"`
	Server ConfigurationServer   `json:"service"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Logger: new(logger.Configuration),
		Meta: ConfigurationMeta{
			File:  new(meta_file.Configuration),
			Mysql: new(meta_mysql.Configuration),
		},
		Server: ConfigurationServer{
			Rest: new(service_rest.Configuration),
			Grpc: new(service_grpc.Configuration),
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
	bytes, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, bytes, fs.FileMode(0644))
}

func (c *Configuration) Default(pwd string) {
	c.Meta.Type = DefaultMetaType
	c.Logger.Default()
	c.Meta.File.Default(pwd)
	c.Meta.Mysql.Default()
	c.Server.Rest.Default()
	c.Server.Grpc.Default()
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if metaType, ok := envs[EnvNameMetaType]; ok {
		c.Meta.Type = meta.AtoType(metaType)
	}
	if s, ok := envs[EnvNameServiceRestEnabled]; ok {
		if restEnabled, err := strconv.ParseBool(s); err == nil {
			c.Server.RestEnabled = restEnabled
		}
	}
	if s, ok := envs[EnvNameServiceGrpcEnabled]; ok {
		if grpcEnabled, err := strconv.ParseBool(s); err == nil {
			c.Server.GrpcEnabled = grpcEnabled
		}
	}
	c.Logger.FromEnv(pwd, envs)
	c.Meta.File.FromEnv(pwd, envs)
	c.Meta.Mysql.FromEnv(pwd, envs)
	c.Server.Rest.FromEnv(pwd, envs)
	c.Server.Grpc.FromEnv(pwd, envs)
}

func (c *Configuration) Validate() (err error) {
	if c.Meta.Type == "" {
		return errors.New(ErrMetaTypeEmpty)
	}
	if err = c.Server.Rest.Validate(); err != nil {
		return
	}
	if err = c.Server.Grpc.Validate(); err != nil {
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
