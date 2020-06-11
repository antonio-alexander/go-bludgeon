package bludgeonserver

import (
	"log"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest_server"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/server/common"
	endpoints "github.com/antonio-alexander/go-bludgeon/bludgeon/server/endpoints"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/functional"
)

func MainRest(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	var config common.Configuration
	var m interface {
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
		bludgeon.MetaOwner
	}
	var logOut, logError = log.New(os.Stdout, "", 0), log.New(os.Stderr, "", 0)
	var chExternal <-chan struct{}

	//get config from environment
	if config, err = common.GetConfigFromEnv(pwd, envs); err != nil {
		return
	}
	//initialize meta
	if m, err = initMeta(config.Meta.Type, config.Meta.Config[config.Meta.Type]); err != nil {
		return
	}
	//create rest server
	r := rest.NewServer()
	//create server
	s := server.NewServer(logOut, logError, m)
	//start the server
	if chExternal, err = s.Start(config); err == nil {
		//build rest routes for server
		if err = r.BuildRoutes(endpoints.BuildRoutes(s)); err == nil {
			//start the rest server
			if err = r.Start(config.Server.Rest.Address, config.Server.Rest.Port); err == nil {
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
