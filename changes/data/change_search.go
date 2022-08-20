package data

import (
	"fmt"
	"strconv"
	"strings"
)

type ChangeSearch struct {
	ChangeIds     []string `json:"change_ids,omitempty"`
	DataIds       []string `json:"data_ids,omitempty"`
	Types         []string `json:"types,omitempty"`
	Actions       []string `json:"actions,omitempty"`
	ServiceNames  []string `json:"service_names,omitempty"`
	LatestVersion *bool    `json:"latest_version,omitempty"`
	Since         *int64   `json:"since,string,omitempty"`
}

func (c *ChangeSearch) ToParams() string {
	const parameterf string = "%s=%s"
	var parameters []string

	if len(c.ChangeIds) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterChangeIds, strings.Join(c.ChangeIds, ",")))
	}
	if len(c.DataIds) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterDataIds, strings.Join(c.DataIds, ",")))
	}
	if len(c.Types) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterTypes, strings.Join(c.Types, ",")))
	}
	if len(c.ServiceNames) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterServiceNames, strings.Join(c.ServiceNames, ",")))
	}
	if len(c.Actions) > 0 {
		parameters = append(parameters,
			fmt.Sprintf(parameterf, ParameterActions, strings.Join(c.Actions, ",")))
	}
	if c.LatestVersion != nil {
		parameters = append(parameters, fmt.Sprint(*c.LatestVersion))
	}
	if c.Since != nil {
		parameters = append(parameters, fmt.Sprint(*c.Since))
	}
	return "?" + strings.Join(parameters, "&")
}

func (c *ChangeSearch) FromParams(params map[string][]string) {
	for key, value := range params {
		switch strings.ToLower(key) {
		case ParameterActions:
			for _, value := range value {
				if value != "" {
					c.Actions = append(c.Actions, strings.Split(value, ",")...)
				}
			}
		case ParameterChangeIds:
			for _, value := range value {
				if value != "" {
					c.ChangeIds = append(c.ChangeIds, strings.Split(value, ",")...)
				}
			}
		case ParameterDataIds:
			for _, value := range value {
				if value != "" {
					c.DataIds = append(c.DataIds, strings.Split(value, ",")...)
				}
			}
		case ParameterTypes:
			for _, value := range value {
				if value != "" {
					c.Types = append(c.Types, strings.Split(value, ",")...)
				}
			}
		case ParameterServiceNames:
			for _, value := range value {
				if value != "" {
					c.ServiceNames = append(c.ServiceNames, strings.Split(value, ",")...)
				}
			}
		case ParameterLatestVersion:
			for _, value := range value {
				latestVersion, err := strconv.ParseBool(value)
				if err != nil {
					continue
				}
				c.LatestVersion = new(bool)
				*c.LatestVersion = latestVersion
				break
			}
		case ParameterSince:
			for _, value := range value {
				since, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					continue
				}
				c.Since = new(int64)
				*c.Since = since
				break
			}
		}
	}
}
