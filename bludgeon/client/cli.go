package bludgeonclient

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func Parse(pwd string, args []string, envs map[string]string) (o Options, err error) {
	var command, objectType string

	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&command, ArgCommand, DefaultCommand, UsageCommand)
	flagSet.StringVar(&objectType, ArgType, DefaultType, UsageType)
	flagSet.StringVar(&o.Configuration, ArgConfiguration, filepath.Join(pwd, bludgeon.DefaultConfigurationFile), UsageConfiguration)
	//timer
	flagSet.StringVar(&o.Timer.UUID, ArgTimerID, DefaultTimerID, UsageTimerID)
	flagSet.Int64Var(&o.Timer.Start, ArgTimerStart, DefaultTimerStart, UsageTimerStart)
	flagSet.Int64Var(&o.Timer.Finish, ArgTimerFinish, DefaultTimerFinish, UsageTimerFinish)
	flagSet.StringVar(&o.Timer.Comment, ArgTimerComment, DefaultTimerComment, UsageTimerComment)
	//parse the arguments
	if err = flagSet.Parse(args); err != nil {
		return
	}
	//after parsing, convert command to something appropriate
	switch strings.ToLower(objectType) {
	case "timer":
		switch strings.ToLower(command) {
		case "create":
			o.Command = bludgeon.CommandClientTimerCreate
		case "start", "resume":
			o.Command = bludgeon.CommandClientTimerStart
		case "pause":
			o.Command = bludgeon.CommandClientTimerPause
		case "read":
			o.Command = bludgeon.CommandClientTimerRead
		case "stop":
			o.Command = bludgeon.CommandClientTimerStop
		case "submit":
			o.Command = bludgeon.CommandClientTimerSubmit
		case "update":
			o.Command = bludgeon.CommandClientTimerUpdate
		default:
			err = fmt.Errorf("Command %s, not supported for object %s", command, objectType)
		}
	default:
		err = fmt.Errorf("No support for object %s", objectType)
	}

	return
}
