package rest

import (
	"net/http"
)

//endpoint constants
const (
	GET           string = http.MethodGet
	PUT           string = http.MethodPut
	POST          string = http.MethodPost
	DELETE        string = http.MethodDelete
	ErrStarted    string = "already started"
	ErrNotStarted string = "not started"
)

type Owner interface {
	Close()
}
