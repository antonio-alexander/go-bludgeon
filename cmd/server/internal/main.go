package internal

import (
	"os"
	"path/filepath"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/meta/file"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
	server "github.com/antonio-alexander/go-bludgeon/server"
	serverrest "github.com/antonio-alexander/go-bludgeon/server/rest"

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
		meta := metafile.New()
		if err := meta.Initialize(config.Meta.File); err != nil {
			return nil, err
		}
		return meta, nil
	case meta.TypeMySQL:
		meta := metamysql.New()
		if err := meta.Initialize(config.Meta.Mysql); err != nil {
			return nil, err
		}
		return meta, nil
	}
}

func getServer(config *Configuration, logic logic.Logic, logger data.Logger) (server.Owner, error) {
	switch v := config.Server.Type; v {
	default:
		return nil, errors.Errorf("unsupported server: %s", v)
	case server.TypeREST:
		server := serverrest.New(logger, logic)
		if err := server.Start(config.Server.Rest); err != nil {
			return nil, err
		}
		return server, nil
	}
}

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	logger := logger.New("bludgeon-server")
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
	logic := logic.New(logger, meta)
	server, err := getServer(config, logic, logger)
	if err != nil {
		return err
	}
	<-chSignalInt
	server.Stop()
	return nil
}
