package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"

	"github.com/pkg/errors"
)

//Stop
func Stop(l common.Logger, f common.FunctionalManage) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		//attempt to execute the timer create
		err = f.Stop()
		//handle errors
		if err = handleResponse(writer, err, nil); err != nil {
			l.Error(errors.Wrap(err, "Stop"))
		}
	}
}

//TimerCreate
func TimerCreate(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timer common.Timer
		var bytes []byte
		var err error

		//attempt to execute the timer create
		if timer, err = f.TimerCreate(); err == nil {
			bytes, err = json.Marshal(timer)
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			l.Error(errors.Wrap(err, "TimerCreate"))
		}
	}
}

//TimerRead
func TimerRead(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
			l.Error(errors.Wrap(err, "TimerRead"))
		}
	}
}

//TimerUpdate
func TimerUpdate(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
			l.Error(errors.Wrap(err, "TimerUpdate"))
		}
	}
}

//TimerDelete
func TimerDelete(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
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
		if err = handleResponse(writer, err, nil); err != nil {
			l.Error(errors.Wrap(err, "TimerDelete"))
		}
	}
}

//TimerStart
func TimerStart(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
			l.Error(errors.Wrap(err, "TimerStart"))
		}
	}
}

//TimerPause
func TimerPause(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
			l.Error(errors.Wrap(err, "TimerPause"))
		}
	}
}

//TimerSubmit
func TimerSubmit(l common.Logger, f common.FunctionalTimer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.Contract
		var timer common.Timer
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
			l.Error(errors.Wrap(err, "TimerSubmit"))
		}
	}
}

//TimeSliceRead
func TimeSliceRead(l common.Logger, f common.FunctionalTimeSlice) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var timeSlice common.TimeSlice
		var contract common.Contract
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
			l.Error(errors.Wrap(err, "TimeSliceRead"))
		}
	}
}
