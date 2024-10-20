package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/pkg/errors"
)

// swagger:route DELETE /changes/registrations/{registration_id} registrations delete_registrations
// Creates an change, the email address and it can't be changed post create.
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
//   200: RegistrationDeleteResponseOK
//   500: RegistrationDeleteResponseError

// When a registration is successfully deleted, no content is returned
// swagger:response RegistrationDeleteResponseOK
type RegistrationDeleteResponseOK struct {
	// in:body
	Body struct{}
}

// This is the general response when a non-specific error occurs
// swagger:response RegistrationDeleteResponseError
type RegistrationDeleteResponseError struct {
	// in:body
	Body errors.Error
}

// swagger:parameters  delete_registrations
type RegistrationDeleteParams struct {
	// The registration id
	// in:path
	RegistrationId string `json:"registration_id"`
}
