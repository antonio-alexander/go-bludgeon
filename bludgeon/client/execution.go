package bludgeonclient

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/server/api"
)

func initMeta(metaType string, config interface{}) (meta interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, err error) {
	// filepath.Join(pwd, "bludgeon.json")
	switch metaType {
	case "json":
		//create metajson
		m := json.NewMetaJSON()
		//initialize metajson
		if err = m.Initialize(config); err != nil {
			return
		}
		meta = m
	case "mysql":
		m := mysql.NewMetaMySQL()
		//connect
		if err = m.Initialize(config); err != nil {
			return
		}
		meta = m
	default:
		err = fmt.Errorf("meta type unsupported: %s", metaType)
	}

	return
}

func initRemote(remoteType string, config interface{}) (remote interface {
	bludgeon.RemoteOwner
	bludgeon.RemoteTimer
	bludgeon.RemoteTimeSlice
}, err error) {
	//switch on rest type
	switch remoteType {
	case "rest":
		//create rest remote
		r := rest.NewRemote()
		if err = r.Initialize(config); err != nil {
			return
		}
		remote = r
	default:
		err = fmt.Errorf("remote type unsupported: %s", remoteType)
	}

	return
}
