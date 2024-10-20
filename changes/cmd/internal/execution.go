package internal

import (
	"path/filepath"
	"strconv"

	logic "github.com/antonio-alexander/go-bludgeon/changes/logic"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	metamemory "github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	metamysql "github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	servicekafka "github.com/antonio-alexander/go-bludgeon/changes/service/kafka"
	servicerest "github.com/antonio-alexander/go-bludgeon/changes/service/rest"

	healthcheckrestservice "github.com/antonio-alexander/go-bludgeon/healthcheck/service/rest"

	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/pkg/meta"
	internal_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/server"
)

func getConfig(pwd string, envs map[string]string) interface{} {
	configFile := filepath.Join(pwd, internal_config.DefaultConfigPath, internal_config.DefaultConfigFile)
	config := new(Configuration)
	config.Default(pwd)
	if err := config.Read(configFile); err == nil {
		return config
	}
	return internal_config.Envs(envs)
}

func parameterize(config interface{}) (interface {
	internal_logger.Logger
	internal_logger.Printer
	common.Configurer
}, []interface{}) {
	var restEnabled, kafkaEnabled bool
	var metaType internal_meta.Type
	var parameters []interface{}
	var changesMeta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		common.Initializer
		common.Configurer
		common.Parameterizer
	}

	switch v := config.(type) {
	case internal_config.Envs:
		restEnabled, _ = strconv.ParseBool(v[EnvNameServiceRestEnabled])
		kafkaEnabled, _ = strconv.ParseBool(v[EnvNameServiceKafkaEnabled])
		metaType = internal_meta.Type(v[EnvNameMetaType])
	case *Configuration:
		restEnabled = v.RestEnabled
		kafkaEnabled = v.KafkaEnabled
		metaType = v.MetaType
	case Configuration:
		restEnabled = v.RestEnabled
		kafkaEnabled = v.KafkaEnabled
		metaType = v.MetaType
	}
	logger := internal_logger.New()
	switch metaType {
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
	if restEnabled {
		changesRestService := servicerest.New()
		changesRestService.SetUtilities(logger)
		changesRestService.SetParameters(changesLogic)
		healthCheckRestService := healthcheckrestservice.New()
		healthCheckRestService.SetUtilities(logger)
		healthCheckRestService.SetParameters(changesLogic)
		restServer := internal_rest.New()
		restServer.SetUtilities(logger)
		restServer.SetParameters(changesRestService, healthCheckRestService)
		parameters = append(parameters, restServer, changesRestService,
			healthCheckRestService)
	}
	if kafkaEnabled {
		kafkaClient := internal_kafka.New()
		kafkaClient.SetUtilities(logger)
		changesKafkaService := servicekafka.New()
		changesKafkaService.SetUtilities(logger)
		changesKafkaService.SetParameters(changesLogic, kafkaClient)
		parameters = append(parameters, kafkaClient,
			changesKafkaService)
	}
	return logger, parameters
}

func configure(pwd string, config interface{}, parameters ...interface{}) error {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case common.Configurer:
			switch v := config.(type) {
			default:
				if err := p.Configure(v); err != nil {
					return err
				}
			case map[string]string:
				if err := p.Configure(internal_config.Envs(v)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func initialize(parameters ...interface{}) error {
	for _, parameter := range parameters {
		if p, ok := parameter.(common.Initializer); ok {
			if err := p.Initialize(); err != nil {
				return err
			}
		}
	}
	return nil
}

func shutdown(parameters ...interface{}) {
	for _, parameter := range parameters {
		if p, ok := parameter.(common.Initializer); ok {
			p.Shutdown()
		}
	}
}
