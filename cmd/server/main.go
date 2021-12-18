package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/antonio-alexander/go-bludgeon/cmd/server/internal"
)

func main() {
	pwd, _ := os.Getwd()
	args := os.Args[1:]
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	chSignalInt := make(chan os.Signal, 1)
	signal.Notify(chSignalInt, os.Interrupt)
	if err := internal.Main(pwd, args, envs, chSignalInt); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	close(chSignalInt)
}
