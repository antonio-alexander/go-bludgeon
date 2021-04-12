package rest

import (
	"os"
	"path/filepath"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/server/config"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"

	"github.com/pkg/errors"
)

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	var configFile = filepath.Join(pwd, config.DefaultConfigPath, config.DefaultConfigFile)
	var chExternal <-chan struct{}
	var meta interface {
		common.MetaTimer
		common.MetaTimeSlice
		common.MetaOwner
	}

	logger := logger.New("bludgeon-server")
	config := config.New()
	config.Default(pwd)
	if err = config.Read(configFile); err != nil {
		logger.Error(errors.Wrap(err, "error attempting to read config"))
		config.FromEnv(pwd, envs)
		if err := config.Write(configFile); err != nil {
			logger.Error(errors.Wrap(err, "error attempting to write config"))
		} else {
			logger.Info("Created config file using environment")
		}
	}
	switch config.MetaType {
	case common.MetaTypeJSON:
		meta = metajson.NewMetaJSON()
		if err = meta.Initialize(*config.MetaJSON); err != nil {
			return
		}
	case common.MetaTypeMySQL:
		meta = metamysql.NewMetaMySQL()
		if err = meta.Initialize(*config.MetaMySQL); err != nil {
			return
		}
	default:
		err = errors.Errorf("unsupported meta: %s", config.MetaType)

		return
	}
	server := New(logger, meta)
	if chExternal, err = server.Start(*config.RemoteRest); err != nil {
		return
	}
	<-chExternal
	err = server.Stop()

	return
}
