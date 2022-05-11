package data

// swagger:model Employee
//Employee represents a person uniquely identified by their email address
type Employee struct {
	// The ID of an employee (v4 UUID)
	// example: 86fa2f09-d260-11ec-bd5d-0242c0a8e002
	ID string `json:"id"`

	// The first name of an employee
	// example: John
	FirstName string `json:"first_name"`

	// The last name of an employee
	// example: Smith
	LastName string `json:"last_name"`

	// The email address of an employee
	// example: John.Smith@foobar.duck
	EmailAddress string `json:"email_address"`

	//Used for accounting of this unique employee
	Audit
}

// swagger:model EmployeePartial
//EmployeePartial provides a way to optionally/partially update different fields of an employee
type EmployeePartial struct {
	//The first name of an employee, this is optional
	// example: Jane
	FirstName *string `json:"first_name,omitempty"`

	//The last name of an employee, this is optional
	// example: Doe
	LastName *string `json:"last_name,omitempty"`

	//The email address of an employee, this is optional, but can't conflict with existing employees
	// example: Jane.Doe@foobar.duck
	EmailAddress *string `json:"email_address,omitempty"`
}
