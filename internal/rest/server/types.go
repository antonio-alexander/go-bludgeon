package restserver

import (
	"net/http"
)

const (
	ErrStarted    string = "already started"
	ErrNotStarted string = "not started"
)

type Router interface {
	HandleFunc(HandleFuncConfig)
}

type HandleFuncConfig struct {
	Route    string
	Method   string
	HandleFx func(http.ResponseWriter, *http.Request)
}

type Owner interface {
	Start(config *Configuration) (err error)
	Stop()
}
