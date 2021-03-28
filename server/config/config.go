package config

import (
	"time"

	"github.com/antonio-alexander/go-bludgeon/common"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
)

type Configuration struct {
	RemoteType common.RemoteType
	Remote     struct {
		Rest Rest
	}
	MetaType common.MetaType
	Meta     struct {
		MySQL metamysql.Configuration
		JSON  metajson.Configuration
	}
}

type Rest struct {
	Address string        `json:"Address"`
	Port    string        `json:"Port"`
	Timeout time.Duration `json:"Timeout"`
}

// func Default() Configuration {
// 	return Configuration{
// 		Address: DefaultAddress,
// 		Port:    DefaultPort,
// 		Timeout: DefaultTimeout,
// 	}
// }

// func FromEnv(pwd string, envs map[string]string, c *Configuration) (err error) {
// 	//Get the address from the environment, then the port
// 	// then the timeout
// 	if address, ok := envs[EnvNameAddress]; ok {
// 		c.Address = address
// 	}
// 	if port, ok := envs[EnvNamePort]; ok {
// 		c.Port = port
// 	}
// 	if timeoutString, ok := envs[EnvNameTimeout]; ok {
// 		if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
// 			if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
// 				c.Timeout = timeout
// 			}
// 		}
// 	}

// 	return
// }

// func Validate(c *Configuration) (err error) {
// 	//validate that the address isn't empty
// 	// check if the port is empty, and then ensure
// 	// that the port is an integer, finally
// 	// check if the timeout is lte 0
// 	if c.Address == "" {
// 		err = errors.New(ErrAddressEmpty)

// 		return
// 	}
// 	if c.Port == "" {
// 		err = errors.New(ErrPortEmpty)

// 		return
// 	}
// 	if _, e := strconv.Atoi(c.Port); e != nil {
// 		err = errors.Errorf(ErrPortBadf, c.Port)

// 		return
// 	}
// 	if c.Timeout <= 0 {
// 		err = errors.Errorf(ErrTimeoutBadf, c.Timeout)
// 	}

// 	return
// }
