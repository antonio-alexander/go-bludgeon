package bludgeonclientcli

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func OptionsToCommand(options Options) (command bludgeon.CommandClient, data interface{}, err error) {
	//swtich on command and populate data
	switch command = options.Command; command {
	case bludgeon.CommandClientTimerCreate:
		//do nothing
	case bludgeon.CommandClientTimerRead, bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerStop:
		//inject timer id
		data = options.Timer.UUID
	case bludgeon.CommandClientTimerUpdate:
		data = options.Timer
	default:
		err = fmt.Errorf("Unsupported command: %s", command.String())
	}

	return
}

func HandleClientResponse(command bludgeon.CommandClient, data interface{}) (err error) {
	//swtich on command and populate data
	switch command {
	case bludgeon.CommandClientTimerCreate, bludgeon.CommandClientTimerRead,
		bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerStop,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerUpdate:
		if timer, ok := data.(bludgeon.Timer); ok {
			fmt.Printf("Timer:\n%s\n", timer)
		}
	default:
		//REVIEW: should we provide positive feedback?
	}

	return
}
