package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route POST /employees employees create
// Creates an employee, the email address and it can't be changed post create.
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
//   200: EmployeePostResponseOK
//   500: EmployeePostResponseError

// This is the response when an Employee is successfully created, it will include all items of Employee that are user-editable as well as other items that are not user editable such as audit information and email address which can't be edited post creation.
// swagger:response EmployeePostResponseOK
type EmployeePostResponseOK struct {
	// in:body
	Body data.Employee
}

// This is the general response when a non-specific error occurs
// swagger:response EmployeePostResponseError
type EmployeePostResponseError struct {
	// in:body
	Body errors.Error
}

//These parameters must be provided for creation, email address is required
// swagger:parameters employees create
type EmployeePostParams struct {
	// This allows you to partially set values for certain properties of an Employee, the only required parameter (specifically for create) is the email address. Any omitted fields (other than email address) will not be set and be null (rather than just empty).
	// in: body
	Body data.EmployeePartial
}
