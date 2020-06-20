package main

import (
	"os"
	"strings"

	internal "github.com/antonio-alexander/go-bludgeon/cmd/client/internal"
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
	if err := internal.Main(pwd, args, envs); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
