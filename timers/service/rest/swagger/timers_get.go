package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /timers/{id} timers read
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
type timersGetResponseOk struct {
	// in:body
	Body data.Timer
}

// swagger:response timersGetResponseNotFound
type timersGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters timers read
type timersGetParams struct {
	// in:path
	ID string `json:"id"`
}
