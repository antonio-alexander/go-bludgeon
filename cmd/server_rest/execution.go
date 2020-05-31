package main

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

func initMeta(config Configuration) (interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, func(), error) {
	// filepath.Join(pwd, "bludgeon.json")
	switch config.Meta.Type {
	case "json":
		//create metajson
		m := json.NewMetaJSON()
		//initialize metajson
		if err := m.Initialize(config.Meta.JSON.File); err != nil {
			return nil, nil, err
		}
		deferFx := func() {
			m.Close()
		}

		return m, deferFx, nil
	case "mysql":
		m := mysql.NewMetaMySQL()
		//connect
		if err := m.Connect(config.Meta.MySQL); err != nil {
			return nil, nil, err
		}
		//create defer function
		deferFx := func() {
			//disconnect
			if err := m.Disconnect(); err != nil {
				fmt.Println(err)
			}
			//close
			m.Close()
		}

		return m, deferFx, nil
	default:
		return nil, func() {}, nil
	}
}
