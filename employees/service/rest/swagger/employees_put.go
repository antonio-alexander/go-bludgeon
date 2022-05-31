package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route PUT /employees/{id} employees update
// Updates an existing employee using their id, the email address cannot be changed.
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
//   200: EmployeePutResponseOK
//   500: EmployeePutResponseError

// This is the response when an Employee is successfully updated, it will include all items of Employee that are user-editable as well as other items that are not user editable such as audit information and email address which can't be edited post creation.
// swagger:response EmployeePutResponseOK
type EmployeePutResponseOK struct {
	// in:body
	Body data.Employee
}

// This is the general response when a non-specific error occurs
// swagger:response EmployeePutResponseError
type EmployeePutResponseError struct {
	// in:body
	Body errors.Error
}

//These parameters must be provided for creation, email address is required
// swagger:parameters employees update
type EmployeePutParams struct {
	// The Employee's id
	// in:path
	ID string `json:"id"`

	// This allows you to partially set values for certain properties of an Employee, the only required parameter (specifically for update) is the email address. Any omitted fields (other than email address) will not be set and be null (rather than just empty).
	// in: body
	Body data.EmployeePartial
}
