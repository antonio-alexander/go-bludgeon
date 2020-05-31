package bludgeonclient_test

import (
	"testing"
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
