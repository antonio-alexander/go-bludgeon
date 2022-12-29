package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /timers/{id} timers delete_timers
// Delete a timer, the id is required.
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
//   204: TimersDeleteResponseNoContent
//   404: TimersDeleteResponseNotFound

// When an timer is successfully deleted, no content is returned
// swagger:response TimersDeleteResponseNoContent
type TimersDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query an timer that doesn't exist
// swagger:response TimersDeleteResponseNotFound
type TimersDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters delete_timers
type TimersDeleteParams struct {
	// in:path
	ID string `json:"id"`
}
