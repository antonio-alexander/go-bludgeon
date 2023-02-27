package grpc

import (
	"context"

	client "github.com/antonio-alexander/go-bludgeon/healthcheck/client"
	data "github.com/antonio-alexander/go-bludgeon/healthcheck/data"
	pb "github.com/antonio-alexander/go-bludgeon/healthcheck/data/pb"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	logger.Logger
	pb.HealthChecksClient
	client interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		grpc.ClientConnInterface
	}
}

// New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	client.Client
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	return &grpcClient{
		Logger: logger.NewNullLogger(),
		client: grpcclient.New(),
	}
}

func (g *grpcClient) SetParameters(parameters ...interface{}) {
	g.client.SetParameters(parameters...)
}

func (g *grpcClient) SetUtilities(parameters ...interface{}) {
	g.client.SetUtilities(parameters...)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			g.Logger = p
		}
	}
}

func (g *grpcClient) Configure(items ...interface{}) error {
	var configuration *Configuration
	var envs map[string]string

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			configuration = v
		}
	}
	if configuration == nil {
		configuration = NewConfiguration()
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	if err := g.client.Configure(configuration.Configuration); err != nil {
		return err
	}
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (g *grpcClient) Initialize() error {
	if err := g.client.Initialize(); err != nil {
		return err
	}
	g.HealthChecksClient = pb.NewHealthChecksClient(g.client)
	return nil
}

func (g *grpcClient) Shutdown() {
	g.client.Shutdown()
}

func (g *grpcClient) HealthCheck(ctx context.Context) (*data.HealthCheck, error) {
	response, err := g.HealthChecksClient.Healthcheck(ctx, &pb.Empty{})
	return pb.ToHealthCheck(response.GetHealthcheck()), err
}
