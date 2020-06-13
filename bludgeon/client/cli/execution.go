package bludgeonclientcli

import (
	"fmt"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/client/common"
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

func HandleClientResponse(command bludgeon.CommandClient, data interface{}, cacheFile string) (err error) {
	//read cache
	cache, _ := common.CacheRead(cacheFile) //swtich on command and populate data
	//
	switch command {
	case bludgeon.CommandClientTimerCreate, bludgeon.CommandClientTimerRead,
		bludgeon.CommandClientTimerStart, bludgeon.CommandClientTimerStop,
		bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerSubmit,
		bludgeon.CommandClientTimerUpdate:
		//cache the timer id
		//print the timer
		if timer, ok := data.(bludgeon.Timer); ok {
			fmt.Printf("Timer:\n%s\n", timer)
			cache.TimerID = timer.UUID
		}
		//write the cache
		err = common.CacheWrite(cacheFile, cache)
	default:
		//REVIEW: should we provide positive feedback?
	}

	return
}
