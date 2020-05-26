package bludgeonconfigserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	mjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
)

func Read(file, jsonFile string) (config Configuration, err error) {
	var bytes []byte

	//check if file exists
	if _, err = os.Stat(file); os.IsNotExist(err) {
		//load default configuration
		config = configDefault(jsonFile)
		//write the default configuration
		err = Write(file, config)
	} else {
		//file exists
		if bytes, err = ioutil.ReadFile(file); err != nil {
			return
		}
		err = json.Unmarshal(bytes, &config)
	}

	return
}

func Write(file string, config Configuration) (err error) {
	var bytes []byte

	//marshal config into bytes
	if bytes, err = json.MarshalIndent(&config, "", "    "); err != nil {
		return
	}
	//write configuration
	err = ioutil.WriteFile(file, bytes, 0644)

	return
}

func FromEnv(pwd string, envs map[string]string) (c Configuration, err error) {
	c = configDefault(pwd)
	//get meta type
	if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
		c.Meta.Type = metaType
	}
	//get rest address
	if restAddress, ok := envs[EnvNameBludgeonRestAddress]; ok {
		c.Rest.Address = restAddress
	}
	//get rest port
	if restPort, ok := envs[EnvNameBludgeonRestPort]; ok {
		c.Rest.Port = restPort
	}
	//
	for _, metaType := range []string{"json", "mysql"} {
		var bytes json.RawMessage

		//get meta configuration
		switch metaType {
		case "json":
			var config mjson.Configuration

			//get json file
			if jsonFile, ok := envs[EnvNameBludgeonMetaJSONFile]; ok {
				config.File = jsonFile
			}
			if bytes, err = json.Marshal(&config); err != nil {
				return
			}
			c.Meta.Config["json"] = bytes
		case "mysql":
			var config mysql.Configuration

			//get database address
			if databaseAddress, ok := envs[EnvNameDatabaseAddress]; ok {
				config.Hostname = databaseAddress
			}
			//get database port
			if databasePort, ok := envs[EnvNameDatabasePort]; ok {
				config.Port = databasePort
			}
			if bytes, err = json.Marshal(&config); err != nil {
				return
			}
			c.Meta.Config["mysql"] = bytes
		default:
			err = fmt.Errorf("Unsupported meta: \"%s\"", metaType)
		}
	}

	return
}
