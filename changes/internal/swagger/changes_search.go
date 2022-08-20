package swagger

import "github.com/antonio-alexander/go-bludgeon/changes/data"

// swagger:route GET /changes/search changes search_changes
// Reads one or more changes using search parameters.
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
//   200: ChangeSearchResponseOk

// This is a valid response that may contain zero or more changes.
// swagger:response ChangeSearchResponseOk
type ChangeSearchResponseOk struct {
	// in:body
	Body data.ChangeDigest
}

// swagger:parameters search_changes
type ChangeSearchParams struct {
	data.ChangeSearch
}
