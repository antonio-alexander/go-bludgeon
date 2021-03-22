package bludgeonclientcli

import (
	"flag"

	common "github.com/antonio-alexander/go-bludgeon/common"
)

func Parse(pwd string, args []string, envs map[string]string) (o Options, err error) {
	var objectType, remoteType, metaType string

	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.Command, ArgCommand, DefaultCommand, UsageCommand)
	flagSet.StringVar(&objectType, ArgObjectType, DefaultObjectType, UsageObjectType)
	flagSet.StringVar(&remoteType, ArgRemoteType, DefaultRemoteType, UsageRemoteType)
	flagSet.StringVar(&metaType, ArgMetaType, DefaultMetaType, UsageMetaType)
	//timer
	flagSet.StringVar(&o.Timer.UUID, ArgTimerID, DefaultTimerID, UsageTimerID)
	flagSet.Int64Var(&o.Timer.Start, ArgTimerStart, DefaultTimerStart, UsageTimerStart)
	flagSet.Int64Var(&o.Timer.Finish, ArgTimerFinish, DefaultTimerFinish, UsageTimerFinish)
	flagSet.StringVar(&o.Timer.Comment, ArgTimerComment, DefaultTimerComment, UsageTimerComment)
	//parse the arguments
	if err = flagSet.Parse(args); err != nil {
		return
	}
	o.ObjectType = common.AtoObjectType(objectType)
	o.MetaType = common.AtoMetaType(metaType)
	o.RemoteType = common.AtoRemoteType(remoteType)

	return
}
