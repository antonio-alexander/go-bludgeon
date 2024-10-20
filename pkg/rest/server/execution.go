package restserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Vars(request *http.Request) map[string]string {
	return mux.Vars(request)
}
