package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /employees/{id} employees delete
// Deletes an employee using id.
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
//   204: EmployeeDeleteResponseNoContent
//   404: EmployeeDeleteResponseNotFound

// When an Employee is successfully deleted, no content is returned
// swagger:response EmployeeDeleteResponseNoContent
type EmployeeDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query an Employee that doesn't exist
// swagger:response EmployeeDeleteResponseNotFound
type EmployeeDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters employees delete
type EmployeeDeleteParams struct {
	// The employee's id
	// in:path
	ID string `json:"id"`
}
