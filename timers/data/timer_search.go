package data

import (
	"fmt"
	"strconv"
	"strings"
)

//TimerSearch can be used to inclusively search for one or more
// timers
type TimerSearch struct {
	//Set to search for timers associated with a specific employee
	// in:query
	EmployeeID *string `json:"employee_id,omitempty"`

	//Set to search for timers associated with one or more employees
	// in:query
	EmployeeIDs []string `json:"employee_ids,omitempty"`

	//Set to search for completed timers only
	// in:query
	Completed *bool `json:"completed,omitempty"`

	//Set to search for archived timers only
	// in:query
	Archived *bool `json:"archived,omitempty"`

	//An array of one or more ids to search for
	// in:query
	IDs []string `json:"ids,omitempty"`
}

//ToParams can be used to generate a parameter string from
// a valid timer search pointer
func (e *TimerSearch) ToParams() string {
	//REVIEW: can we base64 encode the parameters?
	const (
		parameterf     string = "%s=%s"
		parameterBoolf string = "%s=%t"
	)
	var parameters []string
	if len(e.IDs) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterIDs, strings.Join(e.IDs, ",")))
	}
	if employeeID := e.EmployeeID; employeeID != nil && *employeeID != "" {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterEmployeeID, *employeeID))
	}
	if len(e.EmployeeIDs) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterEmployeeIDs, strings.Join(e.EmployeeIDs, ",")))
	}
	if completed := e.Completed; completed != nil {
		parameters = append(parameters,
			fmt.Sprintf(parameterBoolf, ParameterCompleted, *completed))
	}
	if archived := e.Completed; archived != nil {
		parameters = append(parameters,
			fmt.Sprintf(parameterBoolf, ParameterArchived, *archived))
	}
	return "?" + strings.Join(parameters, "&")
}

//FromParams can be used to convert a set of params into a timer
// search pointer
func (e *TimerSearch) FromParams(params map[string][]string) {
	for key, value := range params {
		switch strings.ToLower(key) {
		case ParameterEmployeeID:
			e.EmployeeID = new(string)
			*e.EmployeeID = value[0]
		case ParameterCompleted:
			if completed, err := strconv.ParseBool(value[0]); err == nil {
				e.Completed = new(bool)
				*e.Completed = completed
			}
		case ParameterArchived:
			if archived, err := strconv.ParseBool(value[0]); err == nil {
				e.Archived = new(bool)
				*e.Archived = archived
			}
		case ParameterIDs:
			for _, value := range value {
				e.IDs = append(e.IDs, strings.Split(value, ",")...)
			}
		case ParameterEmployeeIDs:
			for _, value := range value {
				e.EmployeeIDs = append(e.EmployeeIDs, strings.Split(value, ",")...)
			}
		}
	}
}
