package main_test

import (
	"testing"

	client "github.com/antonio-alexander/go-bludgeon/cmd/client_cli"
)

func TestMainCli(t *testing.T) {
	//get actual values
	// aPwd, _ := os.Getwd()
	// aArgs := os.Args[1:]
	// aEnvs := make(map[string]string)
	// for _, env := range os.Environ() {
	// 	if s := strings.Split(env, "="); len(s) > 1 {
	// 		aEnvs[s[0]] = s[1]
	// 	}
	// }
	//create cases
	cases := map[string]struct {
		inPwd  string
		inArgs []string
		inEnvs map[string]string
	}{
		"read": {
			inPwd: "/Users/noobius/source_control/go/src/github.com/antonio-alexander/go-bludgeon/cmd/client_cli",
			inArgs: []string{
				"-command",
				"create",
				"-type",
				"timer",
				"-id",
				"768e8c52-23a6-4146-a65c-446d4f340eff",
			},
		},
	}
	//range over cases
	for _, c := range cases {
		client.Main(c.inPwd, c.inArgs, c.inEnvs)
	}
}
