package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/pkg/errors"
)

// swagger:route GET /changes/{change_id} changes get_change
// Reads an change using its id.
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
//   200: ChangeGetResponseOk
//   404: ChangeGetResponseNotFound

// This response is provided when the change is found
// swagger:response ChangeGetResponseOk
type ChangeGetResponseOk struct {
	// in:body
	Body data.Change
}

// This respons is provided when a change is not found
// swagger:response ChangeGetResponseNotFound
type ChangeGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters get_change
type ChangeGetParams struct {
	// The id of a change
	// in:path
	ChangeId string `json:"change_id"`
}
