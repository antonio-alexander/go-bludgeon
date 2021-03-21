package bludgeonmetamysql

import (
	"encoding/json"
	"fmt"
	"time"
)

//configuration defaults
const (
	DefaultHostname       string        = "127.0.0.1"
	DefaultPort           string        = "3306"
	DefaultUsername       string        = "bludgeon"
	DefaultPassword       string        = "bludgeon"
	DefaultDatabase       string        = "bludgeon"
	DefaultConnectTimeout time.Duration = 5 * time.Second
	DefaultQueryTimeout   time.Duration = 10 * time.Second
	DefaultParseTime      bool          = true
	DefaultDriver         string        = "mysql"
)

func castConfiguration(element interface{}) (c Configuration, err error) {
	switch v := element.(type) {
	case json.RawMessage:
		err = json.Unmarshal(v, &c)
	case *Configuration:
		c = *v
	case Configuration:
		c = v
	default:
		err = fmt.Errorf("Unsupported type: %t", element)
	}

	return
}

//Configuration is a struct that contains al lthe possible configuration for supported database drivers
type Configuration struct {
	Hostname       string        `json:"Hostname"`       //hostame to user to access the database
	Port           string        `json:"Port"`           //port to connect to
	Username       string        `json:"Username"`       //username to authenticate with
	Password       string        `json:"Password"`       //password to authenticate with
	Database       string        `json:"Database"`       //database to connect to
	ConnectTimeout time.Duration `json:"ConnectTimeout"` //how long to wait to connect
	QueryTimeout   time.Duration `json:"QueryTimeout"`   //how long to wait when querying
	ParseTime      bool          `json:"ParseTime"`      //whether or not to parse time
	// Driver         string        `json:"Driver"`         //what driver to use
}

//Validate is used to ensure that the values being configured make sense
// it's not necessarily to prevent a misconfiguration, but to use default values in the
// event a value doesn't exist
func (c *Configuration) Validate() (err error) {
	if c.Port == "" {
		c.Port = DefaultPort
	}
	if c.Database == "" {
		c.Database = DefaultDatabase
	}
	if c.Username == "" {
		c.Username = DefaultUsername
	}
	if c.Password == "" {
		c.Password = DefaultPassword
	}
	if c.Hostname == "" {
		c.Hostname = DefaultHostname
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = DefaultConnectTimeout
	}

	return
}

func (c *Configuration) Default() {
	c.Hostname = DefaultHostname
	c.Port = DefaultPort
	c.Username = DefaultUsername
	c.Password = DefaultPassword
	c.Database = DefaultDatabase
	c.ConnectTimeout = DefaultConnectTimeout
	c.QueryTimeout = DefaultQueryTimeout
	c.ParseTime = DefaultParseTime
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) (err error) {
	//TODO: add functionality

	return
}
