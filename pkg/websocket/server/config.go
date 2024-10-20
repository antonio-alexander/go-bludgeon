package server

import (
	"time"

	"github.com/pkg/errors"
)

const (
	DefaultReadTimeout  time.Duration = 10 * time.Second
	DefaultWriteTimeout time.Duration = 10 * time.Second
	DefaultPingTimeout  time.Duration = 10 * time.Second
	DefaultPingInterval time.Duration = 2 * time.Second
)

const (
	ReadTimeout  string = "read timeout is less than 0"
	WriteTimeout string = "write timeout is less than 0"
	PingTimeout  string = "ping timeout is less than 0"
	PingInterval string = "ping interval is less than 0"
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
