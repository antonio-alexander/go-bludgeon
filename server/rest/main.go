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
	var configPath = filepath.Join(pwd, config.DefaultConfigPath)
	var chExternal <-chan struct{}
	var meta interface {
		common.MetaTimer
		common.MetaTimeSlice
		common.MetaOwner
	}

	logger := logger.New()
	config := config.New()
	config.Default(pwd)
	if err = config.Read(configPath); err != nil {
		logger.Error(err)
		logger.Info("Attempting to read config from env")
		config.FromEnv(pwd, envs)
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
