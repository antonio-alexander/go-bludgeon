package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
)

const urif string = "http://%s:%s%s"

type rest struct {
	restclient.Client
	logger.Logger
	config *restclient.Configuration
}

//Client describes operations that can be done with a rest
// client
type Client interface {
	logic.Logic

	//Initialize can be used to configure and start the business logic
	// of the underlying pointer
	Initialize(config *restclient.Configuration) error
}

//New will create a populated instance of the rest client
// if Configuration is provided as a parameter, it will
// also attempt to initialize (and panic on error)
func New(parameters ...interface{}) Client {
	var config *restclient.Configuration
	r := &rest{
		Client: restclient.New(parameters...),
	}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *restclient.Configuration:
			config = p
		case logger.Logger:
			r.Logger = p
		}
	}
	if config != nil {
		if err := r.Initialize(config); err != nil {
			panic(err)
		}
	}
	return r
}

//Initialize can be used to configure and start the business logic
// of the underlying pointer
func (r *rest) Initialize(config *restclient.Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if err := r.Client.Initialize(config); err != nil {
		return err
	}
	r.config = config
	return nil
}

//TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (r *rest) TimerCreate(ctx context.Context, timerPartial data.TimerPartial) (*data.Timer, error) {
	bytes, err := json.Marshal(&timerPartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteTimers)
	bytes, err = r.DoRequest(ctx, uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (r *rest) TimerRead(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimersRead can be used to read one or more timers depending
// on search values provided
func (r *rest) TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteTimersSearch+search.ToParams())
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var timers = []*data.Timer{}
	if err = json.Unmarshal(bytes, &timers); err != nil {
		return nil, err
	}
	return timers, nil
}

//TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (r *rest) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDStartf, id))
	bytes, err := r.DoRequest(ctx, uri, http.MethodPut, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (r *rest) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDStopf, id))
	bytes, err := r.DoRequest(ctx, uri, http.MethodPut, nil)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerUpdateCommnet will only update the comment for timer with
// the provided id
func (r *rest) TimerUpdateComment(ctx context.Context, id, comment string) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Comment: &comment,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDCommentf, id))
	bytes, err = r.DoRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerArchive will only update the archive for timer with
// the provided id
func (r *rest) TimerArchive(ctx context.Context, id string, archive bool) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Archived: &archive,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDArchivef, id))
	bytes, err = r.DoRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerDelete can be used to delete a timer if it exists
func (r *rest) TimerDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
	if _, err := r.DoRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//TimerSubmit can be used to stop a timer and set completed to true
func (r *rest) TimerSubmit(ctx context.Context, id string, finishTime *time.Time) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.Contract{
		Finish: finishTime.UnixNano(),
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDSubmitf, id))
	bytes, err = r.DoRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timer := new(data.Timer)
	if err = json.Unmarshal(bytes, timer); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimeSliceCreate can be used to create a single time
// slice
func (r *rest) TimeSliceCreate(ctx context.Context, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteTimeSlices)
	bytes, err = r.DoRequest(ctx, uri, http.MethodPost, bytes)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

//TimeSliceRead can be used to read an existing time slice
func (r *rest) TimeSliceRead(ctx context.Context, id string) (*data.TimeSlice, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

//TimeSliceUpdate can be used to update an existing time slice
func (r *rest) TimeSliceUpdate(ctx context.Context, id string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	bytes, err = r.DoRequest(ctx, uri, http.MethodPut, bytes)
	if err != nil {
		return nil, err
	}
	timeSlice := new(data.TimeSlice)
	if err = json.Unmarshal(bytes, timeSlice); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

//TimeSliceDelete can be used to delete an existing time slice
func (r *rest) TimeSliceDelete(ctx context.Context, id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	if _, err := r.DoRequest(ctx, uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (r *rest) TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteTimeSlicesSearch+search.ToParams())
	bytes, err := r.DoRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var timeSlices = []*data.TimeSlice{}
	if err = json.Unmarshal(bytes, &timeSlices); err != nil {
		return nil, err
	}
	return timeSlices, nil
}
