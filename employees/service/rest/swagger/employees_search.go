package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
)

// swagger:route GET /employees/search employees search
// Reads one or more employees using search parameters.
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
//   200: EmployeeSearchResponseOk

// swagger:response EmployeeSearchResponseOk
type EmployeeSearchResponseOk struct {
	// in:body
	Body []data.Employee
}

// swagger:parameters EmployeeSearchParams
type EmployeeSearchParams struct {
	data.EmployeeSearch
}
