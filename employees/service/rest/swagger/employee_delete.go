package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route DELETE /employees/{id} employees delete
// You can delete an existing employee, the id is required.
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
//   204: employeeDeleteResponseNoContent
//   404: employeeDeleteResponseNotFound

// When an employee is successfully deleted, no content is returned
// swagger:response employeeDeleteResponseNoContent
type employeeDeleteResponseNoContent struct {
	// in:body
	Body struct{}
}

// This is the response when you attempt to query an employee that doesn't exist
// swagger:response employeeDeleteResponseNotFound
type employeeDeleteResponseNotFound struct {
	// in:body
	Body errors.Error
}

// swagger:parameters employee delete
type employeeDeleteParams struct {
	// The employee's id
	// in:path
	ID string `json:"id"`
}
