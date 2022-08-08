package client

import (
	"context"
	"errors"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	pb "github.com/antonio-alexander/go-bludgeon/timers/data/pb"

	internal_grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	timersClient     pb.TimersClient
	timeSlicesClient pb.TimeSlicesClient
	grpcClient       interface {
		internal_grpc.Client
		grpc.ClientConnInterface
	}
	internal_logger.Logger
}

//New can be used to create a concrete instance of the client client
// that implements the interfaces of logic.Logic and Owner
func New(parameters ...interface{}) interface {
	Client
} {
	var config *internal_grpc.Configuration

	c := &grpcClient{grpcClient: internal_grpc.New(parameters...)}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *internal_grpc.Configuration:
			config = p
		case internal_logger.Logger:
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
func (c *grpcClient) Initialize(config *internal_grpc.Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if err := c.grpcClient.Initialize(config); err != nil {
		return err
	}
	c.timersClient = pb.NewTimersClient(c.grpcClient)
	c.timeSlicesClient = pb.NewTimeSlicesClient(c.grpcClient)
	return nil
}

//TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (c *grpcClient) TimerCreate(ctx context.Context, timerPartial data.TimerPartial) (*data.Timer, error) {
	response, err := c.timersClient.TimerCreate(ctx, &pb.TimerCreateRequest{
		TimerPartial: pb.FromTimerPartial(&timerPartial),
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (c *grpcClient) TimerRead(ctx context.Context, id string) (*data.Timer, error) {
	response, err := c.timersClient.TimerRead(ctx, &pb.TimerReadRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimersRead can be used to read one or more timers depending
// on search values provided
func (c *grpcClient) TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error) {
	response, err := c.timersClient.TimersRead(ctx, &pb.TimersReadRequest{
		TimerSearch: pb.ToTimerSearch(&search),
	})
	return pb.ToTimers(response.GetTimers()), err
}

//TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (c *grpcClient) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	response, err := c.timersClient.TimerStart(ctx, &pb.TimerStartRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (c *grpcClient) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	response, err := c.timersClient.TimerStop(ctx, &pb.TimerStopRequest{
		Id: id,
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimerDelete can be used to delete a timer if it exists
func (c *grpcClient) TimerDelete(ctx context.Context, id string) error {
	_, err := c.timersClient.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: id,
	})
	return err
}

//TimerSubmit can be used to stop a timer and set completed to true
func (c *grpcClient) TimerSubmit(ctx context.Context, timerID string, finishTime *time.Time) (*data.Timer, error) {
	request := &pb.TimerSubmitRequest{
		Id: timerID,
	}
	if finishTime != nil {
		request.FinishOneof = &pb.TimerSubmitRequest_Finish{
			Finish: finishTime.UnixNano(),
		}
	}
	response, err := c.timersClient.TimerSubmit(ctx, request)
	return pb.ToTimer(response.GetTimer()), err
}

//TimerUpdateCommnet will only update the comment for timer with
// the provided id
func (c *grpcClient) TimerUpdateComment(ctx context.Context, id, comment string) (*data.Timer, error) {
	response, err := c.timersClient.TimerUpdateComment(ctx, &pb.TimerUpdateCommentRequest{
		Id:      id,
		Comment: comment,
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimerArchive will only update the archive for timer with
// the provided id
func (c *grpcClient) TimerArchive(ctx context.Context, id string, archive bool) (*data.Timer, error) {
	response, err := c.timersClient.TimerArchive(ctx, &pb.TimerArchiveRequest{
		Id:      id,
		Archive: archive,
	})
	return pb.ToTimer(response.GetTimer()), err
}

//TimeSliceCreate can be used to create a single time
// slice
func (c *grpcClient) TimeSliceCreate(ctx context.Context, timeslicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	response, err := c.timeSlicesClient.TimeSliceCreate(ctx, &pb.TimeSliceCreateRequest{
		TimeSlicePartial: pb.FromTimeSlicePartial(&timeslicePartial),
	})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

//TimeSliceRead can be used to read an existing time slice
func (c *grpcClient) TimeSliceRead(ctx context.Context, id string) (*data.TimeSlice, error) {
	response, err := c.timeSlicesClient.TimeSliceRead(ctx, &pb.TimeSliceReadRequest{Id: id})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

//TimeSliceUpdate can be used to update an existing time slice
func (c *grpcClient) TimeSliceUpdate(ctx context.Context, id string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	response, err := c.timeSlicesClient.TimeSliceUpdate(ctx, &pb.TimeSliceUpdateRequest{
		Id:               id,
		TimeSlicePartial: pb.FromTimeSlicePartial(&timeSlicePartial)})
	return pb.ToTimeSlice(response.GetTimeSlice()), err
}

//TimeSliceDelete can be used to delete an existing time slice
func (c *grpcClient) TimeSliceDelete(ctx context.Context, id string) error {
	_, err := c.timeSlicesClient.TimeSliceDelete(ctx, &pb.TimeSliceDeleteRequest{Id: id})
	return err
}

//TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (c *grpcClient) TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	response, err := c.timeSlicesClient.TimeSlicesRead(ctx, &pb.TimeSlicesReadRequest{
		TimeSliceSearch: pb.ToTimeSliceSearch(&search),
	})
	return pb.ToTimeSlices(response.GetTimeSlices()), err
}
