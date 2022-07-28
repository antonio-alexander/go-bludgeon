package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	rest_server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"

	"github.com/gorilla/mux"
)

type restServer struct {
	logger.Logger
	logic.Logic
	rest_server.Router
}

type Owner interface {
	Close()
}

func New(parameters ...interface{}) Owner {
	s := &restServer{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			s.Logic = p
		case rest_server.Router:
			s.Router = p
		case logger.Logger:
			s.Logger = p
		}
	}
	switch {
	case s.Logic == nil:
		panic("logic not set")
	case s.Router == nil:
		panic("router not set")
	}
	s.buildRoutes()
	return s
}

func (s *restServer) endpointTimerCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timerPartial data.TimerPartial
		var timer *data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				if timer, err = s.TimerCreate(request.Context(), timerPartial); err == nil {
					bytes, err = json.Marshal(timer)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer create -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerRead(request.Context(), id); err == nil {
			bytes, err = json.Marshal(timer)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer read -  %s", err)
		}
	}
}

func (s *restServer) endpointTimersRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var search data.TimerSearch
		var timers []*data.Timer
		var bytes []byte
		var err error

		search.FromParams(request.URL.Query())
		if timers, err = s.TimersRead(request.Context(), search); err == nil {
			bytes, err = json.Marshal(timers)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer read -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		id := idFromPath(mux.Vars(request))
		err = s.TimerDelete(request.Context(), id)
		if err = handleResponse(writer, err, nil); err != nil {
			s.Error("timer delete -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerStart() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerStart(request.Context(), id); err == nil {
			bytes, err = json.Marshal(&timer)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer start -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerStop() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timer, err = s.TimerStop(request.Context(), id); err == nil {
			bytes, err = json.Marshal(&timer)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer pause -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerSubmit() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error
		var timerPartial data.TimerPartial

		id := idFromPath(mux.Vars(request))
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			defer request.Body.Close()
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				var finishTime time.Time

				if timerPartial.Finish != nil && *timerPartial.Finish > 0 {
					finishTime = time.Unix(0, *timerPartial.Finish)
				} else {
					finishTime = time.Now()
				}
				if timer, err = s.TimerSubmit(request.Context(), id, &finishTime); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerUpdateComment() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error
		var timerPartial data.TimerPartial

		id := idFromPath(mux.Vars(request))
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			defer request.Body.Close()
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				if timerPartial.Comment == nil {
					//REVIEW: should this error be here?
					err = errors.New("No comment provided")
				} else {
					if timer, err = s.TimerUpdateComment(request.Context(), id, *timerPartial.Comment); err == nil {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerUpdateArchive() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer *data.Timer
		var bytes []byte
		var err error
		var timerPartial data.TimerPartial

		id := idFromPath(mux.Vars(request))
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			defer request.Body.Close()
			if err = json.Unmarshal(bytes, &timerPartial); err == nil {
				if timerPartial.Archived == nil {
					//REVIEW: should this error be here?
					err = errors.New("No archive provided")
				} else {
					if timer, err = s.TimerArchive(request.Context(), id, *timerPartial.Archived); err == nil {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSliceCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlicePartial data.TimeSlicePartial
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timeSlicePartial); err == nil {
				if timeSlice, err = s.TimeSliceCreate(request.Context(), timeSlicePartial); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("time slice create -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSliceRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if timeSlice, err = s.TimeSliceRead(request.Context(), id); err == nil {
			bytes, err = json.Marshal(timeSlice)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("time slice read -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSlicesRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlices []*data.TimeSlice
		var search data.TimeSliceSearch
		var bytes []byte
		var err error

		search.FromParams(request.URL.Query())
		if timeSlices, err = s.TimeSlicesRead(request.Context(), search); err == nil {
			bytes, err = json.Marshal(timeSlices)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("time slices read -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSliceUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlicePartial data.TimeSlicePartial
		var timeSlice *data.TimeSlice
		var bytes []byte
		var err error

		id := idFromPath(mux.Vars(request))
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &timeSlicePartial); err == nil {
				if timeSlice, err = s.TimeSliceUpdate(request.Context(), id, timeSlicePartial); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("timer slice create -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSliceDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		id := idFromPath(mux.Vars(request))
		err = s.TimeSliceDelete(request.Context(), id)
		if err = handleResponse(writer, err, nil); err != nil {
			s.Error("timer delete -  %s", err)
		}
	}
}

func (s *restServer) buildRoutes() {
	for _, route := range []rest_server.HandleFuncConfig{
		//timer
		{Route: data.RouteTimers, Method: http.MethodPost, HandleFx: s.endpointTimerCreate()},
		{Route: data.RouteTimersSearch, Method: http.MethodGet, HandleFx: s.endpointTimersRead()},
		{Route: data.RouteTimersID, Method: http.MethodGet, HandleFx: s.endpointTimerRead()},
		{Route: data.RouteTimersID, Method: http.MethodDelete, HandleFx: s.endpointTimerDelete()},
		{Route: data.RouteTimersIDStart, Method: http.MethodPut, HandleFx: s.endpointTimerStart()},
		{Route: data.RouteTimersIDStop, Method: http.MethodPut, HandleFx: s.endpointTimerStop()},
		{Route: data.RouteTimersIDSubmit, Method: http.MethodPut, HandleFx: s.endpointTimerSubmit()},
		{Route: data.RouteTimersIDComment, Method: http.MethodPut, HandleFx: s.endpointTimerUpdateComment()},
		{Route: data.RouteTimersIDArchive, Method: http.MethodPut, HandleFx: s.endpointTimerUpdateArchive()},
		//time slice
		{Route: data.RouteTimeSlices, Method: http.MethodPost, HandleFx: s.endpointTimeSliceCreate()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodGet, HandleFx: s.endpointTimeSliceRead()},
		{Route: data.RouteTimeSlices, Method: http.MethodGet, HandleFx: s.endpointTimeSlicesRead()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodPut, HandleFx: s.endpointTimeSliceUpdate()},
		{Route: data.RouteTimeSlicesID, Method: http.MethodDelete, HandleFx: s.endpointTimeSliceDelete()},
	} {
		s.Router.HandleFunc(route)
	}
}

func (s *restServer) Close() {
	//KIM: we have the option of doing nothing here since upon close this pointer
	// shouldn't be re-used, and we can't ensure that the endpoints aren't being
	// called again (to prevent a panic)
}
