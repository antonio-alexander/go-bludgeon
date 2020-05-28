package bludgeonclient

import (
	"fmt"
	"path/filepath"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	cli "github.com/antonio-alexander/go-bludgeon/bludgeon/cli"
	meta "github.com/antonio-alexander/go-bludgeon/bludgeon/meta"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
	remote "github.com/antonio-alexander/go-bludgeon/bludgeon/remote"
)

func MainCLI(pwd string, args []string, envs map[string]string) {
	var command bludgeon.CommandClient
	var options cli.Options
	var data interface{}
	var err error
	var meta interface {
		meta.MetaTimer
		meta.MetaTimeSlice
		Close()
	}
	var client interface {
		Owner
		Functional
		Manage
		API
	}

	//TODO: generate code to pull this from a json file
	//parse the arguments
	if options, err = cli.ParseClient(pwd, args, envs); err != nil {
		fmt.Println(err)

		return
	}
	//convert options to command/data
	if command, data, err = optionsToCommand(options); err != nil {
		fmt.Println(err)

		return
	}
	//initialize meta
	if meta, err = metaInit(pwd); err != nil {
		fmt.Println(err)

		return
	}
	//initialize client
	if client, err = clientInit(Configuration{
		ServerAddress: "127.0.0.1",
		ServerPort:    "60000",
		ClientAddress: "127.0.0.1",
		ClientPort:    "60000",
		// Task:          0,
		// Employee:      0,
	}, meta, nil); err != nil {
		fmt.Println(err)
	} else {
		//call command handler
		if data, err = client.CommandHandler(command, data); err != nil {
			fmt.Println(err)
		} else {
			//handle response
			if err = handleClientResponse(command, data); err != nil {
				fmt.Println(err)
			}
		}
	}
	//shutdown client
	if err = clientShutdown(client); err != nil {
		fmt.Println(err)
	}
	//shutdown meta
	if err = metaShutdown(meta); err != nil {
		fmt.Println(err)
	}

	return
}

func clientInit(config Configuration, meta interface {
	meta.MetaTimer
	meta.MetaTimeSlice
}, remote interface {
	remote.Remote
}) (client interface {
	Owner
	Manage
	Functional
	API
}, err error) {

	//create the client
	client = NewClient(meta, remote)
	//TODO: de-serialize client information
	//TODO: defer serialize client information
	//start the client
	if err = client.Start(config); err != nil {
		return
	}

	return
}

func clientShutdown(client interface {
	Owner
	Manage
}) (err error) {

	//stop the client
	err = client.Stop()
	//close the client
	client.Close()

	return
}

const metaType string = "mysql"

func metaInit(pwd string) (meta interface {
	meta.MetaTimer
	meta.MetaTimeSlice
	Close()
}, err error) {

	switch metaType {
	case "json":
		//create metajson
		metaJSON := json.NewMetaJSON()
		//initialize metajson
		if err = metaJSON.Initialize(filepath.Join(pwd, "bludgeon.json")); err != nil {
			return
		}
		meta = metaJSON
	case "mysql":
		metaMySQL := mysql.NewMetaMySQL()
		//initialize
		if err = metaMySQL.Connect(mysql.Configuration{
			Driver:          "mysql",
			DataSource:      "bludgeon",
			Hostname:        "127.0.0.1",
			Port:            "306",
			Username:        "bludgeon",
			Password:        "bludgeon",
			Database:        "bludgeon",
			ParseTime:       true,
			UseTransactions: true,
			Timeout:         5 * time.Second,
		}); err != nil {
			return
		}
		//connect
		meta = metaMySQL
	default:
	}

	return
}

func metaShutdown(meta interface {
	Close()
}) (err error) {

	switch metaType {
	case "json":
		meta.Close()
	case "mysql":
		meta.Close()
	}

	return
}
