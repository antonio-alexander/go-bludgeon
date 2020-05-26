package bludgeonserver

import (
	"log"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	config "github.com/antonio-alexander/go-bludgeon/bludgeon/config/server"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest_server"
	endpoints "github.com/antonio-alexander/go-bludgeon/bludgeon/server/endpoints"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/functional"
)

func MainRest(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	var conf config.Configuration
	var m interface {
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
		bludgeon.MetaOwner
	}
	var logOut, logError = log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0)
	var chExternal <-chan struct{}

	//get config from environment
	if conf, err = config.FromEnv(pwd, envs); err != nil {
		return
	}
	//initialize meta
	if m, err = initMeta(conf.Meta.Type, conf.Meta.Config[conf.Meta.Type]); err != nil {
		return
	}
	//create rest server
	r := rest.NewServer()
	//create server
	s := server.NewServer(logOut, logError, m)
	//start the server
	if chExternal, err = s.Start(conf.Server); err == nil {
		//build rest routes for server
		if err = r.BuildRoutes(endpoints.BuildRoutes(s)); err == nil {
			//start the rest server
			if err = r.Start(conf.Rest.Address, conf.Rest.Port); err == nil {
				//block until the signal is killed
				select {
				case <-chSignalInt:
					logOut.Println("Server stopped externally (os interrupt)")
					//stop server
					if err := s.Stop(); err != nil {
						logOut.Println(err)
					}
				case <-chExternal:
					//server was stopped internally
				}
				//stop rest server
				if err := r.Stop(); err != nil {
					logOut.Println(err)
				}
			}
		}
	}
	//close the server and clean-up resources
	r.Close()
	s.Close()

	return
}
