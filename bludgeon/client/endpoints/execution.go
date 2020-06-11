package bludgeonclientendpoints

import "net/http"

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
