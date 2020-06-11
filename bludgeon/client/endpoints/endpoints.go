package bludgeonclientendpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/client/common"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/functional"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest_server"

	errors "github.com/pkg/errors"
)

//BuildRoutes will create all the routes and their functions to execute when received
func BuildRoutes(c interface {
	client.Functional
	bludgeon.Logger
}) []rest.HandleFuncConfig {
	//TODO: replace nil with actual logger
	return []rest.HandleFuncConfig{
		//admin
		//timer
		{Route: common.RouteTimerCreate, Method: POST, HandleFx: ClientTimerCreate(c)},
		{Route: common.RouteTimerRead, Method: POST, HandleFx: ClientTimerRead(c)},
		{Route: common.RouteTimerUpdate, Method: POST, HandleFx: ClientTimerUpdate(c)},
		{Route: common.RouteTimerDelete, Method: POST, HandleFx: ClientTimerDelete(c)},
		{Route: common.RouteTimerStart, Method: POST, HandleFx: ClientTimerStart(c)},
		{Route: common.RouteTimerPause, Method: POST, HandleFx: ClientTimerPause(c)},
		{Route: common.RouteTimerSubmit, Method: POST, HandleFx: ClientTimerSubmit(c)},
	}
}

//ClientTimerCreate
func ClientTimerCreate(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var data interface{}
		var bytes []byte
		var err error

		//attempt to execute the timer create
		if data, err = c.CommandHandler(bludgeon.CommandClientTimerCreate, nil); err == nil {
			//attempt to marshal into bytes
			if timer, ok := data.(bludgeon.Timer); ok {
				bytes, err = json.Marshal(timer)
			} else {
				err = errors.New("unable to cast into timer")
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerCreate")
			c.Error(err)
		}
	}
}

//ClientTimerRead
func ClientTimerRead(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var data interface{}
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = c.CommandHandler(bludgeon.CommandClientTimerRead, common.CommandData{
					ID: contract.ID,
				}); err == nil {
					//attempt to marshal into bytes
					if timer, ok := data.(bludgeon.Timer); ok {
						bytes, err = json.Marshal(timer)
					} else {
						err = errors.New("unable to cast into timer")
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerRead")
			c.Error(err)
		}
	}
}

//ClientTimerUpdate
func ClientTimerUpdate(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var bytes []byte
		var data interface{}
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = c.CommandHandler(bludgeon.CommandClientTimerUpdate, contract.Timer); err == nil {
					if timer, ok := data.(bludgeon.Timer); ok {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerUpdate")
			c.Error(err)
		}
	}
}

//ClientTimerDelete
func ClientTimerDelete(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				_, err = c.CommandHandler(bludgeon.CommandClientTimerDelete, common.CommandData{
					ID: contract.ID,
				})
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerDelete")
			c.Error(err)
		}
	}
}

//ClientTimerStart
func ClientTimerStart(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var data interface{}
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = c.CommandHandler(bludgeon.CommandClientTimerStart, common.CommandData{
					ID:        contract.ID,
					StartTime: time.Unix(0, contract.StartTime),
				}); err == nil {
					if timer, ok := data.(bludgeon.Timer); !ok {
						err = errors.New("unable to cast into timer")
					} else {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerStart")
			c.Error(err)
		}
	}
}

//ClientTimerDelete
func ClientTimerPause(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var data interface{}
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = c.CommandHandler(bludgeon.CommandClientTimerPause, common.CommandData{
					ID:        contract.ID,
					PauseTime: time.Unix(0, contract.PauseTime),
				}); err == nil {
					if timer, ok := data.(bludgeon.Timer); !ok {
						err = errors.New("unable to cast into timer")
					} else {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerPause")
			c.Error(err)
		}
	}
}

//ClientTimerSubmit
func ClientTimerSubmit(c interface {
	client.Functional
	bludgeon.Logger
}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract common.ContractClientIn
		var data interface{}
		var bytes []byte
		var err error

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = c.CommandHandler(bludgeon.CommandClientTimerSubmit, common.CommandData{
					ID:         contract.ID,
					FinishTime: time.Unix(0, contract.FinishTime),
				}); err == nil {
					if timer, ok := data.(bludgeon.Timer); !ok {
						err = errors.New("unable to cast into timer")
					} else {
						bytes, err = json.Marshal(&timer)
					}
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ClientTimerSubmit")
			c.Error(err)
		}
	}
}
