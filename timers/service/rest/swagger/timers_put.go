package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

// swagger:route PUT /timers/{id} timers update
// POST allows you to update an timer, the only unique constraint is email address and that can't be changed post update.
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
//   200: timersPutResponseOK
//   500: timersPutResponseError

// This is the response when an timer is successfully updated, it will include all items of timer that are user-editable as well as other items that are not user editable such as audit information and email address which can't be edited post creation.
// swagger:response timersPutResponseOK
type timersPutResponseOK struct {
	// in:body
	Body data.Timer
}

// This is the general response when a non-specific error occurs
// swagger:response timersPutResponseError
type timersPutResponseError struct {
	// in:body
	Body errors.Error
}

//These parameters must be provided for creation, email address is required
// swagger:parameters timers update
type timerPutParams struct {
	// This allows you to partially set values for certain properties of an timer, the only required parameter (specifically for update) is the email address. Any omitted fields (other than email address) will not be set and be null (rather than just empty).
	// in: body
	Body data.TimerPartial
}
