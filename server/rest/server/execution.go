package restserver

import (
	"log"
	"net/http"

	bludgeon "github.com/antonio-alexander/go-bludgeon/internal/common"
)

func println(log *log.Logger, v ...interface{}) {
	if log != nil {
		log.Println(v...)
	}
}

func printf(log *log.Logger, format string, v ...interface{}) {
	if log != nil {
		log.Printf(format, v...)
	}
}

func print(log *log.Logger, v ...interface{}) {
	if log != nil {
		log.Print(v...)
	}
}

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
