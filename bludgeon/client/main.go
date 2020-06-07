package bludgeonclient

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func MainCli(pwd string, args []string, envs map[string]string) (err error) {
	var options Options
	var config Configuration
	var command bludgeon.CommandClient
	var data interface{}
	var meta bludgeon.Meta
	var remote bludgeon.Remote
	var metaFx, remoteFx, clientFx func()
	var client Functional
	// var cache Cache

	//parse the arguments
	if options, err = Parse(pwd, args, envs); err != nil {
		return
	}
	// //read cache
	// if cache, err = CacheRead(filepath.Join(pwd, "data/bludgeon_cache.json")); err != nil {
	// 	return
	// }
	// //if timer is empty, replace with data from cache
	// if options.Timer.UUID == "" {
	// 	options.Timer.UUID = cache.TimerID
	// }
	//read configuration
	if config, err = ConfigRead(options.Configuration); err != nil {
		return
	}
	//initialize meta
	if meta, metaFx, err = initMeta(config); err != nil {
		return
	}
	defer metaFx()
	//initialize remote
	if remote, remoteFx, err = initRemote(config); err != nil {
		return
	}
	defer remoteFx()
	//initialize client
	if client, clientFx, err = initClient(config, meta, remote); err != nil {
		return
	}
	defer clientFx()
	//convert options to command/data
	if command, data, err = optionsToCommand(options); err != nil {
		return
	}
	//call command handler
	if data, err = client.CommandHandler(command, data); err != nil {
		return
	}
	//handle response
	if err = handleClientResponse(command, data); err != nil {
		return
	}
	//TODO: update cache

	return
}
