package bludgeonclient_test

import (
	"fmt"
	"testing"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
)

//test serialization
//test deserialization

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
				"read",
				"-type",
				"timer",
				"-id",
				"768e8c52-23a6-4146-a65c-446d4f340eff",
			},
		},
	}
	//range over cases
	for _, c := range cases {
		client.MainCLI(c.inPwd, c.inArgs, c.inEnvs)
	}
}

func TestTimerCreate(t *testing.T) {
	//create cache interface
	client := client.NewClient(nil, nil)
	//create a timer
	if timer, err := client.TimerCreate(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", timer)
	}
	//attempt to serialize the data in the cache
	if bytes, err := client.Serialize(); err != nil {
		fmt.Println(bytes)
	} else {
		fmt.Println(bytes)
	}
	//close the cache
	client.Close()
}

func TestTimerRead(t *testing.T) {

	cases := map[string]struct {
		oErr string //error
	}{
		//check for non existent timer
		//check for existent timer
	}
	//range over cases
	for range cases {
		//
	}
}

func TestTimerDelete(t *testing.T) {
	//
}

func TestTimerStart(t *testing.T) {
	//
}

func TestTimerPause(t *testing.T) {
	//
}

func TestTimerSubmit(t *testing.T) {
	//
}
