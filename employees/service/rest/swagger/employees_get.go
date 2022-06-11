package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route GET /employees/{id} employees read
// Reads an employee using their id.
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
//   200: EmployeeGetResponseOk
//   404: EmployeeGetResponseNotFound

// swagger:response EmployeeGetResponseOk
type EmployeeGetResponseOk struct {
	// in:body
	Body data.Employee
}

// swagger:response EmployeeGetResponseNotFound
type EmployeeGetResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters employees read
type EmployeeGetParams struct {
	// in:path
	ID string `json:"id"`
}
