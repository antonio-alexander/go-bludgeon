package bludgeonserver

import (
	"flag"
	"path/filepath"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func Parse(pwd string, args []string, envs map[string]string) (o Options, err error) {
	//create the flagset
	flagSet := flag.NewFlagSet("bludgeon", flag.ExitOnError)
	//options
	flagSet.StringVar(&o.Configuration, ArgConfiguration, filepath.Join(pwd, bludgeon.DefaultConfigurationFile), UsageConfiguration)
	//parse the arguments
	if err = flagSet.Parse(args); err != nil {
		return
	}

	return
}
