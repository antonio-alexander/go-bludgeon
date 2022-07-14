package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route GET /time_slices/search time_slices search_time_slices
// Read one or more time slices.
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
type TimeSlicesSearchResponseOk struct {
	// in:body
	Body []data.TimeSlice
}

// swagger:response TimeSlicesGetResponseNotFound
type TimeSlicesSearchResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters time_slices search_time_slices
type TimeSlicesSearchParams struct {
	data.TimeSliceSearch
}
