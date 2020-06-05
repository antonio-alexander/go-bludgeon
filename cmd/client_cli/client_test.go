package main_test

import (
	"testing"

	client "github.com/antonio-alexander/go-bludgeon/cmd/client_cli"
)

func TestMainCli(t *testing.T) {
	//create cases
	cases := map[string]struct {
		pwd  string
		args []string
		envs map[string]string
	}{
		"": {
			pwd: "/Users/noobius/source_control/go/src/github.com/antonio-alexander/go-bludgeon/cmd/client_cli",
			args: []string{
				"-command",
				"start",
				"-type",
				"timer",
				"-id",
				"59c20b2d-e1a5-42a0-bb98-628d6000d47c",
			},
		},
	}
	//range over cases
	for _, c := range cases {
		if err := client.Main(c.pwd, c.args, c.envs); err != nil {
			t.Fatal(err)
		}
	}
}
