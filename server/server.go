package server

import (
	"os"
	"path/filepath"

	common "github.com/antonio-alexander/go-bludgeon/common"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
	rest "github.com/antonio-alexander/go-bludgeon/server/rest"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"

	"github.com/pkg/errors"
)

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (errs []error) {
	var configFile = filepath.Join(pwd, DefaultConfigPath, DefaultConfigFile)

	logger := logger.New("bludgeon-server")
	config := NewConfiguration()
	config.Default(pwd)
	if err := config.Read(configFile); err != nil {
		logger.Error(errors.Wrap(err, "error attempting to read config"))
		config.FromEnv(pwd, envs)
		if err := config.Write(configFile); err != nil {
			logger.Error(errors.Wrap(err, "error attempting to write config"))
		} else {
			logger.Info("Created config file using environment")
		}
	}
	meta, err := getMeta(config)
	if err != nil {
		errs = append(errs, err)
		return
	}
	logic := logic.New(logger, meta)
	serverRest := rest.New(logger, logic)
	if err := serverRest.Start(*config.Server.Rest); err != nil {
		errs = append(errs, err)
		return
	}
	<-chSignalInt
	if err := serverRest.Stop(); err != nil {
		errs = append(errs, err)
	}

	return
}

func getMeta(config *Configuration) (interface {
	common.MetaTimer
	common.MetaTimeSlice
	common.MetaOwner
}, error) {
	switch config.MetaType {
	default:
		return nil, errors.Errorf("unsupported meta: %s", config.MetaType)
	case common.MetaTypeJSON:
		meta := metajson.NewMetaJSON()
		if err := meta.Initialize(*config.Meta.JSON); err != nil {
			return nil, err
		}
		return meta, nil
	case common.MetaTypeMySQL:
		meta := metamysql.NewMetaMySQL()
		if err := meta.Initialize(*config.Meta.MySQL); err != nil {
			return nil, err
		}
		return meta, nil
	}
}
