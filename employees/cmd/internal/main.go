package internal

import (
	"os"
	"path/filepath"

	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	service "github.com/antonio-alexander/go-bludgeon/employees/service"

	meta_file "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	meta_mysql "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	service_rest "github.com/antonio-alexander/go-bludgeon/employees/service/rest"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	server_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/pkg/errors"
)

func startMeta(config *Configuration, parameters ...interface{}) (interface {
	meta.Employee
	meta.Owner
}, error) {
	switch v := config.Meta.Type; v {
	default:
		return nil, errors.Errorf("unsupported meta: %s", v)
	case meta.TypeFile:
		meta := meta_file.New(parameters...)
		if err := meta.Initialize(config.Meta.File); err != nil {
			return nil, err
		}
		return meta, nil
	case meta.TypeMySQL:
		meta := meta_mysql.New(parameters...)
		if err := meta.Initialize(config.Meta.Mysql); err != nil {
			return nil, err
		}
		return meta, nil
	}
}

func startLogic(config *Configuration, parameters ...interface{}) (interface {
	logic.Logic
}, error) {
	logic := logic.New(parameters...)
	return logic, nil
}

func startService(config *Configuration, parameters ...interface{}) (server_rest.Owner, error) {
	server := server_rest.New(parameters...)
	switch v := config.Server.Type; v {
	default:
		return nil, errors.Errorf("unsupported service: %s", v)
	case service.TypeREST:
		parameters = append(parameters, server)
		service_rest.New(parameters...)
	}
	if err := server.Start(config.Server.Rest); err != nil {
		return nil, err
	}
	return server, nil
}

func getConfig(pwd string, envs map[string]string, logger logger.Logger) *Configuration {
	config := NewConfiguration()
	config.Default(pwd)
	configFile := filepath.Join(pwd, DefaultConfigPath, DefaultConfigFile)
	if err := config.Read(configFile); err == nil {
		return config
	}
	logger.Info("using config from environment")
	config.FromEnv(pwd, envs)
	return config
}

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	logger := logger.New("bludgeon-employees-service")
	logger.Info("version %s (%s@%s)", Version, GitBranch, GitCommit)
	config := getConfig(pwd, envs, logger)
	meta, err := startMeta(config, logger)
	if err != nil {
		return err
	}
	defer meta.Shutdown()
	logic, err := startLogic(config, logger, meta)
	if err != nil {
		return err
	}
	service, err := startService(config, logic, logger)
	if err != nil {
		return err
	}
	defer service.Stop()
	<-chSignalInt
	return nil
}
