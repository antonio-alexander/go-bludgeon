package bludgeonrestserver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	// client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server"
)

//serverTimerCreate
func serverTimerCreate(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error

		//attempt to execute the timer create
		if data, err = server.CommandHandler(bludgeon.CommandServerTimerCreate, nil); err == nil {
			//attempt to marshal into bytes
			if timer, ok := data.(bludgeon.Timer); ok {
				bytes, err = json.Marshal(timer)
			} else {
				err = errors.New("unable to cast into timer")
			}
		}
		//check for errors
		if err == nil {
			writer.Write(bytes)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimerRead
func serverTimerRead(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var data interface{}
		var bytes []byte
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = server.CommandHandler(bludgeon.CommandServerTimerRead, contract.Timer.UUID); err == nil {
					//attempt to marshal into bytes
					if timer, ok := data.(bludgeon.Timer); ok {
						bytes, err = json.Marshal(timer)
					} else {
						err = errors.New("unable to cast into timer")
					}
				}
			}
		}
		//check for errors
		if err == nil {
			writer.Write(bytes)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimerUpdate
func serverTimerUpdate(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				_, err = server.CommandHandler(bludgeon.CommandServerTimerUpdate, contract.Timer)
			}
		}
		//check for errors
		if err == nil {
			writer.Write(nil)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimerDelete
func serverTimerDelete(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				_, err = server.CommandHandler(bludgeon.CommandServerTimerDelete, contract.Timer.UUID)
			}
		}
		//check for errors
		if err == nil {
			writer.Write(nil)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimeSliceCreate
func serverTimeSliceCreate(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = server.CommandHandler(bludgeon.CommandServerTimeSliceCreate, contract.TimeSlice.TimerUUID); err == nil {
					//attempt to marshal into bytes
					if timeSlice, ok := data.(bludgeon.TimeSlice); ok {
						bytes, err = json.Marshal(timeSlice)
					} else {
						err = errors.New("unable to cast into time slice")
					}
				}
			}
		}
		//check for errors
		if err == nil {
			writer.Write(bytes)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimeSliceRead
func serverTimeSliceRead(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				if data, err = server.CommandHandler(bludgeon.CommandServerTimeSliceRead, contract.TimeSlice.UUID); err == nil {
					//attempt to marshal into bytes
					if timeSlice, ok := data.(bludgeon.TimeSlice); ok {
						bytes, err = json.Marshal(timeSlice)
					} else {
						err = errors.New("unable to cast into time slice")
					}
				}
			}
		}
		//check for errors
		if err == nil {
			writer.Write(bytes)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimeSliceUpdate
func serverTimeSliceUpdate(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				_, err = server.CommandHandler(bludgeon.CommandServerTimeSliceUpdate, contract.TimeSlice)
			}
		}
		//check for errors
		if err == nil {
			writer.Write(nil)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimeSliceDelete
func serverTimeSliceDelete(server server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn

		//read bytes from request
		if bytes, err = ioutil.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &contract); err == nil {
				//attempt to execute the timer create
				_, err = server.CommandHandler(bludgeon.CommandServerTimeSliceDelete, contract.TimeSlice.UUID)
			}
		}
		//check for errors
		if err == nil {
			writer.Write(nil)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}
