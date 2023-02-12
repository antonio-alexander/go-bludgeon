package grpc

import (
	"context"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	pb "github.com/antonio-alexander/go-bludgeon/timers/data/pb"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	logger.Logger
	timersClient     pb.TimersClient
	timeSlicesClient pb.TimeSlicesClient
	client           interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		grpc.ClientConnInterface
	}
}

// New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
	client.Client
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
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			g.Logger = p
		}
	}
}

func (g *grpcClient) Configure(items ...interface{}) error {
	var envs map[string]string
	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		c = new(Configuration)
		c.Default()
		c.FromEnv(envs)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	if err := g.client.Configure(&c.Configuration); err != nil {
		return err
	}
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (g *grpcClient) Initialize() error {
	if err := g.client.Initialize(); err != nil {
		return err
	}
	g.timersClient = pb.NewTimersClient(g.client)
	g.timeSlicesClient = pb.NewTimeSlicesClient(g.client)
	return nil
}

func (g *grpcClient) Shutdown() {
	g.client.Shutdown()
}

// TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (g *grpcClient) TimerCreate(ctx context.Context, timerPartial data.TimerPartial) (*data.Timer, error) {
	response, err := g.timersClient.TimerCreate(ctx, &pb.TimerCreateRequest{
		TimerPartial: pb.FromTimerPartial(&timerPartial),
	})
	return pb.ToTimer(response.GetTimer()), err
}

// TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (g *grpcClient) TimerRead(ctx context.Context, id string) (*data.Timer, error) {
	response, err := g.timersClient.TimerRead(ctx, &pb.TimerReadRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

// TimersRead can be used to read one or more timers depending
// on search values provided
func (g *grpcClient) TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error) {
	response, err := g.timersClient.TimersRead(ctx, &pb.TimersReadRequest{
		TimerSearch: pb.ToTimerSearch(&search),
	})
	return pb.ToTimers(response.GetTimers()), err
}

// TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (g *grpcClient) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	response, err := g.timersClient.TimerStart(ctx, &pb.TimerStartRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

// TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (g *grpcClient) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	response, err := g.timersClient.TimerStop(ctx, &pb.TimerStopRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

// TimerDelete can be used to delete a timer if it exists
func (g *grpcClient) TimerDelete(ctx context.Context, id string) error {
	_, err := g.timersClient.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: id,
	})
	return err
}

// TimerSubmit can be used to stop a timer and set completed to true
func (g *grpcClient) TimerSubmit(ctx context.Context, timerID string, finishTime int64) (*data.Timer, error) {
	request := &pb.TimerSubmitRequest{
		Id: timerID,
	}
	request.FinishOneof = &pb.TimerSubmitRequest_Finish{
		Finish: finishTime,
	}
	response, err := g.timersClient.TimerSubmit(ctx, request)
	return pb.ToTimer(response.GetTimer()), err
}

func (g *grpcClient) TimerUpdate(ctx context.Context, id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	response, err := g.timersClient.TimerUpdate(ctx, &pb.TimerUpdateRequest{
		Id:           id,
		TimerPartial: pb.FromTimerPartial(&timerPartial),
	})
	return pb.ToTimer(response.GetTimer()), err
}

// TimeSliceCreate can be used to create a single time
// slice
func (g *grpcClient) TimeSliceCreate(ctx context.Context, timeslicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	response, err := g.timeSlicesClient.TimeSliceCreate(ctx, &pb.TimeSliceCreateRequest{
		TimeSlicePartial: pb.FromTimeSlicePartial(&timeslicePartial),
	})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

// TimeSliceRead can be used to read an existing time slice
func (g *grpcClient) TimeSliceRead(ctx context.Context, id string) (*data.TimeSlice, error) {
	response, err := g.timeSlicesClient.TimeSliceRead(ctx, &pb.TimeSliceReadRequest{Id: id})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

// TimeSliceUpdate can be used to update an existing time slice
func (g *grpcClient) TimeSliceUpdate(ctx context.Context, id string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	response, err := g.timeSlicesClient.TimeSliceUpdate(ctx, &pb.TimeSliceUpdateRequest{
		Id:               id,
		TimeSlicePartial: pb.FromTimeSlicePartial(&timeSlicePartial)})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

// TimeSliceDelete can be used to delete an existing time slice
func (g *grpcClient) TimeSliceDelete(ctx context.Context, id string) error {
	_, err := g.timeSlicesClient.TimeSliceDelete(ctx, &pb.TimeSliceDeleteRequest{Id: id})
	return err
}

// TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (g *grpcClient) TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	response, err := g.timeSlicesClient.TimeSlicesRead(ctx, &pb.TimeSlicesReadRequest{
		TimeSliceSearch: pb.ToTimeSliceSearch(&search),
	})
	return pb.ToTimeSlices(response.GetTimeSlices()), err
}
