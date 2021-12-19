package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
	server "github.com/antonio-alexander/go-bludgeon/server"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type restServer struct {
	sync.RWMutex
	sync.WaitGroup
	logic.Logic
	data.Logger
	config  Configuration
	router  *mux.Router
	server  *http.Server
	stopper chan struct{}
	started bool
}

func New(logger data.Logger, logic logic.Logic) interface {
	Owner
	server.Owner
	logic.Logic
} {
	router := mux.NewRouter()
	return &restServer{
		Logic:  logic,
		Logger: logger,
		router: router,
		server: &http.Server{
			Handler: router,
		},
	}
}

func (s *restServer) handleResponse(writer http.ResponseWriter, errIn error, bytes []byte) error {
	if errIn != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err := writer.Write([]byte(errIn.Error()))
		return err
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := writer.Write(bytes)
	return err
}

func (s *restServer) endpointTimerCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer data.Timer
		var bytes []byte
		var err error

		if timer, err = s.TimerCreate(); err == nil {
			bytes, err = json.Marshal(timer)
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer create -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				if timer, err = s.TimerRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer read -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				if timer, err = s.TimerUpdate(contract.Timer); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer update -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				err = s.TimerDelete(contract.ID)
			}
		}
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error("timer delete -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerStart() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				var timeStart time.Time

				if contract.StartTime > 0 {
					timeStart = time.Unix(0, contract.StartTime)
				} else {
					timeStart = time.Now()
				}
				if timer, err = s.TimerStart(contract.ID, timeStart); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer start -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerPause() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				var pauseTime time.Time

				if contract.PauseTime > 0 {
					pauseTime = time.Unix(0, contract.PauseTime)
				} else {
					pauseTime = time.Now()
				}
				if timer, err = s.TimerPause(contract.ID, pauseTime); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer pause -  %s", err)
		}
	}
}

func (s *restServer) endpointTimerSubmit() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				var finishTime time.Time

				if contract.FinishTime > 0 {
					finishTime = time.Unix(0, contract.FinishTime)
				} else {
					finishTime = time.Now()
				}
				if timer, err = s.TimerSubmit(contract.ID, finishTime); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("timer submit -  %s", err)
		}
	}
}

func (s *restServer) endpointTimeSliceRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice data.TimeSlice
		var contract data.Contract
		var bytes []byte
		var err error

		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				if timeSlice, err = s.TimeSliceRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		if err = s.handleResponse(writer, err, bytes); err != nil {
			s.Error("time slice read -  %s", err)
		}
	}
}

func (s *restServer) buildRoutes() {
	for _, route := range []handleFuncConfig{
		//timer
		{Route: data.RouteTimerCreate, Method: POST, HandleFx: s.endpointTimerCreate()},
		{Route: data.RouteTimerRead, Method: POST, HandleFx: s.endpointTimerRead()},
		{Route: data.RouteTimerUpdate, Method: POST, HandleFx: s.endpointTimerUpdate()},
		{Route: data.RouteTimerDelete, Method: POST, HandleFx: s.endpointTimerDelete()},
		{Route: data.RouteTimerStart, Method: POST, HandleFx: s.endpointTimerStart()},
		{Route: data.RouteTimerPause, Method: POST, HandleFx: s.endpointTimerPause()},
		{Route: data.RouteTimerSubmit, Method: POST, HandleFx: s.endpointTimerSubmit()},
		//time slice
		{Route: data.RouteTimeSliceRead, Method: POST, HandleFx: s.endpointTimeSliceRead()},
	} {
		s.router.HandleFunc(route.Route, route.HandleFx).Methods(route.Method)
	}
}

func (s *restServer) launchServer() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		close(started)
		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.Debug("Httpserver: ListenAndServe() error: %s", err)
			}
		}
		//REVIEW: Do we need to account for a situation where the rest server kills itself
		// unexepctedly?
	}()
	<-started
}

func (s *restServer) Start(config *Configuration) (err error) {
	s.Lock()
	defer s.Unlock()

	if s.started {
		return errors.New(ErrStarted)
	}
	s.config = *config
	s.buildRoutes()
	s.stopper = make(chan struct{})
	s.server.Addr = fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
	s.launchServer()
	s.started = true
	s.Info("started, listening on %s:%s", s.config.Address, s.config.Port)

	return
}

func (s *restServer) Stop() {
	s.Lock()
	defer s.Unlock()

	if !s.started {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), ConfigShutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.Error("shutting down - %s", err)
	}
	close(s.stopper)
	s.Wait()
	s.started = false
	s.Info("stopped")
}
