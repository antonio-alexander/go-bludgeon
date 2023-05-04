package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	"github.com/pkg/errors"
)

type restClient struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	client interface {
		internal.Configurer
		internal.Parameterizer
		restclient.Client
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

// New will create a populated instance of the rest client
// if Configuration is provided as a parameter, it will
// also attempt to initialize (and panic on error)
func New() interface {
	client.Client
	internal.Parameterizer
	internal.Configurer
	internal.Initializer
} {
	return &restClient{
		Logger: logger.NewNullLogger(),
		client: restclient.New(),
	}
}

func (r *restClient) cacheWrite(key string, item interface{}) error {
	if r.configured && r.config.DisableCache {
		return errors.New("cache disabled")
	}
	return r.cache.Write(key, item)
}

func (r *restClient) cacheRead(key string, item interface{}) error {
	if r.configured && r.config.DisableCache {
		return errors.New("cache disabled")
	}
	return r.cache.Read(key, item)
}

func (r *restClient) cacheDelete(key string) error {
	if r.configured && r.config.DisableCache {
		return errors.New("cache disabled")
	}
	return r.cache.Delete(key)
}

func (r *restClient) registrationChangeAcknowledge(serviceName string, changeIds ...string) {
	if len(changeIds) <= 0 {
		return
	}
	r.Add(1)
	go func() {
		defer r.Done()

		ctx, cancel := context.WithTimeout(context.Background(), r.config.ChangesTimeout)
		defer cancel()
		if err := r.changesClient.RegistrationChangeAcknowledge(ctx, r.config.ChangesRegistrationId, changeIds...); err != nil {
			r.Error("error while acknowledging changes: %s", err)
		}
	}()
}

func (r *restClient) handleChanges(changes ...*changesdata.Change) error {
	var changesToAcknowledge []string

	for _, change := range changes {
		switch {
		default:
			fmt.Println("??" + change.DataType)
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		case change.DataType == data.ChangeTypeTimer && ((change.DataAction == data.ChangeActionUpdate) ||
			(change.DataAction == data.ChangeActionStart) || (change.DataAction == data.ChangeActionStop) ||
			(change.DataAction == data.ChangeActionSubmit)):
			failure := false
			timer := &data.Timer{}
			if err := r.cacheRead(change.DataId, timer); err == nil {
				//REVIEW: should this use context.Background()?
				timerRead, err := r.TimerRead(context.Background(), change.DataId)
				if err != nil {
					failure = true
					break
				}
				r.cacheWrite(timerRead.ID, timerRead)
			}
			if failure {
				break
			}
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		case change.DataType == data.ChangeTypeTimer && change.DataAction == data.ChangeActionDelete:
			fmt.Println("!!" + change.DataId)
			r.cacheDelete(change.DataId)
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		}
	}
	r.registrationChangeAcknowledge(r.config.ChangesRegistrationId, changesToAcknowledge...)
	return nil
}

func (r *restClient) launchChangeHandler() {
	started := make(chan struct{})
	r.Add(1)
	go func() {
		defer r.Done()

		checkChangesFx := func() {
			ctx, cancel := context.WithTimeout(context.Background(), r.config.ChangesTimeout)
			defer cancel()
			changesRead, err := r.changesClient.RegistrationChangesRead(ctx, r.config.ChangesRegistrationId)
			if err != nil {
				r.Error("error while reading registration changes: %s", err)
				return
			}
			if len(changesRead) == 0 {
				return
			}
			if err := r.handleChanges(changesRead...); err != nil {
				r.Error("error while reading registration changes: %s", err)
			}
		}
		tCheck := time.NewTicker(r.config.ChangeRateRead)
		defer tCheck.Stop()
		close(started)
		for {
			select {
			case <-r.stopper:
				return
			case <-tCheck.C:
				checkChangesFx()
			}
		}
	}()
	<-started
}

func (r *restClient) launchChangeRegistration() {
	started := make(chan struct{})
	r.Add(1)
	go func() {
		defer r.Done()

		var registered, handlerSet bool
		var err error

		tRegister := time.NewTicker(r.config.ChangeRateRegistration)
		defer tRegister.Stop()
		close(started)
		for {
			select {
			case <-r.stopper:
				return
			case <-tRegister.C:
				if !handlerSet {
					if r.handlerId, err = r.changesHandler.HandlerCreate(r.handleChanges); err != nil {
						r.Error("error while creating change handler: %s", err)
						break
					}
					r.Debug("Change handler created: %s (%s)", r.handlerId, r.config.ChangesRegistrationId)
					handlerSet = true
				}
				if !registered {
					ctx, cancel := context.WithTimeout(context.Background(), r.config.ChangesTimeout)
					defer cancel()
					if err := r.changesClient.RegistrationUpsert(ctx, r.config.ChangesRegistrationId); err != nil {
						r.Error("error while upserting change registration: %s", err)
						break
					}
					r.Debug("Change registration upserted for: %s", r.config.ChangesRegistrationId)
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

func (r *restClient) doRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error) {
	bytes, statusCode, err := r.client.DoRequest(ctx, uri, method, data)
	if err != nil {
		return nil, err
	}
	switch statusCode {
	default:
		return nil, errors.Errorf("error status code: %d", statusCode)
	case http.StatusOK, http.StatusNoContent:
		return bytes, nil
	}
}

func (r *restClient) SetParameters(parameters ...interface{}) {
	r.Lock()
	defer r.Unlock()

	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			cache.Cache
		}:
			r.cache = p
		case interface {
			restclient.Client
			internal.Configurer
			internal.Parameterizer
		}:
			r.client = p
		case interface {
			changesclient.Handler
			changesclient.Client
		}:
			r.changesHandler = p
			r.changesClient = p
		case changesclient.Handler:
			r.changesHandler = p
		case changesclient.Client:
			r.changesClient = p
		}
	}
	switch {
	case r.changesHandler == nil:
		panic("changes handler not set")
	case r.changesClient == nil:
		panic("changes client not set")
	case r.cache == nil:
		panic("cache is nil")
	case r.client == nil:
		panic("client is nil")
	}
	r.client.SetParameters(parameters...)
}

func (r *restClient) SetUtilities(parameters ...interface{}) {
	r.Lock()
	defer r.Unlock()

	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			r.Logger = p
		}
	}
	r.client.SetUtilities(parameters...)
}

