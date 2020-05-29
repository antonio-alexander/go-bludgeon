package main

import (
	"os"
	"strings"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
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
	client.MainCLI(pwd, args, envs)
}
