package bludgeonconfigclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
