package grpc

import (
	"context"
	"sync"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	pb "github.com/antonio-alexander/go-bludgeon/timers/data/pb"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	grpcclient "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type grpcClient struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	timersClient     pb.TimersClient
	timeSlicesClient pb.TimeSlicesClient
	client           interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		grpc.ClientConnInterface
	}
	cache          cache.Cache
	changesClient  changesclient.Client
	changesHandler changesclient.Handler
	config         *Configuration
	configured     bool
	stopper        chan struct{}
	initialized    bool
	handlerId      string
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

func (g *grpcClient) cacheWrite(key string, item interface{}) error {
	if g.configured && g.config.DisableCache {
		return errors.New("cache disabled")
	}
	return g.cache.Write(key, item)
}

func (g *grpcClient) cacheRead(key string, item interface{}) error {
	if g.configured && g.config.DisableCache {
		return errors.New("cache disabled")
	}
	return g.cache.Read(key, item)
}

func (g *grpcClient) cacheDelete(key string) error {
	if g.configured && g.config.DisableCache {
		return errors.New("cache disabled")
	}
	return g.cache.Delete(key)
}

func (g *grpcClient) registrationChangeAcknowledge(serviceName string, changeIds ...string) {
	if len(changeIds) <= 0 {
		return
	}
	g.Add(1)
	go func() {
		defer g.Done()

		ctx, cancel := context.WithTimeout(context.Background(), g.config.ChangesTimeout)
		defer cancel()
		if err := g.changesClient.RegistrationChangeAcknowledge(ctx, g.config.ChangesRegistrationId, changeIds...); err != nil {
			g.Error("error while acknowledging changes: %s", err)
		}
	}()
}

func (g *grpcClient) handleChanges(changes ...*changesdata.Change) error {
	var changesToAcknowledge []string

	for _, change := range changes {
		switch {
		case change.DataType == data.ChangeTypeTimer && ((change.DataAction == data.ChangeActionUpdate) ||
			(change.DataAction == data.ChangeActionStart) || (change.DataAction == data.ChangeActionStop) ||
			(change.DataAction == data.ChangeActionSubmit)):
			failure := false
			timer := &data.Timer{}
			if err := g.cacheRead(change.DataId, timer); err != nil {
				//REVIEW: should this use context.Background()?
				timerRead, err := g.TimerRead(context.Background(), change.DataId)
				if err != nil {
					failure = true
					break
				}
				g.cacheWrite(timerRead.ID, timerRead)
			}
			if failure {
				break
			}
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		case change.DataType == data.ChangeTypeTimer && change.DataAction == data.ChangeActionDelete:
			g.cacheDelete(change.DataId)
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		}
	}
	g.registrationChangeAcknowledge(g.config.ChangesRegistrationId, changesToAcknowledge...)
	return nil
}

func (g *grpcClient) launchChangeHandler() {
	started := make(chan struct{})
	g.Add(1)
	go func() {
		defer g.Done()

		checkChangesFx := func() {
			ctx, cancel := context.WithTimeout(context.Background(), g.config.ChangesTimeout)
			defer cancel()
			changesRead, err := g.changesClient.RegistrationChangesRead(ctx, g.config.ChangesRegistrationId)
			if err != nil {
				g.Error("error while reading registration changes: %s", err)
				return
			}
			if len(changesRead) == 0 {
				return
			}
			if err := g.handleChanges(changesRead...); err != nil {
				g.Error("error while reading registration changes: %s", err)
			}
		}
		tCheck := time.NewTicker(g.config.ChangeRateRead)
		defer tCheck.Stop()
		close(started)
		for {
			select {
			case <-g.stopper:
				return
			case <-tCheck.C:
				checkChangesFx()
			}
		}
	}()
	<-started
}

func (g *grpcClient) launchChangeRegistration() {
	started := make(chan struct{})
	g.Add(1)
	go func() {
		defer g.Done()

		var registered, handlerSet bool
		var err error

		tRegister := time.NewTicker(g.config.ChangeRateRegistration)
		defer tRegister.Stop()
		close(started)
		for {
			select {
			case <-g.stopper:
				return
			case <-tRegister.C:
				if !handlerSet {
					if g.handlerId, err = g.changesHandler.HandlerCreate(g.handleChanges); err != nil {
						g.Error("error while creating change handler: %s", err)
						break
					}
					g.Debug("Change handler created: %s (%s)", g.handlerId, g.config.ChangesRegistrationId)
					handlerSet = true
				}
				if !registered {
					ctx, cancel := context.WithTimeout(context.Background(), g.config.ChangesTimeout)
					defer cancel()
					if err := g.changesClient.RegistrationUpsert(ctx, g.config.ChangesRegistrationId); err != nil {
						g.Error("error while upserting change registration: %s", err)
						break
					}
					g.Debug("Change registration upserted for: %s", g.config.ChangesRegistrationId)
					registered = true
				}
				if handlerSet && registered {
					return
				}
			}
		}
	}()
	<-started
}

func (g *grpcClient) SetParameters(parameters ...interface{}) {
	g.Lock()
	defer g.Unlock()

	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			internal.Configurer
			internal.Initializer
			internal.Parameterizer
			grpc.ClientConnInterface
		}:
			g.client = p
		case interface {
			cache.Cache
		}:
			g.cache = p
		case interface {
			changesclient.Handler
			changesclient.Client
		}:
			g.changesHandler = p
			g.changesClient = p
		case changesclient.Handler:
			g.changesHandler = p
		case changesclient.Client:
			g.changesClient = p
		}
	}
	switch {
	case g.changesHandler == nil:
		panic("changes handler not set")
	case g.changesClient == nil:
		panic("changes client not set")
	case g.cache == nil:
		panic("cache is nil")
	case g.client == nil:
		panic("client is nil")
	}
	g.client.SetParameters(parameters...)
}

func (g *grpcClient) SetUtilities(parameters ...interface{}) {
	g.Lock()
	defer g.Unlock()

	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			g.Logger = p
		}
	}
	g.client.SetUtilities(parameters...)
}

func (g *grpcClient) Configure(items ...interface{}) error {
	g.Lock()
	defer g.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			c = new(Configuration)
			c.FromEnvs(v)
		case *Configuration:
			c = v
		}
	}
	if err := c.Validate(); err != nil {
		return err
	}
	g.config = c
	g.configured = true
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (g *grpcClient) Initialize() error {
	g.Lock()
	defer g.Unlock()

	if err := g.client.Initialize(); err != nil {
		return err
	}
	if g.configured && !g.config.DisableCache {
		g.stopper = make(chan struct{})
		g.launchChangeRegistration()
		g.launchChangeHandler()
	}
	g.timersClient = pb.NewTimersClient(g.client)
	g.timeSlicesClient = pb.NewTimeSlicesClient(g.client)
	g.initialized = true
	return nil
}

func (g *grpcClient) Shutdown() {
	g.Lock()
	defer g.Unlock()

	if !g.initialized {
		return
	}
	if g.configured && !g.config.DisableCache {
		close(g.stopper)
		g.Wait()
	}
	g.client.Shutdown()
	g.initialized = false
	g.configured = false
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
