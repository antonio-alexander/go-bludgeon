package bludgeonserverendpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/server/common"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/functional"

	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest_server"
	errors "github.com/pkg/errors"
)

//BuildRoutes will create all the routes and their functions to execute when received
func BuildRoutes(s interface {
	bludgeon.Logger
	server.Functional
}) []rest.HandleFuncConfig {
	return []rest.HandleFuncConfig{
		//admin
		{Route: common.RouteServerStop, Method: POST, HandleFx: ServerStop(s)},
		//timer
		{Route: common.RouteTimerCreate, Method: POST, HandleFx: ServerTimerCreate(s)},
		{Route: common.RouteTimerRead, Method: POST, HandleFx: ServerTimerRead(s)},
		{Route: common.RouteTimerUpdate, Method: POST, HandleFx: ServerTimerUpdate(s)},
		{Route: common.RouteTimerDelete, Method: POST, HandleFx: ServerTimerDelete(s)},
		{Route: common.RouteTimerStart, Method: POST, HandleFx: ServerTimerStart(s)},
		{Route: common.RouteTimerPause, Method: POST, HandleFx: ServerTimerPause(s)},
		{Route: common.RouteTimerSubmit, Method: POST, HandleFx: ServerTimerSubmit(s)},
		//time slice
		{Route: common.RouteTimeSliceRead, Method: POST, HandleFx: ServerTimeSliceRead(s)},
	}
}

//ServerStop
func ServerStop(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var token bludgeon.Token
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//attempt to execute the timer create
			_, err = s.CommandHandler(bludgeon.CommandServerStop, nil, token)
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerCreate")
			s.Error(err)
		}
	}
}

//ServerTimerCreate
func ServerTimerCreate(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var token bludgeon.Token
		var data interface{}
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//attempt to execute the timer create
			if data, err = s.CommandHandler(bludgeon.CommandServerTimerCreate, nil, token); err == nil {
				//attempt to marshal into bytes
				if timer, ok := data.(bludgeon.Timer); ok {
					bytes, err = json.Marshal(timer)
				} else {
					err = errors.New("unable to cast into timer")
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerCreate")
			s.Error(err)
		}
	}
}

//ServerTimerRead
func ServerTimerRead(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var data interface{}
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerRead, common.CommandData{
						ID: contract.ID,
					}, token); err == nil {
						//attempt to marshal into bytes
						if timer, ok := data.(bludgeon.Timer); ok {
							bytes, err = json.Marshal(timer)
						} else {
							err = errors.New("unable to cast into timer")
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerRead")
			s.Error(err)
		}
	}
}

//ServerTimerUpdate
func ServerTimerUpdate(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var bytes []byte
		var data interface{}
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerUpdate, contract.Timer, token); err == nil {
						if timer, ok := data.(bludgeon.Timer); ok {
							bytes, err = json.Marshal(&timer)
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerUpdate")
			s.Error(err)
		}
	}
}

//ServerTimerDelete
func ServerTimerDelete(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerDelete, common.CommandData{
						ID: contract.ID,
					}, token)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerDelete")
			s.Error(err)
		}
	}
}

//ServerTimerStart
func ServerTimerStart(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var data interface{}
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerStart, common.CommandData{
						ID:        contract.ID,
						StartTime: time.Unix(0, contract.StartTime),
					}, token); err == nil {
						if timer, ok := data.(bludgeon.Timer); !ok {
							err = errors.New("unable to cast into timer")
						} else {
							bytes, err = json.Marshal(&timer)
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerStart")
			s.Error(err)
		}
	}
}

//ServerTimerDelete
func ServerTimerPause(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var data interface{}
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerPause, common.CommandData{
						ID:        contract.ID,
						PauseTime: time.Unix(0, contract.PauseTime),
					}, token); err == nil {
						if timer, ok := data.(bludgeon.Timer); !ok {
							err = errors.New("unable to cast into timer")
						} else {
							bytes, err = json.Marshal(&timer)
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerPause")
			s.Error(err)
		}
	}
}

//ServerTimerSubmit
func ServerTimerSubmit(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractServerIn
		var token bludgeon.Token
		var data interface{}
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerSubmit, common.CommandData{
						ID:         contract.ID,
						FinishTime: time.Unix(0, contract.FinishTime),
					}, token); err == nil {
						if timer, ok := data.(bludgeon.Timer); !ok {
							err = errors.New("unable to cast into timer")
						} else {
							bytes, err = json.Marshal(&timer)
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerSubmit")
			s.Error(err)
		}
	}
}

//ServerTimeSliceRead
func ServerTimeSliceRead(s interface {
	bludgeon.Logger
	server.Functional
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error
		var contract common.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimeSliceRead, common.CommandData{
						ID: contract.ID,
					}, token); err == nil {
						//attempt to marshal into bytes
						if timeSlice, ok := data.(bludgeon.TimeSlice); ok {
							bytes, err = json.Marshal(timeSlice)
						} else {
							err = errors.New("unable to cast into time slice")
						}
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimeSliceRead")
			s.Error(err)
		}
	}
}
