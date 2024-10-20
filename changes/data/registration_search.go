package data

import (
	"fmt"
	"strings"
)

type RegistrationSearch struct {
	RegistrationIds []string `json:"registration_ids,omitempty"`
}

func (r *RegistrationSearch) ToParams() string {
	const parameterf string = "%s=%s"
	var parameters []string

	if len(r.RegistrationIds) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterChangeIds, strings.Join(r.RegistrationIds, ",")))
	}
	return "?" + strings.Join(parameters, "&")
}

func (r *RegistrationSearch) FromParams(params map[string][]string) {
	for key, value := range params {
		switch strings.ToLower(key) {
		case ParameterRegistrationIds:
			for _, value := range value {
				if value != "" {
					r.RegistrationIds = append(r.RegistrationIds, strings.Split(value, ",")...)
				}
			}
		}
	}
}
