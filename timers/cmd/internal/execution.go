package internal

import (
	"errors"
	"path/filepath"

	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	servergrpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	serverrest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	metatimers "github.com/antonio-alexander/go-bludgeon/timers/meta"
	metafile "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	metamemory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	metamysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	servicegrpc "github.com/antonio-alexander/go-bludgeon/timers/service/grpc"
	servicerest "github.com/antonio-alexander/go-bludgeon/timers/service/rest"
)

func getConfig(pwd string, envs map[string]string) *config.Configuration {
	c := config.NewConfiguration()
	c.Default(pwd)
	configFile := filepath.Join(pwd, config.DefaultConfigPath, config.DefaultConfigFile)
	if err := c.Read(configFile); err == nil {
		return c
	}
	c.FromEnv(pwd, envs)
	return c
}

func getLogger(config *config.Configuration) logger.Logger {
	config.Logger.Prefix = serviceName
	logger := logger.New(config.Logger)
	return logger
}

func startMeta(config *config.Configuration, parameters ...interface{}) (interface {
	metatimers.Timer
	metatimers.TimeSlice
	meta.Owner
}, error) {
	switch v := config.Meta.Type; v {
	default:
		return nil, meta.ErrUnsupportedMeta(v)
	case meta.TypeMemory:
		return metamemory.New(parameters...), nil
	case meta.TypeFile:
		meta := metafile.New(parameters...)
		if err := meta.Initialize(config.Meta.File); err != nil {
			return nil, err
		}
		return meta, nil
	case meta.TypeMySQL:
		meta := metamysql.New(parameters...)
		if err := meta.Initialize(config.Meta.Mysql); err != nil {
			return nil, err
		}
		return meta, nil
	}
}

func startLogic(config *config.Configuration, parameters ...interface{}) (interface {
	logic.Logic
}, error) {
	logic := logic.New(parameters...)
	return logic, nil
}

func startServices(config *config.Configuration, parameters ...interface{}) (func(), error) {
	if !config.Server.RestEnabled && !config.Server.GrpcEnabled {
		return nil, errors.New("no servers enabled")
	}
	restServer := serverrest.New(parameters...)
	if config.Server.RestEnabled {
		servicerest.New(append(parameters, restServer)...)
		if err := restServer.Start(config.Server.Rest); err != nil {
			return nil, err
		}
	}
	grpcServer := servergrpc.New(parameters...)
	if config.Server.GrpcEnabled {
		grpcService := servicegrpc.New(append(parameters, grpcServer)...)
		if err := grpcServer.Initialize(config.Server.Grpc,
			grpcService.Register); err != nil {
			return nil, err
		}
	}
	return func() {
		if config.Server.RestEnabled {
			restServer.Stop()
		}
		if config.Server.GrpcEnabled {
			grpcServer.Shutdown()
		}
	}, nil
}
