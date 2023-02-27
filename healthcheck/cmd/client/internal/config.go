package internal

import (
	"encoding/json"
	"flag"
	"io/fs"
	"os"

	grpcclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/grpc"
	restclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/rest"
)

const (
	EnvNameClientType = "BLUDGEON_HEALTHCHECK_CLIENT_TYPE"
)

const (
	DefaultClientType string = "rest"
)

type Configuration struct {
	ClientType string `json:"client_type"`
	Rest       *restclient.Configuration
	Grpc       *grpcclient.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Rest: restclient.NewConfiguration(),
		Grpc: grpcclient.NewConfiguration(),
	}
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
	c.ClientType = DefaultClientType
	c.Grpc.Default()
	c.Rest.Default()
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if s, ok := envs[EnvNameClientType]; ok {
		c.ClientType = s
	}
	c.Grpc.FromEnv(envs)
	c.Rest.FromEnv(envs)
}

func (c *Configuration) FromArgs(pwd string, args []string) error {
	var clientType, clientAddress, clientPort string

	cli := flag.NewFlagSet("healthcheck-client", flag.ExitOnError)
	cli.StringVar(&clientType, "client-type", "rest", "healthcheck client type")
	cli.StringVar(&clientAddress, "client-address", "localhost", "healthcheck client address")
	cli.StringVar(&clientPort, "client-port", "9030", "healthcheck client port")
	if err := cli.Parse(args); err != nil {
		return err
	}
	c.ClientType = clientType
	switch c.ClientType {
	default:
		c.Rest.Address = clientAddress
		c.Rest.Port = clientPort
	case "grpc":
		c.Grpc.Address = clientAddress
		c.Grpc.Port = clientPort
	}
	return nil
}
