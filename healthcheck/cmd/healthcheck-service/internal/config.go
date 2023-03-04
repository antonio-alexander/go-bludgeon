package internal

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"
)

const (
	EnvNameServiceRestEnabled = "BLUDGEON_REST_ENABLED"
	EnvNameServiceGrpcEnabled = "BLUDGEON_GRPC_ENABLED"
)

const (
	DefaultServiceRestEnabled = true
	DefaultServiceGrpcEnabled = true
)

type Configuration struct {
	ServiceGrpcEnabled bool `json:"service_grpc_enabled"`
	ServiceRestEnabled bool `json:"service_rest_enabled"`
}

func (c *Configuration) Read(configFile string) error {
	bytes, err := os.ReadFile(configFile)
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
	return os.WriteFile(configFile, bytes, fs.FileMode(0644))
}

func (c *Configuration) Default(pwd string) {
	c.ServiceRestEnabled = DefaultServiceRestEnabled
	c.ServiceGrpcEnabled = DefaultServiceGrpcEnabled
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if s, ok := envs[EnvNameServiceRestEnabled]; ok {
		if serviceRestEnabled, err := strconv.ParseBool(s); err == nil {
			c.ServiceRestEnabled = serviceRestEnabled
		}
	}
	if s, ok := envs[EnvNameServiceGrpcEnabled]; ok {
		if serviceGrpcEnabled, err := strconv.ParseBool(s); err == nil {
			c.ServiceGrpcEnabled = serviceGrpcEnabled
		}
	}
}
