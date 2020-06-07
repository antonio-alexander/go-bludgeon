package bludgeonserver

import (
	"fmt"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/server"
)

func MainRest(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) (err error) {
	var options Options
	var config Configuration
	var fxMetaCleanup func()
	var m interface {
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
	}
	//parse command line arguments
	if options, err = Parse(pwd, args, envs); err != nil {
		return
	}
	//configuration
	if config, err = ConfigRead(options.Configuration); err != nil {
		return
	}
	//initialize meta
	if m, fxMetaCleanup, err = initMeta(config); err != nil {
		return
	}
	//defer meta cleanup
	defer fxMetaCleanup()
	//create rest server
	r := rest.NewServer()
	//create server
	s := NewServer(m)
	//start the server
	if err = s.Start(config); err == nil {
		//build rest routes for server
		routes := s.BuildRoutes()
		if err = r.BuildRoutes(routes); err == nil {
			//start the rest server
			if err = r.Start(config.Server.Rest.Address, config.Server.Rest.Port); err == nil {
				//block until the signal is killed
				select {
				case <-chSignalInt:
					fmt.Println("Server stopped externally (os interrupt)")
				}
				//stop rest server
				if err := r.Stop(); err != nil {
					fmt.Println(err)
				}
				//stop server
				if err := s.Stop(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	//close the server and clean-up resources
	r.Close()
	s.Close()

	return
}
