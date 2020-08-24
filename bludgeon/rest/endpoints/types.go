package bludgeonrestendpoints

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

var (
	ConfigShutdownTimeout = DefaultShutdownTimeout
)

const (
	//DefaultShutdownTimeout provides a constant duration to be used for the context timeout
	// when shutting down the rest server
	DefaultShutdownTimeout = 10 * time.Second
)
