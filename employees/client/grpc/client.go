package grpc

import (
	"context"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type client struct {
	logger.Logger
	pb.EmployeesClient
	grpcClient interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		grpc.ClientConnInterface
	}
}

// New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	Client
	internal.Initializer
	internal.Configurer
} {
	return &client{
		Logger:     logger.NewNullLogger(),
		grpcClient: grpcclient.New(),
	}
}

func (c *client) SetParameters(parameters ...interface{}) {
	c.grpcClient.SetParameters(parameters...)
}

func (c *client) SetUtilities(parameters ...interface{}) {
	c.grpcClient.SetUtilities(parameters...)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			c.Logger = p
		}
	}
}

func (c *client) Configure(items ...interface{}) error {
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
	if c == nil {
		configuration = new(Configuration)
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	if err := c.grpcClient.Configure(&configuration.Configuration); err != nil {
		return err
	}
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (c *client) Initialize() error {
	if err := c.grpcClient.Initialize(); err != nil {
		return err
	}
	c.EmployeesClient = pb.NewEmployeesClient(c.grpcClient)
	return nil
}

func (c *client) Shutdown() {
	c.grpcClient.Shutdown()
}

// EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (c *client) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeCreate(ctx, &pb.EmployeeCreateRequest{})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeRead can be used to read a single employee given a
// valid id
func (c *client) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeRead(ctx, &pb.EmployeeReadRequest{Id: id})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeUpdate can be used to update the properties of a given employee
func (c *client) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeeUpdate(ctx, &pb.EmployeeUpdateRequest{})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeDelete can be used to delete a single employee given a
// valid id
func (c *client) EmployeeDelete(ctx context.Context, id string) error {
	_, err := c.EmployeesClient.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{Id: id})
	return err
}

// EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (c *client) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	response, err := c.EmployeesClient.EmployeesRead(ctx, &pb.EmployeesReadRequest{
		EmployeeSearch: pb.ToEmployeeSearch(&search),
	})
	return pb.ToEmployees(response.GetEmployees()), err
}
