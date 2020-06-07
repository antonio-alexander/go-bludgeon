package bludgeonclient_test

import (
	"testing"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
)

//test serialization
//test deserialization

// func TestTimerCreate(t *testing.T) {
// 	//create cache interface
// 	client := client.NewClient(nil, nil)
// 	//create a timer
// 	if timer, err := client.TimerCreate(); err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Printf("%#v\n", timer)
// 	}
// 	//attempt to serialize the data in the cache
// 	if bytes, err := client.Serialize(); err != nil {
// 		fmt.Println(bytes)
// 	} else {
// 		fmt.Println(bytes)
// 	}
// 	//close the cache
// 	client.Close()
// }

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
				"--command", "create",
				"--type", "timer",
				"--id", "51db329b-9ab1-4636-840f-da06eb4acaea",
				"--comment=\"This is a better test\"",
			},
		},
	}
	//range over cases
	for _, c := range cases {
		if err := client.MainCli(c.pwd, c.args, c.envs); err != nil {
			t.Fatal(err)
		}
	}
}
