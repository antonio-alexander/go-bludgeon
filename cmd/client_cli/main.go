package main

import (
	"fmt"
	"os"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
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

func Main(pwd string, args []string, envs map[string]string) {
	//parse the arguments
	options, err := cli.ParseClient(pwd, args, envs)
	if err != nil {
		fmt.Println(err)
		return
	}
	//read configuration
	config, err := ConfigRead(options.Configuration)
	if err != nil {
		fmt.Println(err)
		return
	}
	//convert options to command/data
	command, data, err := optionsToCommand(options)
	if err != nil {
		fmt.Println(err)
		return
	}
	//initialize meta
	meta, metaFx, err := initMeta(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer metaFx()
	//initialize remote
	remote, remoteFx, err := initRemote(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer remoteFx()
	//initialize client
	client, clientFx, err := initClient(config, meta, remote)
	if err != nil {
		fmt.Println(err)
	} else {
		defer clientFx()
		//call command handler
		data, err := client.CommandHandler(command, data)
		if err != nil {
			fmt.Println(err)
		} else {
			//handle response
			if err = handleClientResponse(command, data); err != nil {
				fmt.Println(err)
			}
		}
	}

	return
}

func optionsToCommand(options cli.Options) (command bludgeon.CommandClient, data interface{}, err error) {
	//convert command string
	if command = bludgeon.AtoCommandClient(options.Command); command == bludgeon.CommandClientInvalid {
		//TODO: generate error
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
		//TODO: generate error
	}

	return
}

func handleClientResponse(command bludgeon.CommandClient, data interface{}) (err error) {
	//swtich on command and populate data
	switch command {
	case bludgeon.CommandClientTimerRead, bludgeon.CommandClientTimerCreate:
		if timer, ok := data.(bludgeon.Timer); ok {
			fmt.Printf("Timer:\n%s\n", timer)
		}
	case bludgeon.CommandClientTimerStart:
	case bludgeon.CommandClientTimerSubmit:
	case bludgeon.CommandClientTimerPause:
	case bludgeon.CommandClientTimerStop:
	default:
		//TODO: generate error
	}

	return
}
