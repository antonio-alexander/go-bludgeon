package rest

import (
	"net/http"
	"time"
)

//endpoint constants
const (
	GET    = http.MethodGet
	PUT    = http.MethodPut
	POST   = http.MethodPost
	DELETE = http.MethodDelete
)

//errors
const (
	ErrAddressEmpty string = "Address is empty"
	ErrPortEmpty    string = "Port is empty"
	ErrPortBadf     string = "Port is a non-integer: %s"
	ErrTimeoutBadf  string = "Timeout is lte to 0: %v"
	ErrStarted      string = "server started"
	ErrNotStarted   string = "server not started"
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
