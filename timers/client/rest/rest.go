package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/logic"
)

const urif string = "http://%s:%s%s"

type rest struct {
	client.Client
	logger.Logger
	config *client.Configuration
}

//Client describes operations that can be done with a rest
// client
type Client interface {
	logic.Logic

	//Initialize can be used to configure and start the business logic
	// of the underlying pointer
	Initialize(config *client.Configuration) error
}

//New will create a populated instance of the rest client
// if Configuration is provided as a parameter, it will
// also attempt to initialize (and panic on error)
func New(parameters ...interface{}) Client {
	var config *client.Configuration
	r := &rest{
		Client: client.New(parameters...),
	}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *client.Configuration:
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
func (r *rest) Initialize(config *client.Configuration) error {
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
func (r *rest) TimerCreate(timerPartial data.TimerPartial) (*data.Timer, error) {
	bytes, err := json.Marshal(&timerPartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteTimers)
	bytes, err = r.DoRequest(uri, http.MethodPost, bytes)
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
func (r *rest) TimerRead(id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
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
func (r *rest) TimersRead(search data.TimerSearch) ([]*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteTimersSearch+search.ToParams())
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
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
func (r *rest) TimerStart(id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDStartf, id))
	bytes, err := r.DoRequest(uri, http.MethodPut, nil)
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
func (r *rest) TimerStop(id string) (*data.Timer, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDStopf, id))
	bytes, err := r.DoRequest(uri, http.MethodPut, nil)
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
func (r *rest) TimerUpdateComment(id, comment string) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Comment: &comment,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDCommentf, id))
	bytes, err = r.DoRequest(uri, http.MethodPut, bytes)
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
func (r *rest) TimerArchive(id string, archive bool) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Archived: &archive,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDArchivef, id))
	bytes, err = r.DoRequest(uri, http.MethodPut, bytes)
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
func (r *rest) TimerDelete(id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
	if _, err := r.DoRequest(uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//TimerSubmit can be used to stop a timer and set completed to true
func (r *rest) TimerSubmit(id string, finishTime *time.Time) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.Contract{
		Finish: finishTime.UnixNano(),
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDSubmitf, id))
	bytes, err = r.DoRequest(uri, http.MethodPut, bytes)
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
func (r *rest) TimeSliceCreate(timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port, data.RouteTimeSlices)
	bytes, err = r.DoRequest(uri, http.MethodPost, bytes)
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
func (r *rest) TimeSliceRead(id string) (*data.TimeSlice, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
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
func (r *rest) TimeSliceUpdate(id string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	bytes, err := json.Marshal(&timeSlicePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	bytes, err = r.DoRequest(uri, http.MethodPut, bytes)
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
func (r *rest) TimeSliceDelete(id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimeSlicesIDf, id))
	if _, err := r.DoRequest(uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

//TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (r *rest) TimeSlicesRead(search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		data.RouteTimeSlicesSearch+search.ToParams())
	bytes, err := r.DoRequest(uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var timeSlices = []*data.TimeSlice{}
	if err = json.Unmarshal(bytes, &timeSlices); err != nil {
		return nil, err
	}
	return timeSlices, nil
}
