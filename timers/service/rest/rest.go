package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/gorilla/mux"
)

type restService struct {
	logger.Logger
	logic.Logic
	ctx context.Context
}

func New() interface {
	internal.Parameterizer
	internal_rest.RouteBuilder
} {
	return &restService{
		Logger: logger.NewNullLogger(),
		ctx:    context.Background(),
	}
}

func (s *restService) handleResponse(writer http.ResponseWriter, err error, bytes []byte) error {
	if err != nil {
		var e internal_errors.Error

		switch {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, meta.ErrTimerNotFound):
			writer.WriteHeader(http.StatusNotFound)
		case errors.Is(err, meta.ErrTimerNotUpdated):
			writer.WriteHeader(http.StatusNotModified)
		case errors.Is(err, meta.ErrTimerConflictCreate) || errors.Is(err, meta.ErrTimerConflictUpdate):
			writer.WriteHeader(http.StatusConflict)
		}
		switch i := err.(type) {
		case internal_errors.Error:
			e = i
		default:
			e = internal_errors.New(err.Error())
		}
		bytes, err = json.Marshal(&e)
		if err != nil {
			return err
		}
		_, err = writer.Write(bytes)
		return err
	}
	if bytes == nil {
		writer.WriteHeader(http.StatusNoContent)
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = writer.Write(bytes)
	return err
}

func (s *restService) endpointTimerCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timerPartial data.TimerPartial
		var timer *data.Timer
		var bytes []byte
		var err error

		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				if timer, err = s.TimerCreate(request.Context(), timerPartial); err == nil {
					bytes, err = json.Marshal(timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer create -  %s", err)
		}
	}
}

func (s *restService) endpointTimerRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerRead(request.Context(), id); err == nil {
			bytes, err = json.Marshal(timer)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer read -  %s", err)
		}
	}
}

func (s *restService) endpointTimersRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var search data.TimerSearch
		var timers []*data.Timer
		var bytes []byte
		var err error

		search.FromParams(request.URL.Query())
		if timers, err = s.TimersRead(request.Context(), search); err == nil {
			bytes, err = json.Marshal(timers)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer read -  %s", err)
		}
	}
}

func (s *restService) endpointTimerDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		id := idFromPath(mux.Vars(request))
		err = s.TimerDelete(request.Context(), id)
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error("timer delete -  %s", err)
		}
	}
}

func (s *restService) endpointTimerStart() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerStart(request.Context(), id); err == nil {
			bytes, err = json.Marshal(&timer)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer start -  %s", err)
		}
	}
}

func (s *restService) endpointTimerStop() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerStop(request.Context(), id); err == nil {
			bytes, err = json.Marshal(&timer)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer pause -  %s", err)
		}
	}
}

func (s *restService) endpointTimerSubmit() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error
		var timerPartial data.TimerPartial

		id := idFromPath(mux.Vars(request))
		if bytes, err = io.ReadAll(request.Body); err == nil {
			defer request.Body.Close()
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				var finishTime time.Time

				if timerPartial.Finish != nil && *timerPartial.Finish > 0 {
					finishTime = time.Unix(0, *timerPartial.Finish)
				} else {
					finishTime = time.Now()
				}
				if timer, err = s.TimerSubmit(request.Context(), id, finishTime.UnixNano()); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restService) endpointTimerUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error
		var timerPartial data.TimerPartial

		id := idFromPath(mux.Vars(request))
		if bytes, err = io.ReadAll(request.Body); err == nil {
			defer request.Body.Close()
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				if timer, err = s.TimerUpdate(request.Context(), id, timerPartial); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restService) endpointTimeSliceCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlicePartial data.TimeSlicePartial
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timeSlicePartial); err == nil {
				if timeSlice, err = s.TimeSliceCreate(request.Context(), timeSlicePartial); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("time slice create -  %s", err)
		}
	}
}

func (s *restService) endpointTimeSliceRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timeSlice, err = s.TimeSliceRead(request.Context(), id); err == nil {
			bytes, err = json.Marshal(timeSlice)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("time slice read -  %s", err)
		}
	}
}

func (s *restService) endpointTimeSlicesRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlices []*data.TimeSlice
		var search data.TimeSliceSearch
		var bytes []byte
		var err error

		search.FromParams(request.URL.Query())
		if timeSlices, err = s.TimeSlicesRead(request.Context(), search); err == nil {
			bytes, err = json.Marshal(timeSlices)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("time slices read -  %s", err)
		}
	}
}

func (s *restService) endpointTimeSliceUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlicePartial data.TimeSlicePartial
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timeSlicePartial); err == nil {
				if timeSlice, err = s.TimeSliceUpdate(request.Context(), id, timeSlicePartial); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer slice create -  %s", err)
		}
	}
}

func (s *restService) endpointTimeSliceDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		id := idFromPath(mux.Vars(request))
		err = s.TimeSliceDelete(request.Context(), id)
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error("timer delete -  %s", err)
		}
	}
}

func (s *restService) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *restService) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			s.Logic = p
		case context.Context:
			s.ctx = p
		}
	}
	switch {
	case s.Logic == nil:
		panic("logic not set")
	}
}

func (s *restService) BuildRoutes() []internal_rest.HandleFuncConfig {
	return []internal_rest.HandleFuncConfig{
		//timer
		{Route: data.RouteTimers, Method: http.MethodPost, HandleFx: s.endpointTimerCreate()},
		{Route: data.RouteTimersSearch, Method: http.MethodGet, HandleFx: s.endpointTimersRead()},
		{Route: data.RouteTimersID, Method: http.MethodGet, HandleFx: s.endpointTimerRead()},
		{Route: data.RouteTimersID, Method: http.MethodPut, HandleFx: s.endpointTimerUpdate()},
		{Route: data.RouteTimersID, Method: http.MethodDelete, HandleFx: s.endpointTimerDelete()},
		{Route: data.RouteTimersIDStart, Method: http.MethodPut, HandleFx: s.endpointTimerStart()},
		{Route: data.RouteTimersIDStop, Method: http.MethodPut, HandleFx: s.endpointTimerStop()},
		{Route: data.RouteTimersIDSubmit, Method: http.MethodPut, HandleFx: s.endpointTimerSubmit()},
		//time slice
		{Route: data.RouteTimeSlices, Method: http.MethodPost, HandleFx: s.endpointTimeSliceCreate()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodGet, HandleFx: s.endpointTimeSliceRead()},
		{Route: data.RouteTimeSlices, Method: http.MethodGet, HandleFx: s.endpointTimeSlicesRead()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodPut, HandleFx: s.endpointTimeSliceUpdate()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodDelete, HandleFx: s.endpointTimeSliceDelete()},
	}
}
