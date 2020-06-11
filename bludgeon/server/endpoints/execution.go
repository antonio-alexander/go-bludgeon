package bludgeonserverendpoints

import (
	"net/http"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

func getToken(request *http.Request) (token bludgeon.Token, err error) {
	//TODO: get token from request

	return
}

func handleResponse(writer http.ResponseWriter, errIn error, bytes []byte) (err error) {
	//check for errors, if so, write 500 internal server error
	if errIn != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(errIn.Error()))

		return
	}
	//if no error, write bytes
	_, err = writer.Write(bytes)

	return
}
