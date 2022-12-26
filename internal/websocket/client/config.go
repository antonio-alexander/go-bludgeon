package client

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	ReadTimeout  string = "read timeout is less than 0"
	WriteTimeout string = "read timeout is less than 0"
)

const (
	EnvNameReadTimeout  string = "BLUDGEON_WEBSOCKET_READ_TIMEOUT"
	EnvNameWriteTimeout string = "BLUDGEON_WEBSOCKET_WRITE_TIMEOUT"
)

const (
	DefaultReadTimeout  time.Duration = 10 * time.Second
	DefaultWriteTimeout time.Duration = 10 * time.Second
)

var (
	ErrReadTimeout  = errors.New(ReadTimeout)
	ErrWriteTimeout = errors.New(WriteTimeout)
)

type Configuration struct {
	WriteTimeout time.Duration `json:"write_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
}

func (c *Configuration) Default() {
	c.ReadTimeout = DefaultReadTimeout
	c.WriteTimeout = DefaultWriteTimeout
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if readTimeoutString, ok := envs[EnvNameReadTimeout]; ok {
		if readtimeoutInt, err := strconv.Atoi(readTimeoutString); err == nil {
			if readTimeout := time.Duration(readtimeoutInt) * time.Second; readTimeout > 0 {
				c.ReadTimeout = readTimeout
			}
		}
	}
	if writeTimeoutString, ok := envs[EnvNameWriteTimeout]; ok {
		if writeTimeoutInt, err := strconv.Atoi(writeTimeoutString); err == nil {
			if writeTimeout := time.Duration(writeTimeoutInt) * time.Second; writeTimeout > 0 {
				c.WriteTimeout = writeTimeout
			}
		}
	}
}

func (c *Configuration) Validate() (err error) {
	if c.ReadTimeout < 0 {
		return ErrReadTimeout
	}
	if c.WriteTimeout < 0 {
		return ErrWriteTimeout
	}
	return
}
