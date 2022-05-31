package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /time_slices/{id} time_slices read_time_slices
// read a time slice using its id.
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
//   200: TimeSlicesGetResponseOk
//   404: TimeSlicesGetResponseNotFound

// swagger:response TimeSlicesGetResponseOk
type TimeSlicesGetResponseOk struct {
	// in:body
	Body data.TimeSlice
}

// swagger:response TimeSlicesGetResponseNotFound
type TimeSlicesGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters read_time_slices
type TimeSlicesGetParams struct {
	// in:path
	ID string `json:"id"`
}
