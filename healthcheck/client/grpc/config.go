package grpc

import (
	"strconv"

	pkg_grpc "github.com/antonio-alexander/go-bludgeon/pkg/grpc/client"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// environmental variables
const (
	EnvNameGrpcAddress string = "BLUDGEON_HEALTHCHECK_GRPC_ADDRESS"
	EnvNameRestPort    string = "BLUDGEON_HEALTHCHECK_GRPC_PORT"
)

// defaults
const (
	DefaultPort    string = "9031"
	DefaultAddress string = "localhost"
)

type Configuration struct {
	*pkg_grpc.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: new(pkg_grpc.Configuration),
	}
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Options = append(r.Options, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		return errors.New(pkg_grpc.ErrAddressEmpty)
	}
	if r.Port == "" {
		return errors.New(pkg_grpc.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		return errors.Errorf(pkg_grpc.ErrPortBadf, r.Port)
	}
	return nil
}
