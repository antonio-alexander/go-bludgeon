package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/internal/rest/client"
)

const urif string = "http://%s:%s%s"

type rest struct {
	client.Client
	logger.Logger
	config *client.Configuration
}

type Owner interface {
	Initialize(config *client.Configuration) error
}

func New(parameters ...interface{}) interface {
	logic.Logic
	Owner
} {
	var config *client.Configuration
	r := &rest{
		Client: client.New(),
	}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *client.Configuration:
			config = p
		case logger.Logger:
			r.Logger = p
		}
	}
	if config != nil {
		if err := r.Initialize(config); err != nil {
			panic(err)
		}
	}
	return r
}

func (r *rest) Initialize(config *client.Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if err := r.Client.Initialize(config); err != nil {
		return err
	}
	r.config = config
	return nil
}

//EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (r *rest) EmployeeCreate(employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteEmployees)
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.DoRequest(uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

//EmployeeRead can be used to read a single employee given a
// valid id
func (r *rest) EmployeeRead(id string) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

//EmployeeUpdate can be used to update the properties of a given employee
func (r *rest) EmployeeUpdate(id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.DoRequest(uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

//EmployeeDelete can be used to delete a single employee given a
// valid id
func (r *rest) EmployeeDelete(id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	if _, err := r.DoRequest(uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (r *rest) EmployeesRead(search data.EmployeeSearch) ([]*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteEmployeesSearch+search.ToParams())
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var employees []*data.Employee
	if err = json.Unmarshal(bytes, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}
