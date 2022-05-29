package data

import (
	"fmt"
	"strconv"
	"strings"
)

// swagger:model TimerSearch
//TimerSearch can be used to inclusively search for one or more
// employees
type TimerSearch struct {
	//Set to search for timers associated with a specific employee
	// example: "00831cfb-2d37-4027-b150-8f3b2db18c49"
	EmployeeID *string `json:"employee_id,omitempty"`

	//Set to search for timers associated with one or more employees
	// example: ["748d751f-0c89-4b6b-a684-c45d1acac8d8","042dc313-58e7-47fd-9ea5-42c633d61506"]
	EmployeeIDs []string `json:"employee_ids,omitempty"`

	//Set to search for completed timers only
	// example: false
	Completed *bool `json:"completed,omitempty"`

	//Set to search for archived timers only
	// example: false
	Archived *bool `json:"archived,omitempty"`

	//An array of one or more ids to search for
	// example: ["71a91f01-4250-461d-9cbc-6cb6c4dd7d2f","d6d98618-4724-44ba-8193-af6920dcbb3e"]
	IDs []string `json:"ids,omitempty"`
}

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
