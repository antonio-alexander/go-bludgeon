package client

import (
	"context"
	"errors"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type client struct {
	pb.EmployeesClient
	client interface {
		grpcclient.Client
		grpc.ClientConnInterface
	}
	logger.Logger
	config *Configuration
}

//New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New(parameters ...interface{}) interface {
	Client
} {
	var config *Configuration

	c := &client{client: grpcclient.New(parameters...)}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			config = p
		case logger.Logger:
			c.Logger = p
		}
	}
	if config != nil {
		if err := c.Initialize(config); err != nil {
			panic(err)
		}
	}
	return c
}

//Initialize can be used to ready the underlying pointer for use
func (c *client) Initialize(config *Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	c.config = config
	if err := c.client.Initialize(&c.config.Configuration); err != nil {
		return err
	}
	c.EmployeesClient = pb.NewEmployeesClient(c.client)
	return nil
}

//EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (c *client) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeCreate(ctx, &pb.EmployeeCreateRequest{})
	return pb.ToEmployee(response.GetEmployee()), err
}

//EmployeeRead can be used to read a single employee given a
// valid id
func (c *client) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeRead(ctx, &pb.EmployeeReadRequest{Id: id})
	return pb.ToEmployee(response.GetEmployee()), err
}

//EmployeeUpdate can be used to update the properties of a given employee
func (c *client) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeUpdate(ctx, &pb.EmployeeUpdateRequest{})
	return pb.ToEmployee(response.GetEmployee()), err
}

//EmployeeDelete can be used to delete a single employee given a
// valid id
func (c *client) EmployeeDelete(ctx context.Context, id string) error {
	_, err := c.EmployeesClient.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{Id: id})
	return err
}

//EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (c *client) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeesRead(ctx, &pb.EmployeesReadRequest{
		EmployeeSearch: pb.ToEmployeeSearch(&search),
	})
	return pb.ToEmployees(response.GetEmployees()), err
}
