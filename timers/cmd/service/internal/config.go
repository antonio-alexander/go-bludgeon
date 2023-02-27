package internal

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"

	meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
)

const (
	EnvNameMetaType                  = "BLUDGEON_META_TYPE"
	EnvNameServiceRestEnabled        = "BLUDGEON_REST_ENABLED"
	EnvNameServiceGrpcEnabled        = "BLUDGEON_GRPC_ENABLED"
	EnvNameClientChangesRestEnabled  = "BLUDGEON_CHANGES_CLIENT_REST_ENABLED"
	EnvNameClientChangesKafkaEnabled = "BLUDGEON_CHANGES_CLIENT_KAFKA_ENABLED"
)

const (
	DefaultMetaType                  = meta.TypeMySQL
	DefaultServiceRestEnabled        = true
	DefaultServiceGrpcEnabled        = true
	DefaultClientChangesRestEnabled  = true
	DefaultClientChangesKafkaEnabled = true
)

type Configuration struct {
	MetaType                  meta.Type `json:"type"`
	ServiceGrpcEnabled        bool      `json:"service_grpc_enabled"`
	ServiceRestEnabled        bool      `json:"service_rest_enabled"`
	ClientChangesRestEnabled  bool      `json:"client_changes_rest_enabled"`
	ClientChangesKafkaEnabled bool      `json:"client_changes_kafka_enabled"`
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
	c.MetaType = DefaultMetaType
	c.ServiceRestEnabled = DefaultServiceRestEnabled
	c.ServiceGrpcEnabled = DefaultServiceGrpcEnabled
	c.ClientChangesRestEnabled = DefaultClientChangesRestEnabled
	c.ClientChangesKafkaEnabled = DefaultClientChangesKafkaEnabled
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if s, ok := envs[EnvNameMetaType]; ok {
		c.MetaType = meta.AtoType(s)
	}
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
	if s, ok := envs[EnvNameClientChangesRestEnabled]; ok {
		if clientChangesRestEnabled, err := strconv.ParseBool(s); err == nil {
			c.ClientChangesRestEnabled = clientChangesRestEnabled
		}
	}
	if s, ok := envs[EnvNameClientChangesKafkaEnabled]; ok {
		if clientChangesKafkaEnabled, err := strconv.ParseBool(s); err == nil {
			c.ClientChangesKafkaEnabled = clientChangesKafkaEnabled
		}
	}
}
