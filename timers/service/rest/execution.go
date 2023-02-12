package rest

import (
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

func idFromPath(pathVariables map[string]string) string {
	id, ok := pathVariables[data.PathID]
	if !ok {
		return ""
	}
	return id
}
