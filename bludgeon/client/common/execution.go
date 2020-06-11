package bludgeonclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	mjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/server/api"
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
	//meta
	c.Meta.Type = "json"
	c.Meta.Config = map[string]interface{}{
		"json": mjson.Configuration{
			File: "./data/bludgeon.json",
		},
		"mysql": mysql.Configuration{
			Hostname:        "127.0.0.1",
			Port:            "3306",
			Username:        "bludgeon",
			Password:        "bludgeon",
			Database:        "bludgeon",
			ParseTime:       false,
			UseTransactions: true,
			Timeout:         10 * time.Second,
			DataSource:      "",
			FilePath:        "",
		},
	}
	c.Remote.Type = "rest"
	c.Remote.Config = map[string]interface{}{
		"rest": rest.Configuration{
			Address: "127.0.0.1",
			Port:    "8080",
			Timeout: 10 * time.Second,
		},
	}

	return
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
