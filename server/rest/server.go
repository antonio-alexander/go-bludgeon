package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
	rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	logic "github.com/antonio-alexander/go-bludgeon/logic"

	"github.com/pkg/errors"
)

type server struct {
	sync.RWMutex                 //mutex for threadsafe operations
	sync.WaitGroup               //waitgroup to track go routines
	logic.Logic                  //
	common.Logger                //
	config         Configuration //configuration
	stopper        chan struct{} //stopper to stop goRoutines
	server         rest.Server   //
	started        bool          //whether or not the business logic has starting
}

func New(logger common.Logger, logic logic.Logic) interface {
	Owner
	Manage
	logic.Logic
} {
	return &server{
		Logic:  logic,
		Logger: logger,
		server: rest.New(logger),
	}
}

func (s *server) launchPurge() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		tPurge := time.NewTicker(30 * time.Second)
		defer tPurge.Stop()
		close(started)
		for {
			select {
			case <-tPurge.C:
				// if err := s.TokenPurge(); err != nil {
				// 	fmt.Printf(err)
				// }
			case <-s.stopper:
				return
			}
		}
	}()
	<-started
}

//TimerCreate
func (s *server) endpointTimerCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer common.Timer
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

//TimerRead
func (s *server) endpointTimerRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
func (s *server) endpointTimerUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
func (s *server) endpointTimerDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
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
func (s *server) endpointTimerStart() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = s.TimerStart(contract.ID, time.Unix(0, contract.StartTime)); err == nil {
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
func (s *server) endpointTimerPause() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = s.TimerPause(contract.ID, time.Unix(0, contract.PauseTime)); err == nil {
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
func (s *server) endpointTimerSubmit() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = s.TimerSubmit(contract.ID, time.Unix(0, contract.FinishTime)); err == nil {
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
func (s *server) endpointTimeSliceRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice common.TimeSlice
		var contract common.Contract
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
func (s *server) buildRoutes() []rest.HandleFuncConfig {
	//REVIEW: in the future when we add tokens, we'll need to create some way to check tokens for
	// certain functions, we may need to implement varratics to add support for tokens for the
	// server actions, may not be able to re-use the existing endpoints
	return []rest.HandleFuncConfig{
		//timer
		{Route: common.RouteTimerCreate, Method: POST, HandleFx: s.endpointTimerCreate()},
		{Route: common.RouteTimerRead, Method: POST, HandleFx: s.endpointTimerRead()},
		{Route: common.RouteTimerUpdate, Method: POST, HandleFx: s.endpointTimerUpdate()},
		{Route: common.RouteTimerDelete, Method: POST, HandleFx: s.endpointTimerDelete()},
		{Route: common.RouteTimerStart, Method: POST, HandleFx: s.endpointTimerStart()},
		{Route: common.RouteTimerPause, Method: POST, HandleFx: s.endpointTimerPause()},
		{Route: common.RouteTimerSubmit, Method: POST, HandleFx: s.endpointTimerSubmit()},
		//time slice
		{Route: common.RouteTimeSliceRead, Method: POST, HandleFx: s.endpointTimeSliceRead()},
	}
}

func (s *server) Close() {
	s.Lock()
	defer s.Unlock()

	//use the remote and/or storage functions to:
	// attempt to synchronize the current values
	// attempt to serialize what's remaining

	//set internal configuration to default
	s.config = Configuration{}
	//set internal pointers to nil
	s.server = nil
	s.Logic = nil
}

func (s *server) Start(conf Configuration) (err error) {
	s.Lock()
	defer s.Unlock()

	//check if already started, store configuration, create stopper,
	// slaunch go routineet started to true
	if s.started {
		return errors.New(ErrStarted)
	}
	s.config = conf
	s.stopper = make(chan struct{})
	routes := s.buildRoutes()
	if err = s.server.BuildRoutes(routes); err != nil {
		return
	}
	s.launchPurge()
	s.started = true

	return
}

func (s *server) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	//check if not started
	//close the stopper
	//set started flag to false
	if !s.started {
		err = errors.New(ErrNotStarted)

		return
	}
	s.server.Stop()
	close(s.stopper)
	s.Wait()
	s.started = false

	return
}
