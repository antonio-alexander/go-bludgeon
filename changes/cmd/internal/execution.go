package internal

import (
	"path/filepath"

	logic "github.com/antonio-alexander/go-bludgeon/changes/logic"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	metamemory "github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	metamysql "github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	servicekafka "github.com/antonio-alexander/go-bludgeon/changes/service/kafka"
	servicerest "github.com/antonio-alexander/go-bludgeon/changes/service/rest"

	healthcheckrestservice "github.com/antonio-alexander/go-bludgeon/healthcheck/service/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	serverrest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
)

func getConfig(pwd string, envs map[string]string) *Configuration {
	configFile := filepath.Join(pwd, config.DefaultConfigPath, config.DefaultConfigFile)
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
	var parameters []interface{}
	var changesMeta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}

	logger := internal_logger.New()
	switch v := config.MetaType; v {
	case internal_meta.TypeMemory:
		changesMeta = metamemory.New()
	case internal_meta.TypeFile:
		changesMeta = metafile.New()
	case internal_meta.TypeMySQL:
		changesMeta = metamysql.New()
	}
	changesMeta.SetUtilities(logger)
	changesLogic := logic.New()
	changesLogic.SetUtilities(logger)
	changesLogic.SetParameters(changesMeta)
	parameters = append(parameters, changesMeta, changesLogic)
	if config.RestEnabled {
		changesRestService := servicerest.New()
		changesRestService.SetUtilities(logger)
		changesRestService.SetParameters(changesLogic)
		healthCheckRestService := healthcheckrestservice.New()
		healthCheckRestService.SetUtilities(logger)
		healthCheckRestService.SetParameters(changesLogic)
		restServer := serverrest.New()
		restServer.SetUtilities(logger)
		restServer.SetParameters(changesRestService, healthCheckRestService)
		parameters = append(parameters, changesMeta, restServer,
			changesRestService, healthCheckRestService)
	}
	if config.KafkaEnabled {
		kafkaClient := kafka.New()
		kafkaClient.SetUtilities(logger)
		changesKafkaService := servicekafka.New()
		changesKafkaService.SetUtilities(logger)
		changesKafkaService.SetParameters(changesLogic, kafkaClient)
		parameters = append(parameters, changesMeta, kafkaClient, changesKafkaService)
	}
	return logger, parameters
}

func configure(pwd string, envs map[string]string, parameters ...interface{}) error {
	//TODO: allow this to be able to accept configuration from a json
	// file
	for _, p := range parameters {
		switch p := p.(type) {
		case internal.Configurer:
			if err := p.Configure(config.Envs(envs)); err != nil {
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
