package client

import (
	"fmt"
	"strings"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/client"
	clientcli "github.com/antonio-alexander/go-bludgeon/client/cli"
	clientrest "github.com/antonio-alexander/go-bludgeon/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/data"
	logic "github.com/antonio-alexander/go-bludgeon/logic"

	"github.com/pkg/errors"
)

func Main(pwd string, args []string, envs map[string]string) (err error) {
	var options clientcli.Options
	var cacheFile string
	var cache Cache
	var remote logic.Logic

	//
	config := NewConfiguration()
	if _, cacheFile, err = Files(pwd, &config); err != nil {
		return
	}
	if options, err = clientcli.Parse(pwd, args, envs); err != nil {
		return
	}
	if cache, err = CacheRead(cacheFile); err != nil {
		return
	}
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	if err = config.Read("", pwd, envs); err != nil {
		return
	}
	switch config.Client.Type {
	default:
		return errors.Errorf("Unsupported client: %s", config.Client.Type)
	case client.TypeRest:
		remote = clientrest.New(config.Client.Rest.Address, config.Client.Rest.Port)
	}
	switch options.ObjectType {
	default:
		return errors.Errorf("Unsupported object: %s", options.ObjectType)
	case data.ObjectTypeTimer:
		var timer data.Timer

		switch strings.ToLower(options.Command) {
		default:
			err = errors.Errorf("Unsupported command: %s", options.Command)
		case "create":
			timer, err = remote.TimerCreate()
		case "read":
			timer, err = remote.TimerRead(options.Timer.UUID)
		case "update":
			timer, err = remote.TimerUpdate(options.Timer)
		case "delete":
			err = remote.TimerDelete(options.Timer.UUID)
		case "start":
			timer, err = remote.TimerStart(options.Timer.UUID, time.Unix(0, options.Timer.Start))
		case "pause":
			//REVIEW: i don't think this works
			timer, err = remote.TimerPause(options.Timer.UUID, time.Unix(0, options.Timer.Finish))
		case "submit":
			timer, err = remote.TimerSubmit(options.Timer.UUID, time.Unix(0, options.Timer.Finish))
		}
		if err == nil {
			switch strings.ToLower(options.Command) {
			default:
				fmt.Println(timer)
			case "delete":
			}
		}
	}

	return
}
