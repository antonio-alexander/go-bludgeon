package grpc

import (
	"strconv"

	internal_grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"

	"github.com/pkg/errors"
)

// environmental variables
const (
	EnvNameGrpcAddress string = "BLUDGEON_EMPLOYEES_GRPC_ADDRESS"
	EnvNameRestPort    string = "BLUDGEON_EMPLOYEES_GRPC_PORT"
)

// defaults
const (
	DefaultPort    string = "8011"
	DefaultAddress string = "localhost"
)

type Configuration struct {
	internal_grpc.Configuration
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
}

func (r *Configuration) FromEnv(envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameGrpcAddress]; ok {
		r.Address = address
	}
	if port, ok := envs[EnvNameRestPort]; ok {
		r.Port = port
	}
}

func (r *Configuration) Validate() error {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is gt 0
	if r.Address == "" {
		return errors.New(internal_grpc.ErrAddressEmpty)
	}
	if r.Port == "" {
		return errors.New(internal_grpc.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		return errors.Errorf(internal_grpc.ErrPortBadf, r.Port)
	}
	return nil
}
