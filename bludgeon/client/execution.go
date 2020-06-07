package bludgeonclient

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	bjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	bmysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
	brest "github.com/antonio-alexander/go-bludgeon/bludgeon/server/api"
)

func initMeta(config Configuration) (interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, func(), error) {
	// filepath.Join(pwd, "bludgeon.json")
	switch config.Meta.Type {
	case "json":
		//create metajson
		m := bjson.NewMetaJSON()
		//initialize metajson
		if err := m.Initialize(config.Meta.JSON.File); err != nil {
			return nil, nil, err
		}
		deferFx := func() {
			m.Close()
		}

		return m, deferFx, nil
	case "mysql":
		m := bmysql.NewMetaMySQL()
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
		r := brest.NewRemote(config.Remote.RestClient.Address, config.Remote.RestClient.Port, config.Remote.RestClient.Timeout)
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
	Functional
}, func(), error) {

	//create the client
	c := NewClient(m, r)
	//TODO: de-serialize client information
	//start the client
	if err := c.Start(config); err != nil {
		return nil, nil, err
	}
	//create defer function
	deferFx := func() {
		//stop the client
		if err := c.Stop(); err != nil {
			fmt.Println(err)
		}
		//TODO: defer serialize client information
		//close the client
		c.Close()
	}

	return c, deferFx, nil
}

func optionsToCommand(options Options) (command bludgeon.CommandClient, data interface{}, err error) {
	//swtich on command and populate data
	switch command = options.Command; command {
	case bludgeon.CommandClientTimerCreate:
		//do nothing
	case bludgeon.CommandClientTimerRead, bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerStop:
		//inject timer id
		data = options.Timer.UUID
	case bludgeon.CommandClientTimerUpdate:
		data = options.Timer
	default:
		err = fmt.Errorf("Unsupported command: %s", command.String())
	}

	return
}

func handleClientResponse(command bludgeon.CommandClient, data interface{}) (err error) {
	//swtich on command and populate data
	switch command {
	case bludgeon.CommandClientTimerCreate, bludgeon.CommandClientTimerRead,
		bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerStop,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerUpdate:
		if timer, ok := data.(bludgeon.Timer); ok {
			fmt.Printf("Timer:\n%s\n", timer)
		}
	default:
		//REVIEW: should we provide positive feedback?
	}

	return
}
