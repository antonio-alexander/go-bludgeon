package rest

import (
	"os"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/server/config"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"

	"github.com/pkg/errors"
)

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	// var configPath string
	var config config.Configurationo
	var meta interface {
		common.MetaTimer
		common.MetaTimeSlice
		common.MetaOwner
	}
	var logger = logger.New()
	var chExternal <-chan struct{}

	// //
	// if configPath, _, err = common.Files(pwd, &config); err != nil {
	// 	return
	// }
	// if err = config.Read(configPath, pwd, envs, &config); err != nil {
	// 	return
	// }
	// if err = config.Write(configPath, &config); err != nil {
	// 	return
	// }
	switch config.MetaType {
	case common.MetaTypeJSON:
		meta = metajson.NewMetaJSON()
		if err = meta.Initialize(config.Meta.JSON); err != nil {
			return
		}
	case common.MetaTypeMySQL:
		meta = metamysql.NewMetaMySQL()
		if err = meta.Initialize(config.Meta.MySQL); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported meta: %s", config.MetaType)

		return
	}
	server := New(logger, meta)
	if chExternal, err = server.Start(config.Remote.Rest); err != nil {
		return
	}
	<-chExternal
	err = server.Stop()

	return
}
