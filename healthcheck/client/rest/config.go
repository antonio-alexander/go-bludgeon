package rest

import (
	"strconv"

	pkg_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/client"

	"github.com/pkg/errors"
)

// environmental variables
const (
	EnvNameRestAddress string = "BLUDGEON_HEALTHCHECK_REST_ADDRESS"
	EnvNameRestPort    string = "BLUDGEON_HEALTHCHECK_REST_PORT"
)

// defaults
const (
	DefaultPort    string = "9030"
	DefaultAddress string = "localhost"
)

type Configuration struct {
	*pkg_rest.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: new(pkg_rest.Configuration),
	}
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Timeout = pkg_rest.DefaultTimeout
}

func (r *Configuration) FromEnv(envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameRestAddress]; ok {
		r.Address = address
	}
	if port, ok := envs[EnvNameRestPort]; ok {
		r.Port = port
	}
	// if timeoutString, ok := envs[pkg_rest.EnvNameRestTimeout]; ok {
	// 	if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
	// 		if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
	// 			r.Timeout = timeout
	// 		}
	// 	}
	// }
}

func (r *Configuration) Validate() error {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is gt 0
	if r.Address == "" {
		return errors.New(pkg_rest.ErrAddressEmpty)
	}
	if r.Port == "" {
		return errors.New(pkg_rest.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		return errors.Errorf(pkg_rest.ErrPortBadf, r.Port)
	}
	if r.Timeout <= 0 {
		return errors.Errorf(pkg_rest.ErrTimeoutBadf, r.Timeout)
	}
	return nil
}