func (r *restClient) Configure(items ...interface{}) error {
	r.Lock()
	defer r.Unlock()

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
	r.config = c
	r.configured = true
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (r *restClient) Initialize() error {
	r.Lock()
	defer r.Unlock()

	if r.initialized {
		return nil
	}
	if !r.configured {
		return errors.New("not configured")
	}
	if !r.config.DisableCache {
		r.stopper = make(chan struct{})
		r.launchChangeRegistration()
		r.launchChangeHandler()
	}
	r.initialized = true
	return nil
}

func (r *restClient) Shutdown() {
	r.Lock()
	defer r.Unlock()

	if !r.initialized {
		return
	}
	if r.configured && !r.config.DisableCache {
		close(r.stopper)
		r.Wait()
	}
	r.initialized = false
	r.configured = false
}

// TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (r *restClient) TimerCreate(ctx context.Context, timerPartial data.TimerPartial) (*data.Timer, error) {
	bytes, err := json.Marshal(&timerPartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimers, r.config.Address, r.config.Port)
	bytes, err = r.doRequest(ctx, uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

// TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (r *restClient) TimerRead(ctx context.Context, timerId string) (*data.Timer, error) {
	timer := new(data.Timer)
	if err := r.cacheRead(timerId, timer); err == nil {
		return timer, nil
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDf,
		r.config.Address, r.config.Port, timerId)
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimersRead can be used to read one or more timers depending
// on search values provided
func (r *restClient) TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersSearch+"%s",
		r.config.Address, r.config.Port, search.ToParams())
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var timers = []*data.Timer{}
	if err = json.Unmarshal(bytes, &timers); err != nil {
		return nil, err
	}
	return timers, nil
}

// TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (r *restClient) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDStartf,
		r.config.Address, r.config.Port, id)
	bytes, err := r.doRequest(ctx, uri, http.MethodPut, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

// TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (r *restClient) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDStopf,
		r.config.Address, r.config.Port, id)
	bytes, err := r.doRequest(ctx, uri, http.MethodPut, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

func (r *restClient) TimerUpdate(ctx context.Context, id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	bytes, err := json.Marshal(timerPartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDf,
		r.config.Address, r.config.Port, id)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

// TimerArchive will only update the archive for timer with
// the provided id
func (r *restClient) TimerArchive(ctx context.Context, id string, archive bool) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Archived: &archive,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDArchivef,
		r.config.Address, r.config.Port, id)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

// TimerDelete can be used to delete a timer if it exists
func (r *restClient) TimerDelete(ctx context.Context, timerId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDf,
		r.config.Address, r.config.Port, timerId)
	if _, err := r.doRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	if err := r.cacheDelete(timerId); err != nil {
		r.Error("error while deleting timer (%s) from cache: %s", timerId, err)
	}
	return nil
}

// TimerSubmit can be used to stop a timer and set completed to true
func (r *restClient) TimerSubmit(ctx context.Context, timerId string, finishTime int64) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.Contract{
		Finish: finishTime,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDSubmitf,
		r.config.Address, r.config.Port, timerId)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timer.ID, timer); err != nil {
		r.Error("error while writing timer (%s) to cache: %s", timer.ID, err)
	}
	return timer, nil
}

// TimeSliceCreate can be used to create a single time
// slice
func (r *restClient) TimeSliceCreate(ctx context.Context, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlices,
		r.config.Address, r.config.Port)
	bytes, err = r.doRequest(ctx, uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timeSlice.ID, timeSlice); err != nil {
		r.Error("error while writing time slice (%s) to cache: %s", timeSlice.ID, err)
	}
	return timeSlice, nil
}

// TimeSliceRead can be used to read an existing time slice
func (r *restClient) TimeSliceRead(ctx context.Context, timeSliceId string) (*data.TimeSlice, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, timeSliceId)
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timeSlice.ID, timeSlice); err != nil {
		r.Error("error while writing time slice (%s) to cache: %s", timeSlice.ID, err)
	}
	return timeSlice, nil
}

// TimeSliceUpdate can be used to update an existing time slice
func (r *restClient) TimeSliceUpdate(ctx context.Context, timeSliceId string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, timeSliceId)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	if err := r.cacheWrite(timeSlice.ID, timeSlice); err != nil {
		r.Error("error while writing time slice (%s) to cache: %s", timeSlice.ID, err)
	}
	return timeSlice, nil
}

// TimeSliceDelete can be used to delete an existing time slice
func (r *restClient) TimeSliceDelete(ctx context.Context, timeSliceId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, timeSliceId)
	if _, err := r.doRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	if err := r.cacheDelete(timeSliceId); err != nil {
		r.Error("error while writing time slice (%s) to cache: %s", timeSliceId, err)
	}
	return nil
}

// TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (r *restClient) TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesSearch+"%s",
		r.config.Address, r.config.Port, search.ToParams())
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var timeSlices = []*data.TimeSlice{}
	if err = json.Unmarshal(bytes, &timeSlices); err != nil {
		return nil, err
	}
	return timeSlices, nil
}
