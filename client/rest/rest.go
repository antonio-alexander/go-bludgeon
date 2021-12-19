package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
)

type rest struct {
	*http.Client
	config Configuration
}

func New(config Configuration) interface {
	logic.Logic
} {
	return &rest{
		config: config,
		Client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

//doRequest
func (r *rest) doRequest(uri, method string, data []byte) ([]byte, error) {
	request, err := http.NewRequest(method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	response, err := r.Do(request)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		if len(data) > 0 {
			return nil, errors.New(string(data))
		}
		return nil, fmt.Errorf("failure: %d", response.StatusCode)
	}
	return data, nil
}

func (r *rest) TimerCreate() (data.Timer, error) {
	var timer data.Timer

	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerCreate)
	bytes, err := r.doRequest(uri, POST, nil)
	if err != nil {
		return data.Timer{}, err
	}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}
	return timer, nil
}

//TimerRead
func (r *rest) TimerRead(id string) (data.Timer, error) {
	var timer data.Timer

	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerRead)
	bytes, err := json.Marshal(&data.Contract{
		ID: id,
	})
	if err != nil {
		return data.Timer{}, err
	}
	bytes, err = r.doRequest(uri, POST, bytes)
	if err != nil {
		return data.Timer{}, err
	}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}

	return timer, nil
}

//
func (r *rest) TimerUpdate(t data.Timer) (data.Timer, error) {
	bytes, err := json.Marshal(&data.Contract{
		Timer: t,
	})
	if err != nil {
		return data.Timer{}, err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerUpdate)
	bytes, err = r.doRequest(uri, POST, bytes)
	if err != nil {
		return data.Timer{}, err
	}
	timer := data.Timer{}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}
	return timer, nil
}

//
func (r *rest) TimerDelete(id string) error {
	bytes, err := json.Marshal(&data.Contract{
		ID: id,
	})
	if err != nil {
		return err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerDelete)
	if _, err = r.doRequest(uri, POST, bytes); err != nil {
		return err
	}
	return nil
}

//
func (r *rest) TimerStart(timerID string, startTime time.Time) (data.Timer, error) {
	var timer data.Timer

	bytes, err := json.Marshal(&data.Contract{
		ID:        timerID,
		StartTime: startTime.UnixNano(),
	})
	if err != nil {
		return data.Timer{}, err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerStart)
	bytes, err = r.doRequest(uri, POST, bytes)
	if err != nil {
		return data.Timer{}, err
	}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}
	return timer, nil
}

//
func (r *rest) TimerPause(timerID string, pauseTime time.Time) (data.Timer, error) {
	var timer data.Timer

	bytes, err := json.Marshal(&data.Contract{
		ID:        timerID,
		PauseTime: pauseTime.UnixNano(),
	})
	if err != nil {
		return data.Timer{}, err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerPause)
	if bytes, err = r.doRequest(uri, POST, bytes); err != nil {
		return data.Timer{}, err
	}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}
	return timer, nil
}

//
func (r *rest) TimerSubmit(timerID string, finishTime time.Time) (data.Timer, error) {
	var timer data.Timer

	bytes, err := json.Marshal(&data.Contract{
		ID:         timerID,
		FinishTime: finishTime.UnixNano(),
	})
	if err != nil {
		return data.Timer{}, err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimerSubmit)
	if bytes, err = r.doRequest(uri, POST, bytes); err != nil {
		return data.Timer{}, err
	}
	if err = json.Unmarshal(bytes, &timer); err != nil {
		return data.Timer{}, err
	}
	return timer, nil
}

func (r *rest) TimeSliceRead(id string) (data.TimeSlice, error) {
	var timeSlice data.TimeSlice

	bytes, err := json.Marshal(&data.Contract{
		ID: id,
	})
	if err != nil {
		return data.TimeSlice{}, err
	}
	uri := fmt.Sprintf(URIf, r.config.Address, r.config.Port, data.RouteTimeSliceRead)
	if bytes, err = r.doRequest(uri, POST, bytes); err != nil {
		return data.TimeSlice{}, err
	}
	if err = json.Unmarshal(bytes, &timeSlice); err != nil {
		return data.TimeSlice{}, err
	}
	return timeSlice, nil
}
