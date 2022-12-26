package grpcserver

import (
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	ErrPortEmpty string = "port is empty"
	ErrPortBadf  string = "port is a non-integer: %s"
)

const (
	EnvNameAddress string = "BLUDGEON_GRPC_ADDRESS"
	EnvNamePort    string = "BLUDGEON_GRPC_PORT"
)

const (
	DefaultAddress string = ""
	DefaultPort    string = "8081"
)

var (
// DefaultAllowedOrigins = [...]string{"http://localhost:8000"}
// DefaultAllowedMethods = [...]string{http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete}
)

type Configuration struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Options []grpc.ServerOption
}

func (c *Configuration) Default() {
	c.Address = DefaultAddress
	c.Port = DefaultPort
}

func (c *Configuration) FromEnv(envs map[string]string) {
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
	if c.Port == "" {
		return errors.New(ErrPortEmpty)
	}
	if _, e := strconv.Atoi(c.Port); e != nil {
		return errors.Errorf(ErrPortBadf, c.Port)
	}
	return
}
