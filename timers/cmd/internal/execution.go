package internal

import (
	"path/filepath"

	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	metamemory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	metamysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	servicegrpc "github.com/antonio-alexander/go-bludgeon/timers/service/grpc"
	servicerest "github.com/antonio-alexander/go-bludgeon/timers/service/rest"

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
	var timersMeta interface {
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
		meta.Timer
		meta.TimeSlice
	}
	var parameters []interface{}
	var changesClient interface {
		internal.Parameterizer
		internal.Initializer
		internal.Configurer
		changesclient.Client
	}
	var changesHandler interface {
		internal.Parameterizer
		internal.Initializer
		internal.Configurer
		changesclient.Handler
	}

	logger := internal_logger.New()
	switch v := config.MetaType; v {
	case internal_meta.TypeMemory:
		timersMeta = metamemory.New()
	case internal_meta.TypeFile:
		timersMeta = metafile.New()
	case internal_meta.TypeMySQL:
		timersMeta = metamysql.New()
	}
	switch {
	default:
		c := changesclientrest.New()
		changesHandler, changesClient = c, c
	case config.ClientChangesKafkaEnabled:
		switch {
		default:
			changesClient = changesclientrest.New()
		}
		changesHandler = changesclientkafka.New()
	}
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	timersMeta.SetUtilities(logger)
	timersLogic := logic.New()
	timersLogic.SetUtilities(logger)
	timersLogic.SetParameters(timersMeta, changesClient, changesHandler)
	parameters = append(parameters, timersMeta, timersLogic, changesClient, changesHandler)
	if config.ServiceRestEnabled {
		restServer := internal_server_rest.New()
		timersRestService := servicerest.New()
		timersRestService.SetUtilities(logger)
		restServer.SetUtilities(logger)
		restServer.SetParameters(timersRestService)
		timersRestService.SetParameters(restServer, timersLogic)
		parameters = append(parameters, restServer, timersRestService)
	}
	if config.ServiceGrpcEnabled {
		timersGrpcService := servicegrpc.New()
		timersGrpcService.SetUtilities(logger)
		timersGrpcService.SetParameters(timersLogic, timersGrpcService)
		grpcServer := internal_server_grpc.New()
		grpcServer.SetUtilities(logger)
		grpcServer.SetParameters(timersGrpcService)
		parameters = append(parameters, grpcServer, timersGrpcService)
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
