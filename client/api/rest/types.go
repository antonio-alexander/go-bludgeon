package restapi

import (
	"net/http"
	"time"
)

const (
	URIf string = "http://%s:%s%s"
)

//endpoint constants
const (
	GET    = http.MethodGet
	PUT    = http.MethodPut
	POST   = http.MethodPost
	DELETE = http.MethodDelete
)

//configuration defaults
const (
	DefaultTimeout time.Duration = 5 * time.Second
)

//configuration variables
var (
	ConfigTimeout = DefaultTimeout
)
