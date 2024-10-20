package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/pkg/errors"
)

// swagger:route PATCH /changes/registrations registrations patch_registrations
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
//   200: RegistrationUpsertPatchResponseOK
//   500: RegistrationUpsertPatchResponseError

// This is the response when an registration is successfully upserted.
// swagger:response RegistrationUpsertPatchResponseOK
type RegistrationUpsertPatchResponseOK struct {
	// in:body
	Body data.ResponseRegister
}

// This is the general response when a non-specific error occurs
// swagger:response RegistrationUpsertPatchResponseError
type RegistrationUpsertPatchResponseError struct {
	// in:body
	Body errors.Error
}

// swagger:parameters  patch_registrations
type RegistrationUpsertPatchParams struct {
	// This allows you to upsert a registration.
	// in: body
	Body data.RequestRegister
}
