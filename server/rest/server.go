package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
	server "github.com/antonio-alexander/go-bludgeon/server"

	"github.com/pkg/errors"
)

type restServer struct {
	sync.RWMutex                 //mutex for threadsafe operations
	sync.WaitGroup               //waitgroup to track go routines
	logic.Logic                  //
	data.Logger                  //
	config         Configuration //configuration
	stopper        chan struct{} //stopper to stop goRoutines
	server         rest.Server   //
	started        bool          //whether or not the business logic has starting
}

func New(logger data.Logger, logic logic.Logic) interface {
	Owner
	server.Owner
	logic.Logic
} {
	return &restServer{
		Logic:  logic,
		Logger: logger,
		server: rest.New(logger),
	}
}

//TimerCreate
func (s *restServer) endpointTimerCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer data.Timer
		var bytes []byte
		var err error

		//attempt to execute the timer create
		if timer, err = s.TimerCreate(); err == nil {
			bytes, err = json.Marshal(timer)
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerCreate"))
		}
	}
}

func (s *restServer) endpointTimerRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = s.TimerRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerRead"))
		}
	}
}

//TimerUpdate
func (s *restServer) endpointTimerUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timer data.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = s.TimerUpdate(contract.Timer); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerUpdate"))
		}
	}
}

//TimerDelete
func (s *restServer) endpointTimerDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				err = s.TimerDelete(contract.ID)
			}
		}
		//handle errors
		if err = handleResponse(writer, err, nil); err != nil {
			s.Error(errors.Wrap(err, "TimerDelete"))
		}
	}
}

//TimerStart
func (s *restServer) endpointTimerStart() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var timeStart time.Time
		var timer data.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
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
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerStart"))
		}
	}
}

//TimerPause
func (s *restServer) endpointTimerPause() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var pauseTime time.Time
		var timer data.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
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
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerPause"))
		}
	}
}

//TimerSubmit
func (s *restServer) endpointTimerSubmit() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract data.Contract
		var finishTime time.Time
		var timer data.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
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
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimerSubmit"))
		}
	}
}

//TimeSliceRead
func (s *restServer) endpointTimeSliceRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice data.TimeSlice
		var contract data.Contract
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timeSlice, err = s.TimeSliceRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error(errors.Wrap(err, "TimeSliceRead"))
		}
	}
}

//buildRoutes will create all the routes and their functions to execute when received
func (s *restServer) BuildRoutes() error {
	routes := []rest.HandleFuncConfig{
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
	}
	return s.server.BuildRoutes(routes)
}

func (s *restServer) Start(config *Configuration) (err error) {
	s.Lock()
	defer s.Unlock()

	if s.started {
		return errors.New(ErrStarted)
	}
	s.config = *config
	if err = s.BuildRoutes(); err != nil {
		return
	}
	s.stopper = make(chan struct{})
	if err = s.server.Start(s.config.Address, s.config.Port); err != nil {
		return
	}
	s.started = true
	s.Info("bludgeon server started, listening on %s:%s", s.config.Address, s.config.Port)

	return
}

func (s *restServer) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	if !s.started {
		err = errors.New(ErrNotStarted)

		return
	}
	s.server.Stop()
	close(s.stopper)
	s.Wait()
	s.started = false
	s.Info("bludgeon server stopped")

	return
}
