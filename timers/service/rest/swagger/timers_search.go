package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /timers/search timers search
// GET allows you to read an timer, the id is required.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
// responses:
//   200: timersGetResponseOk
//   404: timersGetResponseNotFound

// swagger:response timersGetResponseOk
type timerSearchResponseOk struct {
	// in:body
	Body data.Timer
}

// swagger:response timersGetResponseNotFound
type timerSearchResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters timers search
type timerSearchParams struct {
	// in:path
	ID string `json:"id"`
}
