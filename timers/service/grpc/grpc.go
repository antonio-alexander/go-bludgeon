package service

import (
	"context"
	"sync"
	"time"

	pb "github.com/antonio-alexander/go-bludgeon/timers/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcService struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	pb.UnimplementedTimersServer
	pb.UnimplementedTimeSlicesServer
	logic logic.Logic
}

// KIM: we don't need to expose this interface, but we need
// to implement it for grpc's sake
var (
	_ pb.TimersServer     = &grpcService{}
	_ pb.TimeSlicesServer = &grpcService{}
)

func New(parameters ...interface{}) interface {
	internal.Parameterizer
	server.Registerer
} {
	return &grpcService{
		Logger: logger.NewNullLogger(),
	}
}

func (s *grpcService) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
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
}

func (s *grpcService) Register(server grpc.ServiceRegistrar) {
	pb.RegisterTimersServer(server, s)
	pb.RegisterTimeSlicesServer(server, s)
}

func (s *grpcService) TimerCreate(ctx context.Context, request *pb.TimerCreateRequest) (*pb.TimerCreateResponse, error) {
	timer, err := s.logic.TimerCreate(ctx, *pb.ToTimerPartial(request.TimerPartial))
	return &pb.TimerCreateResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimerRead(ctx context.Context, request *pb.TimerReadRequest) (*pb.TimerReadResponse, error) {
	timer, err := s.logic.TimerRead(ctx, request.GetId())
	return &pb.TimerReadResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimerUpdate(ctx context.Context, request *pb.TimerUpdateRequest) (*pb.TimerUpdateResponse, error) {
	timer, err := s.logic.TimerUpdate(ctx, request.GetId(), *pb.ToTimerPartial(request.GetTimerPartial()))
	return &pb.TimerUpdateResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimerDelete(ctx context.Context, request *pb.TimerDeleteRequest) (*pb.TimerDeleteResponse, error) {
	err := s.logic.TimerDelete(ctx, request.GetId())
	return &pb.TimerDeleteResponse{}, err
}

func (s *grpcService) TimersRead(ctx context.Context, request *pb.TimersReadRequest) (*pb.TimersReadResponse, error) {
	timers, err := s.logic.TimersRead(ctx, *pb.FromTimerSearch(request.GetTimerSearch()))
	return &pb.TimersReadResponse{Timers: pb.FromTimers(timers)}, err
}

func (s *grpcService) TimerStart(ctx context.Context, request *pb.TimerStartRequest) (*pb.TimerStartResponse, error) {
	timer, err := s.logic.TimerStart(ctx, request.GetId())
	return &pb.TimerStartResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimerStop(ctx context.Context, request *pb.TimerStopRequest) (*pb.TimerStopResponse, error) {
	timer, err := s.logic.TimerStop(ctx, request.GetId())
	return &pb.TimerStopResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimerSubmit(ctx context.Context, request *pb.TimerSubmitRequest) (*pb.TimerSubmitResponse, error) {
	finish := time.Now()
	if request.FinishOneof != nil {
		finish = time.Unix(0, request.GetFinish())
	}
	timer, err := s.logic.TimerSubmit(ctx, request.GetId(), finish.UnixNano())
	return &pb.TimerSubmitResponse{Timer: pb.FromTimer(timer)}, err
}

func (s *grpcService) TimeSliceCreate(ctx context.Context, request *pb.TimeSliceCreateRequest) (*pb.TimeSliceCreateResponse, error) {
	timeSlice, err := s.logic.TimeSliceCreate(ctx, *pb.ToTimeSlicePartial(request.GetTimeSlicePartial()))
	return &pb.TimeSliceCreateResponse{TimeSlice: pb.FromTimeSlice(timeSlice)}, err
}

func (s *grpcService) TimeSliceRead(ctx context.Context, request *pb.TimeSliceReadRequest) (*pb.TimeSliceReadResponse, error) {
	timeSlice, err := s.logic.TimeSliceRead(ctx, request.GetId())
	return &pb.TimeSliceReadResponse{TimeSlice: pb.FromTimeSlice(timeSlice)}, err
}

func (s *grpcService) TimeSliceUpdate(ctx context.Context, request *pb.TimeSliceUpdateRequest) (*pb.TimeSliceUpdateResponse, error) {
	timeSlice, err := s.logic.TimeSliceUpdate(ctx, request.GetId(), *pb.ToTimeSlicePartial(request.GetTimeSlicePartial()))
	return &pb.TimeSliceUpdateResponse{TimeSlice: pb.FromTimeSlice(timeSlice)}, err
}

func (s *grpcService) TimeSliceDelete(ctx context.Context, request *pb.TimeSliceDeleteRequest) (*pb.TimeSliceDeleteResponse, error) {
	err := s.logic.TimeSliceDelete(ctx, request.GetId())
	return &pb.TimeSliceDeleteResponse{}, err
}

func (s *grpcService) TimeSlicesRead(ctx context.Context, request *pb.TimeSlicesReadRequest) (*pb.TimeSlicesReadResponse, error) {
	timeSlices, err := s.logic.TimeSlicesRead(ctx, *pb.FromTimeSliceSearch(request.GetTimeSliceSearch()))
	return &pb.TimeSlicesReadResponse{TimeSlices: pb.FromTimeSlices(timeSlices)}, err
}
