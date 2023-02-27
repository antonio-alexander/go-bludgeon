package internal

import (
	"flag"

	grpcclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/grpc"
	restclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/rest"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
)

const (
	EnvNameClientType = "BLUDGEON_HEALTHCHECK_CLIENT_TYPE"
)

const (
	DefaultClientType string       = "rest"
	DefaultLogPrefix  string       = "healthcheck-client"
	DefaultLogLevel   logger.Level = logger.Info
)

type Configuration struct {
	ClientType string
	Logger     *logger.Configuration
	Rest       *restclient.Configuration
	Grpc       *grpcclient.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Logger: new(logger.Configuration),
		Rest:   restclient.NewConfiguration(),
		Grpc:   grpcclient.NewConfiguration(),
	}
}

func (c *Configuration) Default(pwd string) {
	c.ClientType = DefaultClientType
	c.Logger.Level = DefaultLogLevel
	c.Logger.Prefix = DefaultLogPrefix
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
	c.Rest.Address = clientAddress
	c.Rest.Port = clientPort
	c.Grpc.Address = clientAddress
	c.Grpc.Port = clientPort
	return nil
}
