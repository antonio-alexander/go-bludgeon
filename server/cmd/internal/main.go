package internal

import (
	"os"
	"path/filepath"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logic "github.com/antonio-alexander/go-bludgeon/internal/logic"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	server "github.com/antonio-alexander/go-bludgeon/server"

	logger_simple "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	logic_simple "github.com/antonio-alexander/go-bludgeon/internal/logic/simple"
	meta_file "github.com/antonio-alexander/go-bludgeon/meta/file"
	meta_mysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
	server_rest "github.com/antonio-alexander/go-bludgeon/server/rest"

	"github.com/pkg/errors"
)

func getMeta(config *Configuration) (interface {
	meta.Timer
	meta.TimeSlice
	meta.Owner
}, error) {
	switch v := config.Meta.Type; v {
	default:
		return nil, errors.Errorf("unsupported meta: %s", v)
	case meta.TypeFile:
		meta := meta_file.New()
		if err := meta.Initialize(config.Meta.File); err != nil {
			return nil, err
		}
		return meta, nil
	case meta.TypeMySQL:
		meta := meta_mysql.New()
		if err := meta.Initialize(config.Meta.Mysql); err != nil {
			return nil, err
		}
		return meta, nil
	}
}

func getLogic(config *Configuration, logger logger.Logger, meta interface {
	meta.Timer
	meta.TimeSlice
}) (interface {
	logic.Logic
	logic.Functional
}, error) {
	switch v := config.Server.Type; v {
	default:
		return nil, errors.Errorf("unsupported server: %s", v)
	case server.Type(logic.TypeSimple):
		logic := logic_simple.New(logger, meta)
		if err := logic.Start(); err != nil {
			return nil, err
		}
		return logic, nil
	}
}

func getServer(config *Configuration, logic logic.Logic, logger logger.Logger) (server.Owner, error) {
	switch v := config.Server.Type; v {
	default:
		return nil, errors.Errorf("unsupported server: %s", v)
	case server.TypeREST:
		server := server_rest.New(logger, logic)
		if err := server.Start(config.Server.Rest); err != nil {
			return nil, err
		}
		return server, nil
	}
}

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	logger := logger_simple.New("bludgeon-server")
	logger.Info("version %s (%s@%s)", Version, GitBranch, GitCommit)
	config := NewConfiguration()
	config.Default(pwd)
	configFile := filepath.Join(pwd, DefaultConfigPath, DefaultConfigFile)
	if err := config.Read(configFile); err != nil {
		logger.Error("failed to read config - %v", err)
		config.FromEnv(pwd, envs)
	}
	meta, err := getMeta(config)
	if err != nil {
		return err
	}
	logic, err := getLogic(config, logger, meta)
	if err != nil {
		return err
	}
	server, err := getServer(config, logic, logger)
	if err != nil {
		return err
	}
	<-chSignalInt
	server.Stop()
	return nil
}
