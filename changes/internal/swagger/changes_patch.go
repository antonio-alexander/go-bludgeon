package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route PATCH /changes changes patch_changes
// Upserts a change, DataId and DataType are required.
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
//   200: ChangePatchResponseOK
//   500: ChangePatchResponseError

// This is the response when an Change is successfully upserted, it will include all items of Change that are user-editable as well as other items that are not user editable such as audit information which can't be edited post creation.
// swagger:response ChangePatchResponseOK
type ChangePatchResponseOK struct {
	// in:body
	Body data.ResponseChange
}

// This is the general response when a non-specific error occurs
// swagger:response ChangePatchResponseError
type ChangePatchResponseError struct {
	// in:body
	Body errors.Error
}

// swagger:parameters patch_changes
type ChangePatchParams struct {
	// This allows you to partially set values for certain properties of an Change.
	// in: body
	Body data.RequestChange
}
