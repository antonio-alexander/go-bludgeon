package swagger

import "github.com/antonio-alexander/go-bludgeon/changes/data"

// swagger:route GET /changes/registrations/search registrations search_registrations
// Reads one or more registrations using search parameters.
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
//   200: RegistrationSearchResponseOk

// This is a valid response that may contain zero or more registrations.
// swagger:response RegistrationSearchResponseOk
type RegistrationSearchResponseOk struct {
	// in:body
	Body []*data.Registration
}

// swagger:parameters search_registrations
type RegistrationSearchParams struct {
	data.RegistrationSearch
}
