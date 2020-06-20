package bludgeonclientcli

import (
	"flag"

	client "github.com/antonio-alexander/go-bludgeon/client"
	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
)

func Parse(pwd string, args []string, envs map[string]string) (o Options, err error) {
	var objectType, remoteType, metaType string

	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon-client", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.Command, ArgCommand, DefaultCommand, UsageCommand)
	flagSet.StringVar(&objectType, ArgObjectType, DefaultObjectType, UsageObjectType)
	flagSet.StringVar(&remoteType, ArgClientType, DefaultClientType, UsageClientType)
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
	o.ObjectType = data.AtoObjectType(objectType)
	o.MetaType = meta.AtoType(metaType)
	o.ClientType = client.AtoType(remoteType)

	return
}
