package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /timers/{id} timers delete
// DELETE allows you to delete an existing timer, the id is required.
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
//   204: timersDeleteResponseNoContent
//   404: timersDeleteResponseNotFound

// When an timer is successfully deleted, no content is returned
// swagger:response timersDeleteResponseNoContent
type timersDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query an timer that doesn't exist
// swagger:response timersDeleteResponseNotFound
type timersDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters timer delete
type timersDeleteParams struct {
	// The timer's id
	// in:path
	ID string `json:"id"`
}
