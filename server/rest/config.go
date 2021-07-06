package rest

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

//errors
const (
	ErrAddressEmpty string = "Address is empty"
	ErrPortEmpty    string = "Port is empty"
	ErrPortBadf     string = "Port is a non-integer: %s"
	ErrTimeoutBadf  string = "Timeout is lte to 0: %v"
)

//environmental variables
const (
	EnvNameRestAddress string = "BLUDGEON_REST_ADDRESS"
	EnvNameRestPort    string = "BLUDGEON_REST_PORT"
	EnvNameRestTimeout string = "BLUDGEON_REST_TIMEOUT"
)

//defaults
const (
	DefaultAddress string        = "127.0.0.1"
	DefaultPort    string        = "8080"
	DefaultTimeout time.Duration = 5 * time.Second
)

type Configuration struct {
	Address string        `json:"Address"`
	Port    string        `json:"Port"`
	Timeout time.Duration `json:"Timeout"`
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Timeout = DefaultTimeout
}

func (r *Configuration) FromEnv(pwd string, envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameRestAddress]; ok {
		r.Address = address
	}
	if port, ok := envs[EnvNameRestPort]; ok {
		r.Port = port
	}
	if timeoutString, ok := envs[EnvNameRestTimeout]; ok {
		if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
			if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
				r.Timeout = timeout
			}
		}
	}
}

func (r *Configuration) Validate() (err error) {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is lte 0
	if r.Address == "" {
		err = errors.New(ErrAddressEmpty)

		return
	}
	if r.Port == "" {
		err = errors.New(ErrPortEmpty)

		return
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		err = errors.Errorf(ErrPortBadf, r.Port)

		return
	}
	if r.Timeout <= 0 {
		err = errors.Errorf(ErrTimeoutBadf, r.Timeout)
	}

	return
}
