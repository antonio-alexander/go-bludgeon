package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
)

// swagger:route GET /employees/search employees search
// You can read one or more employees by providing zero or no search parameters.
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
//   200: employeeSearchResponseOk
//   500: employeeGetResponseNotFound

// swagger:response employeeGetResponseOk
type employeeSearchResponseOk struct {
	// in:body
	Body []data.Employee
}

// swagger:parameters employee read
type employeeSearchParams struct {
	// in:body
	data.EmployeeSearch
}
