package internal

import (
	"path/filepath"

	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	metamemory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	metamysql "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	servicegrpc "github.com/antonio-alexander/go-bludgeon/employees/service/grpc"
	servicerest "github.com/antonio-alexander/go-bludgeon/employees/service/rest"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_config "github.com/antonio-alexander/go-bludgeon/internal/config"
	internal_server_grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_server_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
)

func getConfig(pwd string, envs map[string]string) *Configuration {
	configFile := filepath.Join(pwd, internal_config.DefaultConfigPath, internal_config.DefaultConfigFile)
	config := new(Configuration)
	config.Default(pwd)
	if err := config.Read(configFile); err == nil {
		return config
	}
	config.FromEnv(pwd, envs)
	return config
}

func parameterize(config *Configuration) (interface {
	internal_logger.Logger
	internal_logger.Printer
	internal.Configurer
}, []interface{}) {
	var employeesMeta interface {
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
		meta.Employee
	}
	var changesClient interface {
		changesclient.Client
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}
	var changesHandler interface {
		changesclient.Handler
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}
	var parameters []interface{}

	logger := internal_logger.New()
	switch v := config.MetaType; v {
	case internal_meta.TypeMemory:
		employeesMeta = metamemory.New()
	case internal_meta.TypeFile:
		employeesMeta = metafile.New()
	case internal_meta.TypeMySQL:
		employeesMeta = metamysql.New()
	}
	employeesMeta.SetUtilities(logger)
	switch {
	default:
		c := changesclientrest.New()
		changesClient, changesHandler = c, c
	case config.ClientChangesKafkaEnabled:
		switch {
		default:
			changesClient = changesclientrest.New()
		}
		changesHandler = changesclientkafka.New()
	}
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	employeesLogic := logic.New()
	employeesLogic.SetUtilities(logger)
	employeesLogic.SetParameters(employeesMeta, changesClient, changesHandler)
	parameters = append(parameters, employeesMeta, employeesLogic, changesClient, changesHandler)
	if config.ServiceRestEnabled {
		employeesRestService := servicerest.New()
		employeesRestService.SetUtilities(logger)
		employeesRestService.SetParameters(employeesLogic)
		restServer := internal_server_rest.New()
		restServer.SetUtilities(logger)
		restServer.SetParameters(employeesRestService)
		parameters = append(parameters, restServer, employeesRestService)
	}
	if config.ServiceGrpcEnabled {
		employeesGrpcService := servicegrpc.New()
		employeesGrpcService.SetUtilities(logger)
		employeesGrpcService.SetParameters(employeesLogic)
		grpcServer := internal_server_grpc.New()
		grpcServer.SetUtilities(logger)
		grpcServer.SetParameters(employeesGrpcService)
		parameters = append(parameters, grpcServer, employeesGrpcService)
	}
	return logger, parameters
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
