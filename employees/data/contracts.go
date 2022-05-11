package data

//rest routes for employees
const (
	RouteEmployees       string = "/api/v1/employees"
	RouteEmployeesSearch string = RouteEmployees + "/search"
	RouteEmployeesID     string = RouteEmployees + "/{id}"
	RouteEmployeesIDf    string = RouteEmployees + "/%s"
)

const PathID string = "id"

//patameters for employees service
const (
	ParameterIDs            string = "ids"
	ParameterFirstName      string = "first_name"
	ParameterFirstNames     string = "first_names"
	ParameterLastName       string = "last_name"
	ParameterLastNames      string = "last_names"
	ParameterEmailAddress   string = "email_address"
	ParameterEmailAddresses string = "email_addresses"
)
