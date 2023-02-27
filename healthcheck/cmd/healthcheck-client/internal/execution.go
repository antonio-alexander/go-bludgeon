package internal

import (
	client "github.com/antonio-alexander/go-bludgeon/healthcheck/client"
	grpcclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/grpc"
	restclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

func getConfig(pwd string, args []string, envs map[string]string) *Configuration {
	config := NewConfiguration()
	config.Default(pwd)
	switch {
	default:
		config.FromEnv(pwd, envs)
	case len(args) > 0:
		config.FromArgs(pwd, args)
	}
	return config
}

func parameterize(config *Configuration) (interface {
	internal_logger.Logger
	internal_logger.Printer
}, interface {
	internal.Initializer
	client.Client
}, error) {
	logger := internal_logger.New()
	logger.Configure(config.Logger)
	switch config.ClientType {
	default:
		client := restclient.New()
		client.SetUtilities(logger)
		if err := client.Configure(config.Rest); err != nil {
			return nil, nil, err
		}
		return logger, client, nil
	case "grpc":
		client := grpcclient.New()
		client.SetUtilities(logger)
		if err := client.Configure(config.Grpc); err != nil {
			return nil, nil, err
		}
		return logger, client, nil
	}
}

func initialize(parameters ...interface{}) error {
	for _, p := range parameters {
		if p, ok := p.(internal.Initializer); ok {
			if err := p.Initialize(); err != nil {
				return err
			}
		}
	}
	return nil
}

func shutdown(parameters ...interface{}) {
	for _, p := range parameters {
		if p, ok := p.(internal.Initializer); ok {
			p.Shutdown()
		}
	}
}
