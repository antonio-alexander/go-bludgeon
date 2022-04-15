package data

//Employee
type Employee struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	Audit
}

type EmployeePartial struct {
	FirstName    *string
	LastName     *string
	EmailAddress *string
}

type EmployeeSearch struct {
	IDs            []string
	FirstName      *string
	FirstNames     []string
	LastName       *string
	LastNames      []string
	EmailAddress   *string
	EmailAddresses []string
}
