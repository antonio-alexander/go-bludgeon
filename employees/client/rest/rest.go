package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
)

const urif string = "http://%s:%s%s"

type rest struct {
	internal_logger.Logger
	restClient interface {
		internal.Configurer
		internal.Parameterizer
		internal_rest.Client
	}
	config *Configuration
}

// New can be used to create a concrete instance of the rest client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	internal.Configurer
	internal.Parameterizer
	Client
} {
	return &rest{
		Logger:     internal_logger.NewNullLogger(),
		restClient: internal_rest.New(),
	}
}

func (r *rest) doRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error) {
	return r.restClient.DoRequest(ctx, uri, method, data)
}

func (r *rest) SetParameters(parameters ...interface{}) {
	r.restClient.SetParameters(parameters...)
}

func (r *rest) SetUtilities(parameters ...interface{}) {
	r.restClient.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case internal_logger.Logger:
			r.Logger = p
		}
	}
}

// Initialize can be used to ready the underlying pointer for use
func (r *rest) Configure(items ...interface{}) error {
	var configuration *Configuration
	var envs map[string]string

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			configuration = v
		}
	}
	if configuration == nil {
		configuration = new(Configuration)
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	r.config = configuration
	if err := r.restClient.Configure(&r.config.Configuration); err != nil {
		return err
	}
	return nil
}

// EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (r *rest) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteEmployees)
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.doRequest(ctx, uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

// EmployeeRead can be used to read a single employee given a
// valid id
func (r *rest) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

// EmployeeUpdate can be used to update the properties of a given employee
func (r *rest) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	bytes, err := json.Marshal(employeePartial)
	if err != nil {
		return nil, err
	}
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{}
	if err = json.Unmarshal(bytes, employee); err != nil {
		return nil, err
	}
	return employee, nil
}

// EmployeeDelete can be used to delete a single employee given a
// valid id
func (r *rest) EmployeeDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteEmployeesIDf, id))
	if _, err := r.doRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

// EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (r *rest) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteEmployeesSearch+search.ToParams())
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var employees []*data.Employee
	if err = json.Unmarshal(bytes, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}
