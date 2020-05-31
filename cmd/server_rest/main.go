package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/server"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server"
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
	//execute the client main for cli
	Main(pwd, args, envs)
}

func Main(pwd string, args []string, envs map[string]string) {
	//parse command line arguments
	options, err := cli.ParseServer(pwd, args, envs)
	if err != nil {
		fmt.Println(err)
		return
	}
	//configuration
	config, err := ConfigRead(options.Configuration)
	if err != nil {
		fmt.Println(err)
		return
	}
	//initialize meta
	m, metaFx, err := initMeta(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer meta cleanup
	defer metaFx()
	//create rest server
	r := rest.NewServer()
	//create server
	s := server.NewServer(m)
	//start the server
	if err = s.Start(server.Configuration{
		TokenWait: config.Server.TokenWait,
	}); err == nil {
		//build rest routes for server
		if err = r.BuildRoutes(s); err == nil {
			//start the rest server
			if err = r.Start(config.RestServer.Address, config.RestServer.Port); err == nil {
				//create channel to listen for ctrl+c
				chSignalInt := make(chan os.Signal, 1)
				signal.Notify(chSignalInt, os.Interrupt)
				//block until the signal is killed
				select {
				case <-chSignalInt:
					fmt.Println("Server stopped externally (os interrupt)")
				}
				close(chSignalInt)
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
	//perform error checking
	if err != nil {
		fmt.Println(err)
	}
	//close the server and clean-up resources
	r.Close()
	s.Close()
}
