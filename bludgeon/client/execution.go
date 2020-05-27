package bludgeonclient

import (
	"fmt"

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
