package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Cache struct {
	TimerID string `json:"TimerID"`
}

func CacheRead(file string) (cache Cache, err error) {
	var bytes []byte

	//check if file exists
	if _, err = os.Stat(file); os.IsNotExist(err) {
		//load default configuration
		cache = CacheDefault()
		//write the default configuration
		err = CacheWrite(file, cache)
	} else {
		//file exists
		if bytes, err = ioutil.ReadFile(file); err != nil {
			return
		}
		err = json.Unmarshal(bytes, &cache)
	}
	return
}

func CacheWrite(file string, cache Cache) (err error) {
	var bytes []byte

	//marshal config into bytes
	if bytes, err = json.MarshalIndent(&cache, "", "    "); err != nil {
		return
	}
	//write configuration
	err = ioutil.WriteFile(file, bytes, 0644)

	return
}

func CacheDefault() (c Cache) {
	c.TimerID = ""

	return
}
