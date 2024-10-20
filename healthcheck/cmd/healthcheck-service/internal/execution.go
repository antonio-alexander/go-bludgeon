package internal

import (
	"path/filepath"

	logic "github.com/antonio-alexander/go-bludgeon/healthcheck/logic"
	servicegrpc "github.com/antonio-alexander/go-bludgeon/healthcheck/service/grpc"
	servicerest "github.com/antonio-alexander/go-bludgeon/healthcheck/service/rest"

	common "github.com/antonio-alexander/go-bludgeon/common"
	pkg_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	pkg_server_grpc "github.com/antonio-alexander/go-bludgeon/pkg/grpc/server"
	pkg_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	pkg_server_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/server"
)

func getConfig(pwd string, envs map[string]string) *Configuration {
	configFile := filepath.Join(pwd, pkg_config.DefaultConfigPath, pkg_config.DefaultConfigFile)
	config := new(Configuration)
	config.Default(pwd)
	if err := config.Read(configFile); err == nil {
		return config
	}
	config.FromEnv(pwd, envs)
	return config
}

func parameterize(config *Configuration) (interface {
	pkg_logger.Logger
	pkg_logger.Printer
	common.Configurer
}, []interface{}) {
	var parameters []interface{}

	logger := pkg_logger.New()
	healthCheckLogic := logic.New()
	healthCheckLogic.SetUtilities(logger)
	healthCheckLogic.SetParameters()
	parameters = append(parameters, healthCheckLogic)
	if config.ServiceRestEnabled {
		restServer := pkg_server_rest.New()
		healthCheckRestService := servicerest.New()
		healthCheckRestService.SetUtilities(logger)
		restServer.SetUtilities(logger)
		restServer.SetParameters(healthCheckRestService)
		healthCheckRestService.SetParameters(restServer, healthCheckLogic)
		parameters = append(parameters, restServer, healthCheckRestService)
	}
	if config.ServiceGrpcEnabled {
		healthCheckGrpcService := servicegrpc.New()
		healthCheckGrpcService.SetUtilities(logger)
		healthCheckGrpcService.SetParameters(healthCheckLogic, healthCheckGrpcService)
		grpcServer := pkg_server_grpc.New()
		grpcServer.SetUtilities(logger)
		grpcServer.SetParameters(healthCheckGrpcService)
		parameters = append(parameters, grpcServer, healthCheckGrpcService)
	}
	return logger, parameters
}

func configure(pwd string, envs map[string]string, parameters ...interface{}) error {
	//TODO: allow this to be able to accept configuration from a json
	// file
	for _, p := range parameters {
		switch p := p.(type) {
		case common.Configurer:
			if err := p.Configure(pkg_config.Envs(envs)); err != nil {
				return err
			}
		}
	}
	return nil
}

func initialize(parameters ...interface{}) error {
	for _, p := range parameters {
		if p, ok := p.(common.Initializer); ok {
			if err := p.Initialize(); err != nil {
				return err
			}
		}
	}
	return nil
}

func shutdown(parameters ...interface{}) {
	for _, p := range parameters {
		if p, ok := p.(common.Initializer); ok {
			p.Shutdown()
		}
	}
}
