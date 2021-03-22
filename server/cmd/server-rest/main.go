package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	config "github.com/antonio-alexander/go-bludgeon/bludgeon/config"
	metajson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
	endpoints "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/endpoints"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/server"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server"

	"github.com/pkg/errors"
)

func main() {
	//get environment
	pwd, _ := os.Getwd()
	args := os.Args[1:]
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = s[1]
		}
	}
	//create channel to listen for ctrl+c
	chSignalInt := make(chan os.Signal, 1)
	signal.Notify(chSignalInt, os.Interrupt)
	//execute the client main for cli
	if err := Main(pwd, args, envs, chSignalInt); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	//close signal
	close(chSignalInt)
}

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	var configPath string
	var conf config.Server
	var meta interface {
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
		bludgeon.MetaOwner
	}
	var logOut, logError = log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0)
	var chExternal <-chan struct{}

	//
	if configPath, _, err = bludgeon.Files(pwd, &conf); err != nil {
		return
	}
	if err = config.Read(configPath, pwd, envs, &conf); err != nil {
		return
	}
	if err = config.Write(configPath, &conf); err != nil {
		return
	}
	switch conf.MetaType {
	case bludgeon.MetaTypeJSON:
		meta = metajson.NewMetaJSON()
		if err = meta.Initialize(conf.Meta[bludgeon.MetaTypeJSON]); err != nil {
			return
		}
	case bludgeon.MetaTypeMySQL:
		meta = metamysql.NewMetaMySQL()
		if err = meta.Initialize(conf.Meta[bludgeon.MetaTypeMySQL]); err != nil {
			return
		}
	default:
		err = errors.Errorf("Unsupported meta: %s", conf.MetaType)

		return
	}
	r, s := rest.NewServer(), server.NewServer(logOut, logError, meta)
	if chExternal, err = s.Start(conf.Server); err == nil {
		routes := endpoints.BuildRoutes(nil, s)
		if err = r.BuildRoutes(routes); err == nil {
			if err = r.Start(conf.Rest.Address, conf.Rest.Port); err == nil {
				logOut.Printf("Rest Server started on %s:%s", conf.Rest.Address, conf.Rest.Port)
				select {
				case <-chSignalInt:
					logOut.Println("Server stopped externally (os interrupt)")
					if err := s.Stop(); err != nil {
						logOut.Println(err)
					}
				case <-chExternal:
					logOut.Println("Server stopped internally")
				}
			}
		}
	}
	r.Close()
	s.Close()

	return
}
