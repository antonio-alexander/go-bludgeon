package internal

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"

	meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
)

const (
	EnvNameMetaType            string = "BLUDGEON_META_TYPE"
	EnvNameServiceRestEnabled  string = "BLUDGEON_REST_ENABLED"
	EnvNameServiceGrpcEnabled  string = "BLUDGEON_GRPC_ENABLED"
	EnvNameServiceKafkaEnabled string = "BLUDGEON_KAFKA_ENABLED"
)

const (
	DefaultRestEnabled  bool      = false
	DefaultGrpcEnabled  bool      = true
	DefaultKafkaEnabled bool      = true
	DefaultMetaType     meta.Type = meta.TypeMemory
)

type Configuration struct {
	MetaType     meta.Type `json:"type"`
	GrpcEnabled  bool      `json:"grpc_enabled"`
	RestEnabled  bool      `json:"rest_enabled"`
	KafkaEnabled bool      `json:"kafka_enabled"`
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
	c.RestEnabled = DefaultRestEnabled
	c.KafkaEnabled = DefaultKafkaEnabled
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if s, ok := envs[EnvNameMetaType]; ok {
		c.MetaType = meta.AtoType(s)
	}
	if s, ok := envs[EnvNameServiceRestEnabled]; ok {
		if restEnabled, err := strconv.ParseBool(s); err == nil {
			c.RestEnabled = restEnabled
		}
	}
	if s, ok := envs[EnvNameServiceKafkaEnabled]; ok {
		if kafkaEnabled, err := strconv.ParseBool(s); err == nil {
			c.KafkaEnabled = kafkaEnabled
		}
	}
}
