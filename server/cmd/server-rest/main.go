package main

import (
	"os"
	"os/signal"
	"strings"

	rest "github.com/antonio-alexander/go-bludgeon/server/rest"
)

func main() {
	//Get the present working directory, the args
	// and then grab the environment, create a signal
	// channel and look for ctrl+C and start the Main()
	pwd, _ := os.Getwd()
	args := os.Args[1:]
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = s[1]
		}
	}
	chSignalInt := make(chan os.Signal, 1)
	signal.Notify(chSignalInt, os.Interrupt)
	//execute the client main for cli
	if err := rest.Main(pwd, args, envs, chSignalInt); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	close(chSignalInt)
}
