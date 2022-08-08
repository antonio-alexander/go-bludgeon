package service

import (
	"encoding/json"
	"net/http"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"

	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"

	"github.com/pkg/errors"
)

func handleResponse(writer http.ResponseWriter, err error, bytes []byte) error {
	if err != nil {
		switch {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, meta.ErrEmployeeNotFound):
			writer.WriteHeader(http.StatusNotFound)
		case errors.Is(err, meta.ErrEmployeeNotUpdated):
			writer.WriteHeader(http.StatusNotModified)
		case errors.Is(err, meta.ErrEmployeeConflictCreate) || errors.Is(err, meta.ErrEmployeeConflictUpdate):
			writer.WriteHeader(http.StatusConflict)
		}
		bytes, err = json.Marshal(&internal_errors.Error{Error: err.Error()})
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

func idFromPath(pathVariables map[string]string) string {
	id, ok := pathVariables[data.PathID]
	if !ok {
		return ""
	}
	return id
}
