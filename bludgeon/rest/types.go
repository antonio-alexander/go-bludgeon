package bludgeonrest

import (
	"net/http"
)

type HandleFuncConfig struct {
	Route    string
	Method   string
	HandleFx func(writer http.ResponseWriter, request *http.Request)
}
