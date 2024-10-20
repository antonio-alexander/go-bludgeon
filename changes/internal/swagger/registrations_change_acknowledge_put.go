package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/pkg/errors"
)

// swagger:route POST /changes/registrations/{registration_id}/acknowledge registrations registrations_acknowledge
// Acknowledges one or more changes for a given registration
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
//   200: RegistrationChangesAcknowledgePostResponseOK
//   500: RegistrationChangesAcknowledgePostResponseError

// This is the response when one or more changes are successfully acknowledged
// swagger:response RegistrationChangesAcknowledgePostResponseOK
type RegistrationChangesAcknowledgePostResponseOK struct {
	// in:body
	Body data.ResponseAcknowledge
}

// This is the general response when a non-specific error occurs
// swagger:response RegistrationChangesAcknowledgePostResponseError
type RegistrationChangesAcknowledgePostResponseError struct {
	// in:body
	Body errors.Error
}

// swagger:parameters  registrations_acknowledge
type RegistrationChangesAcknowledgePostParams struct {
	// The registration id
	// in:path
	RegistrationId string `json:"registration_id"`

	// This allows you to partially set values for certain properties of an Change.
	// in: body
	Body data.RequestAcknowledge
}
