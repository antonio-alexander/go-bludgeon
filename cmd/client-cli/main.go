package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	config "github.com/antonio-alexander/go-bludgeon/bludgeon/config"
	api "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/api"
	client "github.com/antonio-alexander/go-bludgeon/internal/client"
	cli "github.com/antonio-alexander/go-bludgeon/internal/client/cli"
	bludgeon "github.com/antonio-alexander/go-bludgeon/internal/common"
	json "github.com/antonio-alexander/go-bludgeon/internal/meta/json"

	"github.com/pkg/errors"
)

func main() {
	//get environment
	pwd, _ := os.Getwd()
	args := os.Args[1:]
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = s[1]
		}
	}
	//execute the client main for cli
	if err := Main(pwd, args, envs); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func Main(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var conf config.Client
	var cacheFile string
	var cache client.Cache
	var meta interface {
		bludgeon.MetaOwner
		bludgeon.MetaTimer
	}
	var remote interface {
		bludgeon.FunctionalOwner
		bludgeon.FunctionalTimer
		bludgeon.FunctionalTimeSlice
	}

	//
	if _, cacheFile, err = bludgeon.Files(pwd, &conf); err != nil {
		return
	}
	if options, err = cli.Parse(pwd, args, envs); err != nil {
		return
	}
	if cache, err = client.CacheRead(cacheFile); err != nil {
		return
	}
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	if err = config.Read("", pwd, envs, &conf); err != nil {
		return
	}
	switch conf.RemoteType {
	case bludgeon.RemoteTypeRest:
		remote = api.NewFunctional()
		//TODO: verify that remote is legit
		if err = remote.Initialize(conf.Remote[bludgeon.RemoteTypeRest]); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported remote: %s", conf.RemoteType)

		return
	}
	switch conf.MetaType {
	case bludgeon.MetaTypeJSON:
		meta = json.NewMetaJSON()
		//TODO: verify that meta is legit
		if err = meta.Initialize(conf.Meta[bludgeon.MetaTypeJSON]); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported meta: %s", conf.MetaType)

		return
	}
	c := client.NewClient(log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0), meta, remote)
	switch options.ObjectType {
	case bludgeon.ObjectTypeTimer:
		var timer bludgeon.Timer

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
