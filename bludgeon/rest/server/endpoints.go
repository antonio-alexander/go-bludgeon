package bludgeonrestserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server"
	errors "github.com/pkg/errors"
)

//serverTimerCreate
func serverTimerCreate(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
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
			log.Println(err)
		}
	}
}

//serverTimerRead
func serverTimerRead(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
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
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerRead, server.CommandData{
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
			log.Println(err)
		}
	}
}

//serverTimerUpdate
func serverTimerUpdate(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
		var token bludgeon.Token
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerUpdate, contract.Timer, token)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerUpdate")
			log.Println(err)
		}
	}
}

//serverTimerDelete
func serverTimerDelete(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
		var token bludgeon.Token
		var bytes []byte
		var err error

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerDelete, server.CommandData{
						ID: contract.ID,
					}, token)
				}
			}
		}
		//handle errors
		if err = handleResponse(writer, err, bytes); err != nil {
			err = errors.Wrap(err, "ServerTimerDelete")
			log.Println(err)
		}
	}
}

//serverTimerStart
func serverTimerStart(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
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
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerStart, server.CommandData{
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
			log.Println(err)
		}
	}
}

//serverTimerDelete
func serverTimerPause(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
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
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerPause, server.CommandData{
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
			log.Println(err)
		}
	}
}

//serverTimerSubmit
func serverTimerSubmit(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var contract rest.ContractServerIn
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
					if data, err = s.CommandHandler(bludgeon.CommandServerTimerSubmit, server.CommandData{
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
			log.Println(err)
		}
	}
}

//serverTimeSliceRead
func serverTimeSliceRead(s server.Functional, log *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					if data, err = s.CommandHandler(bludgeon.CommandServerTimeSliceRead, server.CommandData{
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
			log.Println(err)
		}
	}
}
