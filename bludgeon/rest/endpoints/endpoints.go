package bludgeonrestendpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest"

	"github.com/pkg/errors"
)

//BuildRoutes will create all the routes and their functions to execute when received
func BuildRoutes(logger bludgeon.Logger, functional interface{}) (routes []rest.HandleFuncConfig) {
	//get the routes for the functional timer, and then the timeslice
	if f, ok := functional.(bludgeon.FunctionalTimer); ok {
		routes = append(routes, []rest.HandleFuncConfig{
			{Route: rest.RouteTimerCreate, Method: POST, HandleFx: TimerCreate(logger, f)},
			{Route: rest.RouteTimerRead, Method: POST, HandleFx: TimerRead(logger, f)},
			{Route: rest.RouteTimerUpdate, Method: POST, HandleFx: TimerUpdate(logger, f)},
			{Route: rest.RouteTimerDelete, Method: POST, HandleFx: TimerDelete(logger, f)},
			{Route: rest.RouteTimerStart, Method: POST, HandleFx: TimerStart(logger, f)},
			{Route: rest.RouteTimerPause, Method: POST, HandleFx: TimerPause(logger, f)},
			{Route: rest.RouteTimerSubmit, Method: POST, HandleFx: TimerSubmit(logger, f)},
		}...)
	}
	if f, ok := functional.(bludgeon.FunctionalTimeSlice); ok {
		routes = append(routes, []rest.HandleFuncConfig{
			{Route: rest.RouteTimeSliceRead, Method: POST, HandleFx: TimeSliceRead(logger, f)},
		}...)
	}
	// if f, ok := functional.(bludgeon.FunctionalOwner); ok {
	// 	routes = append(routes, []rest.HandleFuncConfig{
	// 		{Route: rest.RouteStop, Method: POST, HandleFx: Stop(logger, f)},
	// 	}...)
	// }

	return
}

// //ServerStop
// func Stop(l bludgeon.Logger, f 	bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
// 	return func(writer http.ResponseWriter, request *http.Request) {
// 		var token bludgeon.Token
// 		var bytes []byte
// 		var err error

// 		//get token
// 		if token, err = getToken(request); err == nil {
// 			//attempt to execute the timer create
// 			_, err = f.Stop(bludgeon.CommandServerStop, nil, token)
// 		}
// 		//handle errors
// 		if err = handleResponse(writer, err, bytes); err != nil {
// 			err = errors.Wrap(err, "Stop")
// 			l.Error(err)
// 		}
// 	}
// }

//TimerCreate
func TimerCreate(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//attempt to execute the timer create
		if timer, err = f.TimerCreate(); err == nil {
			bytes, err = json.Marshal(timer)
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerCreate")
			l.Error(err)
		}
	}
}

//TimerRead
func TimerRead(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = f.TimerRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerRead")
			l.Error(err)
		}
	}
}

//TimerUpdate
func TimerUpdate(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = f.TimerUpdate(contract.Timer); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerUpdate")
			l.Error(err)
		}
	}
}

//TimerDelete
func TimerDelete(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				err = f.TimerDelete(contract.ID)
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerDelete")
			l.Error(err)
		}
	}
}

//TimerStart
func TimerStart(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = f.TimerStart(contract.ID, time.Unix(0, contract.StartTime)); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerStart")
			l.Error(err)
		}
	}
}

//TimerPause
func TimerPause(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = f.TimerPause(contract.ID, time.Unix(0, contract.PauseTime)); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerPause")
			l.Error(err)
		}
	}
}

//TimerSubmit
func TimerSubmit(l bludgeon.Logger, f bludgeon.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.Contract
		var timer bludgeon.Timer
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timer, err = f.TimerSubmit(contract.ID, time.Unix(0, contract.FinishTime)); err == nil {
					bytes, err = json.Marshal(&timer)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimerSubmit")
			l.Error(err)
		}
	}
}

//TimeSliceRead
func TimeSliceRead(l bludgeon.Logger, f bludgeon.FunctionalTimeSlice) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice bludgeon.TimeSlice
		var contract rest.Contract
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if timeSlice, err = f.TimeSliceRead(contract.ID); err == nil {
					bytes, err = json.Marshal(timeSlice)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "TimeSliceRead")
			l.Error(err)
		}
	}
}
