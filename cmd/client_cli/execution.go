package main

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/remote/rest"
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

func initRemote(config Configuration) (interface {
	bludgeon.Remote
}, func(), error) {
	//switch on rest type
	switch config.Remote.Type {
	case "rest":
		//create rest remote
		r := rest.NewRemote(config.Remote.RestClient.Address, config.Remote.RestClient.Port, config.Remote.RestClient.Timeout)
		//create defer func
		deferFx := func() {
			//
		}

		return r, deferFx, nil
	default:
		return nil, func() {}, nil
	}
}

func initClient(config Configuration, m interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, r interface {
	bludgeon.Remote
}) (interface {
	client.Functional
}, func(), error) {

	//create the client
	c := client.NewClient(m, r)
	//TODO: de-serialize client information
	//TODO: defer serialize client information
	//start the client
	if err := c.Start(config.Client); err != nil {
		return nil, nil, err
	}
	//create defer function
	deferFx := func() {
		//stop the client
		if err := c.Stop(); err != nil {
			fmt.Println(err)
		}
		//close the client
		c.Close()
	}

	return c, deferFx, nil
}
