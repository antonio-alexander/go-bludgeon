package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /time_slices/{id} time_slices delete_time_slices
// Delete a time slice, the id is required.
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
//   204: TimeSlicesDeleteResponseNoContent
//   404: TimeSlicesDeleteResponseNotFound

// When an time slice is successfully deleted, no content is returned
// swagger:response TimeSlicesDeleteResponseNoContent
type TimeSlicesDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query a time slice that doesn't exist
// swagger:response TimeSlicesDeleteResponseNotFound
type TimeSlicesDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters delete_time_slices
type TimeSlicesDeleteParams struct {
	// in:path
	ID string `json:"id"`
}
