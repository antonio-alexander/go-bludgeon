package service

import (
	"context"
	"sync"

	pb "github.com/antonio-alexander/go-bludgeon/healthcheck/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/healthcheck/logic"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcService struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	pb.UnimplementedHealthChecksServer
	logic logic.Logic
}

// KIM: we don't need to expose this interface, but we need
// to implement it for grpc's sake
var _ pb.HealthChecksServer = &grpcService{}

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
	pb.RegisterHealthChecksServer(server, s)
}

func (s *grpcService) Healthcheck(ctx context.Context, empty *pb.Empty) (*pb.HealthCheckResponse, error) {
	healthcheck, err := s.logic.HealthCheck(ctx)
	return &pb.HealthCheckResponse{Healthcheck: pb.FromHealthCheck(healthcheck)}, err
}
