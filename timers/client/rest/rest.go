package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/pkg/errors"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
)

type restClient struct {
	logger.Logger
	client interface {
		internal.Configurer
		internal.Parameterizer
		restclient.Client
	}
	config     *Configuration
	configured bool
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
	r.client.SetParameters(parameters...)
}

func (r *restClient) SetUtilities(parameters ...interface{}) {
	r.client.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			r.Logger = p
		}
	}
}

func (r *restClient) Configure(items ...interface{}) error {
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
	if err := r.client.Configure(&c.Configuration); err != nil {
		return err
	}
	r.config = c
	r.configured = true
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (r *restClient) Initialize() error {
	return nil
}

func (r *restClient) Shutdown() {}

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
	return timer, nil
}

// TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (r *restClient) TimerRead(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDf,
		r.config.Address, r.config.Port, id)
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
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
	return timer, nil
}

// TimerDelete can be used to delete a timer if it exists
func (r *restClient) TimerDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDf,
		r.config.Address, r.config.Port, id)
	if _, err := r.doRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

// TimerSubmit can be used to stop a timer and set completed to true
func (r *restClient) TimerSubmit(ctx context.Context, id string, finishTime int64) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.Contract{
		Finish: finishTime,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimersIDSubmitf,
		r.config.Address, r.config.Port, id)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
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
	return timeSlice, nil
}

// TimeSliceRead can be used to read an existing time slice
func (r *restClient) TimeSliceRead(ctx context.Context, id string) (*data.TimeSlice, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, id)
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

// TimeSliceUpdate can be used to update an existing time slice
func (r *restClient) TimeSliceUpdate(ctx context.Context, id string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, id)
	bytes, err = r.doRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

// TimeSliceDelete can be used to delete an existing time slice
func (r *restClient) TimeSliceDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteTimeSlicesIDf,
		r.config.Address, r.config.Port, id)
	if _, err := r.doRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
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
