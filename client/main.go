package client

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	api "github.com/antonio-alexander/go-bludgeon/client/api/rest"
	cli "github.com/antonio-alexander/go-bludgeon/client/cli"
	common "github.com/antonio-alexander/go-bludgeon/common"
	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"

	"github.com/pkg/errors"
)

func Main(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var config = &Configuration{}
	var cacheFile string
	var cache Cache
	var meta interface {
		common.MetaOwner
		common.MetaTimer
	}
	var remote interface {
		common.FunctionalOwner
		common.FunctionalTimer
		common.FunctionalTimeSlice
	}

	//
	if _, cacheFile, err = Files(pwd, &config); err != nil {
		return
	}
	if options, err = cli.Parse(pwd, args, envs); err != nil {
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
	switch config.RemoteType {
	case common.RemoteTypeRest:
		remote = api.NewFunctional()
		//TODO: verify that remote is legit
		if err = remote.Initialize(config.Remote[common.RemoteTypeRest]); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported remote: %s", config.RemoteType)

		return
	}
	switch config.MetaType {
	case common.MetaTypeJSON:
		meta = metajson.NewMetaJSON()
		//TODO: verify that meta is legit
		if err = meta.Initialize(config.Meta[common.MetaTypeJSON]); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported meta: %s", config.MetaType)

		return
	}
	c := NewClient(log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0), meta, remote)
	switch options.ObjectType {
	case common.ObjectTypeTimer:
		var timer common.Timer

		switch strings.ToLower(options.Command) {
		case "create":
			timer, err = c.TimerCreate()
		case "read":
			timer, err = c.TimerRead(options.Timer.UUID)
		case "update":
			timer, err = c.TimerUpdate(options.Timer)
		case "delete":
			err = c.TimerDelete(options.Timer.UUID)
		case "start":
			timer, err = c.TimerStart(options.Timer.UUID, time.Unix(0, options.Timer.Start))
		case "pause":
			//REVIEW: i don't think this works
			timer, err = c.TimerPause(options.Timer.UUID, time.Unix(0, options.Timer.Finish))
		case "submit":
			timer, err = c.TimerSubmit(options.Timer.UUID, time.Unix(0, options.Timer.Finish))
		default:
			err = errors.Errorf("Unsupported command: %s", options.Command)
		}
		if err == nil {
			switch strings.ToLower(options.Command) {
			case "delete":
			default:
				fmt.Println(timer)
			}
		}
	default:
		err = errors.Errorf("Unsupported object: %s", options.ObjectType)
	}
	if err := meta.Shutdown(); err != nil {
		fmt.Println(errors.Wrap(err, "Meta Shutdown"))
	}
	if err := remote.Shutdown(); err != nil {
		fmt.Println(errors.Wrap(err, "Remote Shutdown"))
	}
	c.Close()

	return
}
