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
func serverTimerCreate(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var data interface{}
		var err error
		var token bludgeon.Token

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
func serverTimerRead(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var data interface{}
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

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
func serverTimerUpdate(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

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
func serverTimerDelete(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

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
		//check for errors
		if err == nil {
			writer.Write(nil)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

//serverTimerStart
func serverTimerStart(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerStart, server.CommandData{
						ID: contract.ID,
					}, token)
				}
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
func serverTimerPause(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerPause, server.CommandData{
						ID:        contract.ID,
						PauseTime: contract.PauseTime,
					}, token)
				}
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

//serverTimerSubmit
func serverTimerSubmit(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimerDelete, server.CommandData{
						ID:         contract.ID,
						FinishTime: contract.FinishTime,
					}, token)
				}
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
func serverTimeSliceCreate(s server.Functional) func(http.ResponseWriter, *http.Request) {
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
					if data, err = s.CommandHandler(bludgeon.CommandServerTimeSliceCreate, server.CommandData{
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
func serverTimeSliceRead(s server.Functional) func(http.ResponseWriter, *http.Request) {
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
func serverTimeSliceUpdate(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimeSliceUpdate, contract.TimeSlice, token)
				}
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
func serverTimeSliceDelete(s server.Functional) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error
		var contract rest.ContractServerIn
		var token bludgeon.Token

		//get token
		if token, err = getToken(request); err == nil {
			//read bytes from request
			if bytes, err = ioutil.ReadAll(request.Body); err == nil {
				if err = json.Unmarshal(bytes, &contract); err == nil {
					//attempt to execute the timer create
					_, err = s.CommandHandler(bludgeon.CommandServerTimeSliceDelete, server.CommandData{
						ID: contract.ID,
					}, token)
				}
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
