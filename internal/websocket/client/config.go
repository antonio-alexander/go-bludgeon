package client

import (
	"time"

	"github.com/pkg/errors"
)

const (
	ReadTimeout  string = "read timeout is less than 0"
	WriteTimeout string = "read timeout is less than 0"
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
