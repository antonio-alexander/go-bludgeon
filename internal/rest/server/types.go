package restserver

import (
	"net/http"
)

const (
	logAlias      string = "[rest_server] "
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
