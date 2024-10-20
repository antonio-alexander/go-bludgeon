package restclient

import (
	"time"
)

const (
	DefaultTimeout time.Duration = 10 * time.Second
)

const (
	ErrAddressEmpty string = "Address is empty"
	ErrPortEmpty    string = "Port is empty"
	ErrPortBadf     string = "Port is a non-integer: %s"
	ErrTimeoutBadf  string = "Timeout is lte to 0: %v"
)

type Configuration struct {
	Address string        `json:"address"`
	Port    string        `json:"port"`
	Timeout time.Duration `json:"timeout"`
}
