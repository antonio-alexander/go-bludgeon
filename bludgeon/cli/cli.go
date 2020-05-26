package bludgeoncli

import (
	"flag"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func ParseClient(pwd string, args []string, envs map[string]string) (o Options, err error) {
	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.command, ArgCommand, DefaultCommand, UsageCommand)
	flagSet.StringVar(&o.objectType, ArgType, DefaultType, UsageType)
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
	switch strings.ToLower(o.objectType) {
	case "timer":
		switch strings.ToLower(o.command) {
		case "create":
			o.Command = bludgeon.CommandClientTimerCreate.String()
		case "delete":
			o.Command = bludgeon.CommandClientTimerDelete.String()
		case "update":
			o.Command = bludgeon.CommandClientTimerUpdate.String()
		case "read":
			o.Command = bludgeon.CommandClientTimerRead.String()
		default:
			//TODO: generate error
		}
	default:
		//TODO: generate error
	}

	return
}
