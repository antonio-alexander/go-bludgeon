package bludgeoncli

import (
	"flag"
	"path/filepath"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func ParseClient(pwd string, args []string, envs map[string]string) (o Options, err error) {
	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.command, ArgCommand, DefaultCommand, UsageCommand)
	flagSet.StringVar(&o.objectType, ArgType, DefaultType, UsageType)
	flagSet.StringVar(&o.Configuration, ArgConfiguration, filepath.Join(pwd, bludgeon.DefaultConfigurationFile), UsageConfiguration)
	//timer
	flagSet.StringVar(&o.Timer.UUID, ArgTimerID, DefaultTimerID, UsageTimerID)
	flagSet.Int64Var(&o.Timer.Start, ArgTimerStart, DefaultTimerStart, UsageTimerStart)
	flagSet.Int64Var(&o.Timer.Finish, ArgTimerFinish, DefaultTimerFinish, UsageTimerFinish)
	// flagSet.StringVar(&o.Timer.Comment, ArgTimerComment, DefaultTimerComment, UsageTimerComment)
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
		case "start", "resume":
			o.Command = bludgeon.CommandClientTimerStart.String()
		case "pause":
			o.Command = bludgeon.CommandClientTimerPause.String()
		case "read":
			o.Command = bludgeon.CommandClientTimerRead.String()
		case "stop":
			o.Command = bludgeon.CommandClientTimerStop.String()
		case "submit":
			o.Command = bludgeon.CommandClientTimerSubmit.String()
		default:
			//TODO: generate error
		}
	default:
		//TODO: generate error
	}

	return
}

func ParseServer(pwd string, args []string, envs map[string]string) (o Options, err error) {
	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.objectType, ArgType, DefaultType, UsageType)
	flagSet.StringVar(&o.Configuration, ArgConfiguration, filepath.Join(pwd, bludgeon.DefaultConfigurationFile), UsageConfiguration)
	//parse the arguments
	if err = flagSet.Parse(args); err != nil {
		return
	}

	return
}
