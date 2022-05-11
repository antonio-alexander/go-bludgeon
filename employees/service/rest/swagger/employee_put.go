package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// swagger:route PUT /employees/{id} employees update
// POST allows you to update an existing employee, the id is required and the email address cannot be changed.
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
//   200: employeePutResponseOK
//   500: employeePutResponseError

// This is the response when an employee is successfully updated, it will include all items of employee that are user-editable as well as other items that are not user editable such as audit information and email address which can't be edited post creation.
// swagger:response employeePutResponseOK
type employeePutResponseOK struct {
	// in:body
	Body data.Employee
}

// This is the general response when a non-specific error occurs
// swagger:response employeePutResponseError
type employeePutResponseError struct {
	// in:body
	Body errors.Error
}

//These parameters must be provided for creation, email address is required
// swagger:parameters employee update
type employeePutParams struct {
	// This allows you to partially set values for certain properties of an employee, the only required parameter (specifically for update) is the email address. Any omitted fields (other than email address) will not be set and be null (rather than just empty).
	// in: body
	Body data.EmployeePartial
}
