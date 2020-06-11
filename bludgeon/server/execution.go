package bludgeonserver

import (
	"fmt"
	"strings"

	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func initMeta(metaType string, config interface{}) (meta interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
	bludgeon.MetaOwner
}, err error) {

	switch strings.ToLower(metaType) {
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
		err = fmt.Errorf("Unsupported meta: %s", metaType)
	}

	return
}
