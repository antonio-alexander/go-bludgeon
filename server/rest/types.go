package rest

import (
	"net/http"
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
	ErrStarted    string = "rest server started"
	ErrNotStarted string = "rest server not started"
)

type Owner interface {
	Start(config *Configuration) (err error)
}
