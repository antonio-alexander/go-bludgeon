package bludgeonclientcli

import (
	"flag"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
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
	o.ObjectType = bludgeon.AtoObjectType(objectType)
	o.MetaType = bludgeon.AtoMetaType(metaType)
	o.RemoteType = bludgeon.AtoRemoteType(remoteType)

	return
}
