package grpc

import (
	"context"

	client "github.com/antonio-alexander/go-bludgeon/employees/client"
	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	logger.Logger
	pb.EmployeesClient
	client interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		grpc.ClientConnInterface
	}
}

// New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	client.Client
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	return &grpcClient{
		Logger: logger.NewNullLogger(),
		client: grpcclient.New(),
	}
}

func (g *grpcClient) SetParameters(parameters ...interface{}) {
	g.client.SetParameters(parameters...)
}

func (g *grpcClient) SetUtilities(parameters ...interface{}) {
	g.client.SetUtilities(parameters...)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			g.Logger = p
		}
	}
}

func (g *grpcClient) Configure(items ...interface{}) error {
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
	if g == nil {
		configuration = new(Configuration)
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	if err := g.client.Configure(&configuration.Configuration); err != nil {
		return err
	}
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (g *grpcClient) Initialize() error {
	if err := g.client.Initialize(); err != nil {
		return err
	}
	g.EmployeesClient = pb.NewEmployeesClient(g.client)
	return nil
}

func (g *grpcClient) Shutdown() {
	g.client.Shutdown()
}

// EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (g *grpcClient) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := g.EmployeesClient.EmployeeCreate(ctx, &pb.EmployeeCreateRequest{
		EmployeePartial: pb.FromEmployeePartial(&employeePartial),
	})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeRead can be used to read a single employee given a
// valid id
func (g *grpcClient) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	response, err := g.EmployeesClient.EmployeeRead(ctx, &pb.EmployeeReadRequest{Id: id})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeUpdate can be used to update the properties of a given employee
func (g *grpcClient) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	response, err := g.EmployeesClient.EmployeeUpdate(ctx, &pb.EmployeeUpdateRequest{
		Id:              id,
		EmployeePartial: pb.FromEmployeePartial(&employeePartial),
	})
	return pb.ToEmployee(response.GetEmployee()), err
}

// EmployeeDelete can be used to delete a single employee given a
// valid id
func (g *grpcClient) EmployeeDelete(ctx context.Context, id string) error {
	_, err := g.EmployeesClient.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{Id: id})
	return err
}

// EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (g *grpcClient) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	response, err := g.EmployeesClient.EmployeesRead(ctx, &pb.EmployeesReadRequest{
		EmployeeSearch: pb.ToEmployeeSearch(&search),
	})
	return pb.ToEmployees(response.GetEmployees()), err
}
