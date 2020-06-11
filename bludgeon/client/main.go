package bludgeonclient

import (
	"fmt"
	"path/filepath"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/client/cli"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/client/common"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/functional"
)

func MainCli(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var config common.Configuration
	var command bludgeon.CommandClient
	var data interface{}
	var meta interface {
		bludgeon.MetaOwner
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
	}
	var remote interface {
		bludgeon.RemoteOwner
		bludgeon.RemoteTimer
		bludgeon.RemoteTimeSlice
	}
	var cache common.Cache

	//parse the arguments
	if options, err = cli.Parse(pwd, args, envs); err != nil {
		return
	}
	//read cache
	if cache, err = common.CacheRead(filepath.Join(pwd, "data/bludgeon_cache.json")); err != nil {
		return
	}
	//if timer is empty, replace with data from cache
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	//read configuration
	if config, err = common.ConfigRead(options.Configuration); err != nil {
		return
	}
	//initialize meta
	if meta, err = initMeta(config.Meta.Type, nil); err != nil {
		return
	}
	//initialize remote
	if remote, err = initRemote(config.Remote.Type, nil); err != nil {
		return
	}
	//convert options to command/data
	if command, data, err = cli.OptionsToCommand(options); err != nil {
		return
	}
	c := client.NewClient(nil, nil, meta, remote)
	//call command handler
	if data, err = c.CommandHandler(command, data); err == nil {
		//handle response
		err = cli.HandleClientResponse(command, data)
	}
	//
	if err := meta.Shutdown(); err != nil {
		fmt.Println(err)
	}
	//
	if err := remote.Shutdown(); err != nil {
		fmt.Println(err)
	}
	//close the client
	c.Close()

	return
}
