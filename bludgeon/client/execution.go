package bludgeonclient

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
)

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
	case bludgeon.CommandClientTimerRead:
		//inject timer id
		data = options.Timer.UUID
	default:
		//TODO: generate error
	}

	return
}
