package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /timers/search timers search_timers
// Read one or more timers using search parameters.
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
type TimersSearchResponseOk struct {
	// in:body
	Body []data.Timer
}

// swagger:response TimersGetResponseNotFound
type TimersSearchResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters timers search_timers
type TimersSearchParams struct {
	data.TimerSearch
}
