package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
)

const urif string = "http://%s:%s%s"

type rest struct {
	restclient.Client
	logger.Logger
	config *Configuration
}

//New can be used to create a concrete instance of the rest client
// that implements the interfaces of logic.Logic and Owner
func New(parameters ...interface{}) interface {
	Client
} {
	var config *Configuration
	//create a rest pointer, and range over the parameters
	// if configuration is provided, initialize (panic on
	// error)
	r := &rest{Client: restclient.New()}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
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

//Initialize can be used to ready the underlying pointer for use
func (r *rest) Initialize(config *Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if err := r.Client.Initialize(&config.Configuration); err != nil {
		return err
	}
	r.config = config
	return nil
}

//EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (r *rest) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteEmployees)
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.DoRequest(ctx, uri, http.MethodPost, bytes)
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
func (r *rest) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
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
func (r *rest) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.DoRequest(ctx, uri, http.MethodPut, bytes)
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
func (r *rest) EmployeeDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	if _, err := r.DoRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (r *rest) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteEmployeesSearch+search.ToParams())
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var employees []*data.Employee
	if err = json.Unmarshal(bytes, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}
