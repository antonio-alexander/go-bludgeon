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
	ErrStarted    string = "server started"
	ErrNotStarted string = "server not started"
)
