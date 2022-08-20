package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /changes/{change_id} changes delete_changes
// Deletes a change using id.
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
//   204: ChangesDeleteResponseNoContent
//   404: ChangesDeleteResponseNotFound

// When an change is successfully deleted, no content is returned
// swagger:response ChangesDeleteResponseNoContent
type ChangesDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query a change that doesn't exist
// swagger:response ChangesDeleteResponseNotFound
type ChangesDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters changes delete_changes
type ChangesDeleteParams struct {
	// The change's id
	// in:path
	ChangeId string `json:"change_id"`
}
