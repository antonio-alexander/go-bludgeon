package main

import (
	"os"
	"os/signal"
	"strings"

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
	//create channel to listen for ctrl+c
	chSignalInt := make(chan os.Signal, 1)
	signal.Notify(chSignalInt, os.Interrupt)
	//execute the client main for cli
	if err := server.MainRest(pwd, args, envs, chSignalInt); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	//close signal
	close(chSignalInt)
}
