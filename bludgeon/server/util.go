package bludgeonserver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func ConfigRead(file string) (config Configuration, err error) {
	var bytes []byte

	//check if file exists
	if _, err = os.Stat(file); os.IsNotExist(err) {
		//load default configuration
		config = ConfigDefault()
		//write the default configuration
		err = ConfigWrite(file, config)
	} else {
		//file exists
		if bytes, err = ioutil.ReadFile(file); err != nil {
			return
		}
		err = json.Unmarshal(bytes, &config)
	}
	return
}

func ConfigWrite(file string, config Configuration) (err error) {
	var bytes []byte

	//marshal config into bytes
	if bytes, err = json.MarshalIndent(&config, "", "    "); err != nil {
		return
	}
	//write configuration
	err = ioutil.WriteFile(file, bytes, 0644)

	return
}

func ConfigDefault() (c Configuration) {
	c.Server.TokenWait = time.Duration(30) * time.Minute
	c.Server.Rest.Address = ""
	c.Server.Rest.Port = "8080"

	return
}
