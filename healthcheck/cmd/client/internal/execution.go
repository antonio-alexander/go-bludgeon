package internal

import (
	"path/filepath"

	client "github.com/antonio-alexander/go-bludgeon/healthcheck/client"
	grpcclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/grpc"
	restclient "github.com/antonio-alexander/go-bludgeon/healthcheck/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_config "github.com/antonio-alexander/go-bludgeon/internal/config"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

func getConfig(pwd string, args []string, envs map[string]string) *Configuration {
	configFile := filepath.Join(pwd, internal_config.DefaultConfigPath, internal_config.DefaultConfigFile)
	config := NewConfiguration()
	config.Default(pwd)
	if err := config.Read(configFile); err == nil {
		return config
	}
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
}) {
	var client interface {
		internal.Configurer
		internal.Parameterizer
		internal.Initializer
		client.Client
	}

	logger := internal_logger.New()
	switch config.ClientType {
	default:
		client = restclient.New()
		client.SetUtilities(logger)
		client.SetParameters()
	case "grpc":
		client = grpcclient.New()
		client.SetUtilities(logger)
		client.SetParameters()
	}
	client.SetParameters(logger)
	return logger, client
}

func configure(pwd string, envs map[string]string, parameters ...interface{}) error {
	//TODO: allow this to be able to accept configuration from a json
	// file
	for _, p := range parameters {
		switch p := p.(type) {
		case internal.Configurer:
			if err := p.Configure(internal_config.Envs(envs)); err != nil {
				return err
			}
		}
	}
	return nil
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
