package service

import (
	"context"
	"sync"

	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcService struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	pb.UnimplementedEmployeesServer
	logic logic.Logic
}

// KIM: we don't need to expose this interface, but we need
// to implement it for grpc's sake
var _ pb.EmployeesServer = &grpcService{}

func New() interface {
	internal.Parameterizer
	server.Registerer
} {
	return &grpcService{
		Logger: logger.NewNullLogger(),
	}
}

func (s *grpcService) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *grpcService) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			s.logic = p
		}
	}
	switch {
	case s.logic == nil:
		panic("logic not set")
	}
}

func (s *grpcService) Register(server grpc.ServiceRegistrar) {
	pb.RegisterEmployeesServer(server, s)
}

func (s *grpcService) EmployeeCreate(ctx context.Context, request *pb.EmployeeCreateRequest) (*pb.EmployeeCreateResponse, error) {
	employee, err := s.logic.EmployeeCreate(ctx, *pb.ToEmployeePartial(request.GetEmployeePartial()))
	return &pb.EmployeeCreateResponse{
		Employee: pb.FromEmployee(employee),
	}, err
}

func (s *grpcService) EmployeeRead(ctx context.Context, request *pb.EmployeeReadRequest) (*pb.EmployeeReadResponse, error) {
	employee, err := s.logic.EmployeeRead(ctx, request.GetId())
	return &pb.EmployeeReadResponse{
		Employee: pb.FromEmployee(employee),
	}, err
}

func (s *grpcService) EmployeesRead(ctx context.Context, request *pb.EmployeesReadRequest) (*pb.EmployeesReadResponse, error) {
	employees, err := s.logic.EmployeesRead(ctx, *pb.FromEmployeeSearch(request.GetEmployeeSearch()))
	return &pb.EmployeesReadResponse{
		Employees: pb.FromEmployees(employees),
	}, err
}

func (s *grpcService) EmployeeUpdate(ctx context.Context, request *pb.EmployeeUpdateRequest) (*pb.EmployeeUpdateResponse, error) {
	employee, err := s.logic.EmployeeUpdate(ctx, request.GetId(), *pb.ToEmployeePartial(request.GetEmployeePartial()))
	return &pb.EmployeeUpdateResponse{
		Employee: pb.FromEmployee(employee),
	}, err
}

func (s *grpcService) EmployeeDelete(ctx context.Context, request *pb.EmployeeDeleteRequest) (*pb.EmployeeDeleteResponse, error) {
	err := s.logic.EmployeeDelete(ctx, request.GetId())
	return &pb.EmployeeDeleteResponse{}, err
}
