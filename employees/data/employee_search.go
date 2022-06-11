package data

import (
	"fmt"
	"strings"
)

// swagger:model EmployeeSearch
//EmployeeSearch can be used to search for an employee using one or more properties
type EmployeeSearch struct {
	//An array of one or more ids to search for
	// example: ["5afbea80-f36e-4e20-8763-64b5badbaf7d","a059dd1e-8406-4f8e-ace0-6cd2a5d54166"]
	// in:query
	IDs []string `json:"ids,omitempty"`

	//A single first name to search for
	// example: Jane
	// in:query
	FirstName *string `json:"first_name,omitempty"`

	//An array of one or more first names to search for
	// example: ["John", "Jane"]
	// in:query
	FirstNames []string `json:"first_names,omitempty"`

	//A single last name to search for
	// example: "Doe"
	// in:query
	LastName *string `json:"last_name,omitempty"`

	//An array of one or more first names to search for
	// example: ["John", "Jane"]
	// in:query
	LastNames []string `json:"last_names,omitempty"`

	//A single email address to search for
	// example: name@company.com
	// in:query
	EmailAddress *string `json:"email_address,omitempty"`

	//An array of one or more email addresses to search for
	// example: ["john.doe@company.com", "jane.doe@company.com"]
	// in:query
	EmailAddresses []string `json:"email_addresses,omitempty"`
}

func (e *EmployeeSearch) ToParams() string {
	//REVIEW: can we base64 encode the parameters?
	const parameterf string = "%s=%s"
	var parameters []string
	if len(e.IDs) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterIDs, strings.Join(e.IDs, ",")))
	}
	if firstName := e.FirstName; firstName != nil && *firstName != "" {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterFirstName, *firstName))
	}
	if len(e.FirstNames) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterFirstNames, strings.Join(e.FirstNames, ",")))
	}
	if lastName := e.LastName; lastName != nil && *lastName != "" {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterLastName, *lastName))
	}
	if len(e.LastNames) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterLastNames, strings.Join(e.LastNames, ",")))
	}
	if emailAddress := e.EmailAddress; emailAddress != nil && *emailAddress != "" {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterEmailAddress, *emailAddress))
	}
	if len(e.EmailAddresses) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterEmailAddresses, strings.Join(e.EmailAddresses, ",")))
	}
	return "?" + strings.Join(parameters, "&")
}

func (e *EmployeeSearch) FromParams(params map[string][]string) {
	for key, value := range params {
		switch strings.ToLower(key) {
		case ParameterIDs:
			for _, value := range value {
				e.IDs = strings.Split(value, ",")
				if len(e.IDs) > 0 {
					break
				}
			}
		case ParameterFirstName:
			e.FirstName = new(string)
			*e.FirstName = value[0]
		case ParameterFirstNames:
			for _, value := range value {
				e.FirstNames = strings.Split(value, ",")
				if len(e.FirstNames) > 0 {
					break
				}
			}
		case ParameterLastName:
			e.LastName = new(string)
			*e.LastName = value[0]
		case ParameterLastNames:
			for _, value := range value {
				e.LastNames = strings.Split(value, ",")
				if len(e.LastNames) > 0 {
					break
				}
			}
		case ParameterEmailAddress:
			e.EmailAddress = new(string)
			*e.EmailAddress = value[0]
		case ParameterEmailAddresses:
			for _, value := range value {
				e.EmailAddresses = strings.Split(value, ",")
				if len(e.EmailAddresses) > 0 {
					break
				}
			}
		}
	}
}
