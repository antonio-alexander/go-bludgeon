package bludgeonrestapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest"
)

func TimerCreate(address, port string) (timer bludgeon.Timer, err error) {
	var response *http.Response
	var bytes []byte

	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerCreate)
	//execute request and get response
	if response, err = doRequest(uri, POST, nil); err != nil {
		return
	}
	//get bytes from response
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	//close the response body
	response.Body.Close()
	//attempt to unmarshal json
	err = json.Unmarshal(bytes, &timer)

	return
}

//
func TimerRead(address, port, id string) (timer bludgeon.Timer, err error) {
	var response *http.Response
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerRead)
	//execute request and get response
	if response, err = doRequest(uri, POST, bytes); err != nil {
		return
	}
	//get bytes from response
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	//close the response body
	response.Body.Close()
	//attempt to unmarshal json
	err = json.Unmarshal(bytes, &timer)

	return
}

//
func TimerUpdate(address, port string, timer bludgeon.Timer) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.Timer = timer
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//marshal timer
	if bytes, err = json.Marshal(&timer); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerUpdate)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimerDelete(address, port, id string) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerDelete)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimerStart(address, port string, timerID string, startTime time.Time) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = timerID
	contract.StartTime = startTime
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerStart)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimerPause(address, port, timerID string, pauseTime time.Time) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = timerID
	contract.PauseTime = pauseTime
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerPause)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimerSubmit(address, port, timerID string, finishTime time.Time) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = timerID
	contract.FinishTime = finishTime
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerSubmit)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimeSliceCreate(address, port, id string) (timeSlice bludgeon.TimeSlice, err error) {
	var response *http.Response
	var bytes []byte
	var contract rest.ContractServerIn

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimeSliceCreate)
	//execute request and get response
	if response, err = doRequest(uri, POST, bytes); err != nil {
		return
	}
	//get bytes from response
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	//close the response body
	response.Body.Close()
	//attempt to unmarshal json
	err = json.Unmarshal(bytes, &timeSlice)

	return
}

//
func TimeSliceRead(address, port, id string) (timeSlice bludgeon.TimeSlice, err error) {
	var response *http.Response
	var bytes []byte
	var contract rest.ContractServerIn

	//set timeslice id
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimeSliceRead)
	//execute request and get response
	if response, err = doRequest(uri, POST, bytes); err != nil {
		return
	}
	//get bytes from response
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	//close the response body
	response.Body.Close()
	//attempt to unmarshal json
	err = json.Unmarshal(bytes, &timeSlice)

	return
}

//
func TimeSliceUpdate(address, port string, timeSlice bludgeon.TimeSlice) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//set timeslice id
	contract.TimeSlice = timeSlice
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerUpdate)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimeSliceDelete(address, port, id string) (err error) {
	var bytes []byte
	var contract rest.ContractServerIn

	//set timeslice id
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, rest.RouteTimerDelete)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}
