package grpcclient

import (
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ErrAddressEmpty string = "address is empty"
	ErrPortEmpty    string = "port is empty"
	ErrPortBadf     string = "port is a non-integer: %s"
)

const (
	EnvNameAddress string = "BLUDGEON_GRPC_ADDRESS"
	EnvNamePort    string = "BLUDGEON_GRPC_PORT"
)

const (
	DefaultAddress string = "127.0.0.1"
	DefaultPort    string = "8081"
)

var (
	DefaultOptions = []grpc.DialOption{
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	}
)

type Configuration struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Options []grpc.DialOption
}

func (c *Configuration) Default() {
	c.Address = DefaultAddress
	c.Port = DefaultPort
	c.Options = DefaultOptions
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameAddress]; ok {
		c.Address = address
	}
	if port, ok := envs[EnvNamePort]; ok {
		c.Port = port
	}
}

func (c *Configuration) Validate() (err error) {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is lte 0
	if c.Address == "" {
		return errors.New(ErrAddressEmpty)
	}
	if c.Port == "" {
		return errors.New(ErrPortEmpty)
	}
	if _, e := strconv.Atoi(c.Port); e != nil {
		return errors.Errorf(ErrPortBadf, c.Port)
	}
	return
}
