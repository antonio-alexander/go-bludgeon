package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route GET /changes/registrations/{registration_id}/changes registrations get_registrations_changes
// Reads all the changes associated with a registration.
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
//   200: RegistrationChangesGetResponseOk
//   404: RegistrationChangesGetResponseNotFound

// This is the response for a successful registration changes read
// swagger:response RegistrationChangesGetResponseOk
type RegistrationChangesGetResponseOk struct {
	// in:body
	Body data.ChangeDigest
}

// This is the response when the registration isn't found
// swagger:response RegistrationChangesGetResponseNotFound
type RegistrationChangesGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters get_registrations_changes
type RegistrationChangesGetParams struct {
	// The registration id
	// in:path
	RegistrationId string `json:"registration_id"`
}
