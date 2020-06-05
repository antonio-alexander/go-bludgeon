package main

import (
	"fmt"
	"os"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
	bludgeonclient "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
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
	Main(pwd, args, envs)
}

func Main(pwd string, args []string, envs map[string]string) (err error) {
	var options cli.Options
	var config Configuration
	var command bludgeon.CommandClient
	var data interface{}
	var meta bludgeon.Meta
	var remote bludgeon.Remote
	var metaFx, remoteFx, clientFx func()
	var client bludgeonclient.Functional
	// var cache Cache

	//parse the arguments
	if options, err = cli.ParseClient(pwd, args, envs); err != nil {
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

func optionsToCommand(options cli.Options) (command bludgeon.CommandClient, data interface{}, err error) {
	//convert command string
	if command = bludgeon.AtoCommandClient(options.Command); command == bludgeon.CommandClientInvalid {
		err = fmt.Errorf("Invalid command: %s", options.Command)

		return
	}
	//swtich on command and populate data
	switch command {
	case bludgeon.CommandClientTimerCreate:
		//do nothing
	case bludgeon.CommandClientTimerRead, bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerStop:
		//inject timer id
		data = options.Timer.UUID
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
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerSubmit:
		if timer, ok := data.(bludgeon.Timer); ok {
			fmt.Printf("Timer:\n%s\n", timer)
		}
	default:
		//REVIEW: should we provide positive feedback?
	}

	return
}
