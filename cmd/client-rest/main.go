package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	config "github.com/antonio-alexander/go-bludgeon/bludgeon/config"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	api "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/api"
	endpoints "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/endpoints"
	restserver "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/server"
)

func main() {
	//
	pwd, _ := os.Getwd()
	args := os.Args[1:]
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = s[1]
		}
	}
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	if err := Main(pwd, args, envs, osSignal); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func Main(pwd string, args []string, envs map[string]string, osSignal chan os.Signal) (err error) {
	// var options cli.Options
	var conf config.Client
	var configFile string

	//
	if configFile, _, err = bludgeon.Files(pwd, &conf); err != nil {
		return
	}
	// if options, err = cli.Parse(pwd, args, envs); err != nil {
	// 	return
	// }
	if err = config.Read(configFile, pwd, envs, &conf); err != nil {
		return
	}
	remote := api.NewFunctional()
	//TODO: verify that remote is legit
	if err = remote.Initialize(conf.Remote[bludgeon.RemoteTypeRest]); err != nil {
		return
	}
	meta := json.NewMetaJSON()
	//TODO: verify that meta is legit
	if err = meta.Initialize(conf.Meta[bludgeon.MetaTypeJSON]); err != nil {
		return
	}
	c := client.NewClient(log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0), meta, remote)
	router := restserver.NewServer()
	routes := endpoints.BuildRoutes(nil, c)
	if err = router.BuildRoutes(routes); err == nil {
		if err = router.Start(conf.Rest.Address, conf.Rest.Port); err == nil {
			<-osSignal
			if err := router.Stop(); err != nil {
				fmt.Println(err)
			}
		}
	}
	if err := meta.Shutdown(); err != nil {
		fmt.Println(err)
	}
	if err := remote.Shutdown(); err != nil {
		fmt.Println(err)
	}
	c.Close()

	return
}
