package bludgeonrestconfig

import (
	"time"
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
	EnvNameAddress string = "BLUDGEON_REST_ADDRESS"
	EnvNamePort    string = "BLUDGEON_REST_PORT"
	EnvNameTimeout string = "BLUDGEON_REST_TIMEOUT"
)

//defaults
const (
	DefaultAddress string        = "127.0.0.1"
	DefaultPort    string        = "8080"
	DefaultTimeout time.Duration = 5 * time.Second
)

//Configuration
type Configuration struct {
	Address string        `json:"Address"`
	Port    string        `json:"Port"`
	Timeout time.Duration `json:"Timeout"`
}
