package rest

import (
	"encoding/json"
	"net/http"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"

	"github.com/pkg/errors"
)

func handleResponse(writer http.ResponseWriter, err error, bytes []byte) error {
	if err != nil {
		switch {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, meta.ErrTimerNotFound):
			writer.WriteHeader(http.StatusNotFound)
		case errors.Is(err, meta.ErrTimerNotUpdated):
			writer.WriteHeader(http.StatusNotModified)
		case errors.Is(err, meta.ErrTimerConflictCreate) || errors.Is(err, meta.ErrTimerConflictUpdate):
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
