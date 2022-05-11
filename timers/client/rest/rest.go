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

type Owner interface {
	Initialize(config *client.Configuration) error
}

func New(parameters ...interface{}) interface {
	logic.Logic
	Owner
} {
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

func (r *rest) TimerUpdateComment(id, comment string) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Comment: &comment,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
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

func (r *rest) TimerArchive(id string, archive bool) (*data.Timer, error) {
	bytes, err := json.Marshal(&data.TimerPartial{
		Archived: &archive,
	})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
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

func (r *rest) TimerDelete(id string) error {
	uri := fmt.Sprintf(urif, r.config.Address, r.config.Port,
		fmt.Sprintf(data.RouteTimersIDf, id))
	if _, err := r.DoRequest(uri, http.MethodDelete, nil); err != nil {
		return err
	}
	return nil
}

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
