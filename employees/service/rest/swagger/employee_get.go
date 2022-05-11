package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route GET /employees/{id} employees read
// You can read an existing employee, the id is required.
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
//   200: employeeGetResponseOk
//   404: employeeGetResponseNotFound

// swagger:response employeeGetResponseOk
type employeeGetResponseOk struct {
	// in:body
	Body data.Employee
}

// swagger:response employeeGetResponseNotFound
type employeeGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters employee read
type employeeGetParams struct {
	// in:path
	ID string `json:"id"`
}
