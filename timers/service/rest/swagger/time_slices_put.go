package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route PUT /time_slices/{id} time_slices update_time_slices
// Update a time slice.
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
//   200: TimeSlicesPutResponseOK
//   500: TimeSlicesPutResponseError

// This is the response when an timer is successfully updated, it will include all items of timer that are user-editable as well as other items that are not user editable such as audit information and email address which can't be edited post creation.
// swagger:response TimeSlicesPutResponseOK
type TimeSlicesPutResponseOK struct {
	// in:body
	Body data.TimeSlice
}

// This is the general response when a non-specific error occurs
// swagger:response TimeSlicesPutResponseError
type TimeSlicesPutResponseError struct {
	// in:body
	Body errors.Error
}

//These parameters must be provided for creation, email address is required
// swagger:parameters update_time_slices
type TimeSlicesPutParams struct {
	// in:path
	ID string `json:"id"`

	// This allows you to partially set values for certain properties of an timer, the only required parameter (specifically for update) is the email address. Any omitted fields (other than email address) will not be set and be null (rather than just empty).
	// in: body
	Body data.TimeSlicePartial
}
