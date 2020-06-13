package bludgeonclient

import (
	"fmt"
	"log"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/client/cli"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/client/common"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/functional"
	config "github.com/antonio-alexander/go-bludgeon/bludgeon/config/client"
)

func MainCli(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var conf config.Configuration
	var command bludgeon.CommandClient
	var cacheFile, configFile, jsonFile string
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

	if configFile, cacheFile, jsonFile, err = initClientFolders(pwd); err != nil {
		return
	}
	//parse the arguments
	if options, err = cli.Parse(pwd, args, envs); err != nil {
		return
	}
	//read cache
	if cache, err = common.CacheRead(cacheFile); err != nil {
		return
	}
	//if timer is empty, replace with data from cache
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	//read configuration
	if conf, err = config.Read(configFile, jsonFile); err != nil {
		return
	}
	//initialize meta
	if meta, err = initMeta(conf.Meta.Type, conf.Meta.Config[conf.Meta.Type]); err != nil {
		return
	}
	//initialize remote
	if remote, err = initRemote(conf.Remote.Type, conf.Remote.Config[conf.Remote.Type]); err != nil {
		return
	}
	//convert options to command/data
	if command, data, err = cli.OptionsToCommand(options); err != nil {
		return
	}
	c := client.NewClient(log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0), meta, remote)
	//call command handler
	if data, err = c.CommandHandler(command, data); err == nil {
		//handle response
		err = cli.HandleClientResponse(command, data, cacheFile)
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

func MainRest(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var conf config.Configuration
	var command bludgeon.CommandClient
	var cacheFile, configFile, jsonFile string
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

	if configFile, cacheFile, jsonFile, err = initClientFolders(pwd); err != nil {
		return
	}
	//parse the arguments
	if options, err = cli.Parse(pwd, args, envs); err != nil {
		return
	}
	//read cache
	if cache, err = common.CacheRead(cacheFile); err != nil {
		return
	}
	//if timer is empty, replace with data from cache
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	//read configuration
	if conf, err = config.Read(configFile, jsonFile); err != nil {
		return
	}
	//initialize meta
	if meta, err = initMeta(conf.Meta.Type, conf.Meta.Config[conf.Meta.Type]); err != nil {
		return
	}
	//initialize remote
	if remote, err = initRemote(conf.Remote.Type, conf.Remote.Config[conf.Remote.Type]); err != nil {
		return
	}
	//convert options to command/data
	if command, data, err = cli.OptionsToCommand(options); err != nil {
		return
	}
	c := client.NewClient(log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0), meta, remote)
	//call command handler
	if data, err = c.CommandHandler(command, data); err == nil {
		//handle response
		err = cli.HandleClientResponse(command, data, cacheFile)
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
