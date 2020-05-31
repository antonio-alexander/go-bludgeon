package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

//Configuration
type Configuration struct {
	Meta struct {
		Type string `json:"Meta"`
		JSON struct {
			File string `json:"File"`
		} `json:"JSON"`
		MySQL mysql.Configuration `json:"Mysql"`
	}
	Remote struct {
		Type       string `json:"Type"`
		RestClient struct {
			Address string        `json:"Address"`
			Port    string        `json:"Port"`
			Timeout time.Duration `json:"Timeout"`
		} `json:"RestClient"`
	} `json:"Remote"`
	Client client.Configuration `json:"Client"`
}

//default constants
const (
//
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
	// json
	c.Meta.JSON.File = ""
	// mysql
	c.Meta.MySQL.Driver = "mysql"
	// c.Meta.MySQL.DataSource = ""
	// c.Meta.MySQL.FilePath = ""
	c.Meta.MySQL.Hostname = "127.0.0.1"
	c.Meta.MySQL.Port = "3306"
	c.Meta.MySQL.Username = "bludgeon"
	c.Meta.MySQL.Password = "bludgeon"
	c.Meta.MySQL.Database = "bludgeon"
	c.Meta.MySQL.ParseTime = false
	c.Meta.MySQL.UseTransactions = true
	c.Meta.MySQL.Timeout = 10 * time.Second
	// remote
	c.Remote.Type = "rest"
	c.Remote.RestClient.Address = "127.0.0.1"
	c.Remote.RestClient.Port = "8080"
	c.Remote.RestClient.Timeout = 10 * time.Second

	return
}
