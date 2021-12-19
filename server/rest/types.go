package rest

import (
	"net/http"
	"time"
)

//endpoint constants
const (
	GET                    string        = http.MethodGet
	PUT                    string        = http.MethodPut
	POST                   string        = http.MethodPost
	DELETE                 string        = http.MethodDelete
	ErrStarted             string        = "rest server started"
	ErrNotStarted          string        = "rest server not started"
	DefaultShutdownTimeout time.Duration = 10 * time.Second
)

var (
	ConfigShutdownTimeout = DefaultShutdownTimeout
)

type handleFuncConfig struct {
	Route    string
	Method   string
	HandleFx func(writer http.ResponseWriter, request *http.Request)
}

type Owner interface {
	Start(config *Configuration) (err error)
}
