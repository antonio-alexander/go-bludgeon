package rest

import (
	"strconv"
	"time"

	internal_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/client"

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
	*internal_rest.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: new(internal_rest.Configuration),
	}
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Timeout = internal_rest.DefaultTimeout
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
	if timeoutString, ok := envs[internal_rest.EnvNameRestTimeout]; ok {
		if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
			if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
				r.Timeout = timeout
			}
		}
	}
}

func (r *Configuration) Validate() error {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is gt 0
	if r.Address == "" {
		return errors.New(internal_rest.ErrAddressEmpty)
	}
	if r.Port == "" {
		return errors.New(internal_rest.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		return errors.Errorf(internal_rest.ErrPortBadf, r.Port)
	}
	if r.Timeout <= 0 {
		return errors.Errorf(internal_rest.ErrTimeoutBadf, r.Timeout)
	}
	return nil
}
