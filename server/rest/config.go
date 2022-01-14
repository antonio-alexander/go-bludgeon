package rest

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	ErrAddressEmpty        string        = "address is empty"
	ErrPortEmpty           string        = "port is empty"
	ErrPortBadf            string        = "port is a non-integer: %s"
	ErrTimeoutBadf         string        = "timeout is lte to 0: %v"
	EnvNameAddress         string        = "BLUDGEON_REST_ADDRESS"
	EnvNamePort            string        = "BLUDGEON_REST_PORT"
	EnvNameTimeout         string        = "BLUDGEON_REST_TIMEOUT"
	EnvNameShutdownTimeout string        = "BLUDGEON_REST_SHUTDOWN_TIMEOUT"
	DefaultAddress         string        = "127.0.0.1"
	DefaultPort            string        = "8080"
	DefaultTimeout         time.Duration = 5 * time.Second
	DefaultShutdownTimeout time.Duration = 10 * time.Second
)

type Configuration struct {
	Address         string        `json:"address"`
	Port            string        `json:"port"`
	Timeout         time.Duration `json:"timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Timeout = DefaultTimeout
	r.ShutdownTimeout = DefaultShutdownTimeout
}

func (r *Configuration) FromEnv(pwd string, envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameAddress]; ok {
		r.Address = address
	}
	if port, ok := envs[EnvNamePort]; ok {
		r.Port = port
	}
	if timeoutString, ok := envs[EnvNameTimeout]; ok {
		if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
			if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
				r.Timeout = timeout
			}
		}
	}
	if shutdownTimeoutString, ok := envs[EnvNameShutdownTimeout]; ok {
		if shutdownTimeoutInt, err := strconv.Atoi(shutdownTimeoutString); err == nil {
			if timeout := time.Duration(shutdownTimeoutInt) * time.Second; timeout > 0 {
				r.ShutdownTimeout = timeout
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
		return errors.New(ErrAddressEmpty)
	}
	if r.Port == "" {
		return errors.New(ErrPortEmpty)
	}
	if _, e := strconv.Atoi(r.Port); e != nil {
		return errors.Errorf(ErrPortBadf, r.Port)
	}
	if r.Timeout <= 0 {
		return errors.Errorf(ErrTimeoutBadf, r.Timeout)
	}

	return
}
