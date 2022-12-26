package server

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	ReadTimeout  string = "read timeout is less than 0"
	WriteTimeout string = "write timeout is less than 0"
	PingTimeout  string = "ping timeout is less than 0"
	PingInterval string = "ping interval is less than 0"
)

const (
	EnvNameReadTimeout  string = "BLUDGEON_WEBSOCKET_READ_TIMEOUT"
	EnvNameWriteTimeout string = "BLUDGEON_WEBSOCKET_WRITE_TIMEOUT"
	EnvNamePingTimeout  string = "BLUDGEON_WEBSOCKET_PING_TIMEOUT"
	EnvNamePingInterval string = "BLUDGEON_WEBSOCKET_PING_INTERVAL"
)

const (
	DefaultReadTimeout  time.Duration = 10 * time.Second
	DefaultWriteTimeout time.Duration = 10 * time.Second
	DefaultPingTimeout  time.Duration = 20 * time.Second
	DefaultPingInterval time.Duration = 5 * time.Second
)

var (
	ErrReadTimeout  = errors.New(ReadTimeout)
	ErrWriteTimeout = errors.New(WriteTimeout)
	ErrPingTimeout  = errors.New(PingTimeout)
	ErrPingInterval = errors.New(PingInterval)
)

type Configuration struct {
	WriteTimeout time.Duration `json:"write_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	PingInterval time.Duration `json:"ping_interval"`
	PingTimeout  time.Duration `json:"ping_timeout"`
}

func (c *Configuration) Default() {
	c.ReadTimeout = DefaultReadTimeout
	c.WriteTimeout = DefaultWriteTimeout
	c.PingInterval = DefaultPingInterval
	c.PingTimeout = DefaultPingTimeout
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if s, ok := envs[EnvNameReadTimeout]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			if readTimeout := time.Duration(i) * time.Second; readTimeout > 0 {
				c.ReadTimeout = readTimeout
			}
		}
	}
	if s, ok := envs[EnvNameWriteTimeout]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			if writeTimeout := time.Duration(i) * time.Second; writeTimeout > 0 {
				c.WriteTimeout = writeTimeout
			}
		}
	}
	if s, ok := envs[EnvNamePingTimeout]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			if pingTimeout := time.Duration(i) * time.Second; pingTimeout > 0 {
				c.PingTimeout = pingTimeout
			}
		}
	}
	if s, ok := envs[PingInterval]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			if pingInterval := time.Duration(i) * time.Second; pingInterval > 0 {
				c.PingInterval = pingInterval
			}
		}
	}
}

func (c *Configuration) Validate() error {
	if c.ReadTimeout < 0 {
		return ErrReadTimeout
	}
	if c.WriteTimeout < 0 {
		return ErrWriteTimeout
	}
	if c.PingTimeout < 0 {
		return ErrPingTimeout
	}
	if c.PingInterval < 0 {
		return ErrPingInterval
	}
	return nil
}
