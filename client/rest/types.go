package rest

import (
	"net/http"
)

const (
	GET    string = http.MethodGet
	PUT    string = http.MethodPut
	POST   string = http.MethodPost
	DELETE string = http.MethodDelete
	URIf   string = "http://%s:%s%s"
)
