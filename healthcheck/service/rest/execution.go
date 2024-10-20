package service

import (
	"encoding/json"
	"net/http"

	pkg_errors "github.com/antonio-alexander/go-bludgeon/pkg/errors"
)

func handleResponse(writer http.ResponseWriter, err error, bytes []byte) error {
	if err != nil {
		var e pkg_errors.Error

		switch {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		}
		switch v := err.(type) {
		default:
			e = pkg_errors.New(err.Error())
		case pkg_errors.Error:
			e = v
		}
		bytes, err = json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = writer.Write(bytes)
		return err
	}
	if bytes == nil {
		writer.WriteHeader(http.StatusNoContent)
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = writer.Write(bytes)
	return err
}
