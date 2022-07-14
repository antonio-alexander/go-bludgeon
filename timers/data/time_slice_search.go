package data

import (
	"fmt"
	"strconv"
	"strings"
)

// swagger:model TimeSliceSearch
//TimeSliceSearch can be used to search for one or more
// time slices
type TimeSliceSearch struct {
	//Set to search for completed time slices only
	// in:query
	Completed *bool

	//Set to search for time slices associated with a specific timer
	// in:query
	TimerID *string

	//Set to search for time slices asociated with one or more timers
	// in:query
	TimerIDs []string

	//Set to search for time slices with one or more ids
	// in:query
	IDs []string
}

//ToParams can be used to generate a parameter string from
// a valid time slice search pointer
func (e *TimeSliceSearch) ToParams() string {
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
	if timerID := e.TimerID; timerID != nil && *timerID != "" {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterTimerID, *timerID))
	}
	if len(e.TimerIDs) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterTimerIDs, strings.Join(e.TimerIDs, ",")))
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

//FromParams can be used to convert a set of params into a time slice
// search pointer
func (e *TimeSliceSearch) FromParams(params map[string][]string) {
	for key, value := range params {
		switch strings.ToLower(key) {
		case ParameterTimerID:
			e.TimerID = new(string)
			*e.TimerID = value[0]
		case ParameterCompleted:
			if completed, err := strconv.ParseBool(value[0]); err == nil {
				e.Completed = new(bool)
				*e.Completed = completed
			}
		case ParameterIDs:
			for _, value := range value {
				e.IDs = append(e.IDs, strings.Split(value, ",")...)
			}
		case ParameterTimerIDs:
			for _, value := range value {
				e.TimerIDs = append(e.TimerIDs, strings.Split(value, ",")...)
			}
		}
	}
}
