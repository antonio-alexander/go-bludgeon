package bludgeonclient

import (
	"fmt"
	"path/filepath"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
)

func MainCLI(pwd string, args []string, envs map[string]string) {
	var command bludgeon.CommandClient
	var options cli.Options
	var data interface{}
	var err error
	//TODO: generate code to pull this from a json file
	var configuration = Configuration{
		ServerAddress: "127.0.0.1",
		ServerPort:    "60000",
		ClientAddress: "127.0.0.1",
		ClientPort:    "60000",
		// Task:          0,
		// Employee:      0,
	}

	//defer function to print errors
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
	//parse the arguments
	if options, err = cli.ParseClient(pwd, args, envs); err != nil {
		return
	}
	//create metajson
	metaJSON := json.NewMetaJSON()
	//defer close function
	defer metaJSON.Close()
	//initialize metajson
	if err = metaJSON.Initialize(filepath.Join(pwd, "bludgeon.json")); err != nil {
		return
	}
	//convert options to command/data
	if command, data, err = optionsToCommand(options); err != nil {
		return
	}
	//create the client
	client := NewClient(metaJSON, nil)
	//defer function to close client
	defer client.Close()
	//TODO: de-serialize client information
	//TODO: defer serialize client information
	//start the client
	if err = client.Start(configuration); err != nil {
		return
	}
	//defer function to stop client
	defer func() {
		if err := client.Stop(); err != nil {
			fmt.Printf("%s\n", err)
		}
	}()
	//call command handler
	if data, err = client.CommandHandler(command, data); err != nil {
		return
	}
	//handle response
	if err = handleClientResponse(command, data); err != nil {
		return
	}

	return
}
