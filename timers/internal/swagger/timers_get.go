package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /timers/{id} timers read_timers
// Read a timer using its id.
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
//   200: TimersGetResponseOk
//   404: TimersGetResponseNotFound

// swagger:response TimersGetResponseOk
type TimersGetResponseOk struct {
	// in:body
	Body data.Timer
}

// swagger:response TimersGetResponseNotFound
type TimersGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters read_timers
type TimersGetParams struct {
	// in:path
	ID string `json:"id"`
}
