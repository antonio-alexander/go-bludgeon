package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logic "github.com/antonio-alexander/go-bludgeon/logic"
)

type rest struct {
	sync.RWMutex
	address string
	port    string
}

func New(address, port string) interface {
	logic.Logic
} {
	return &rest{
		address: address,
		port:    port,
	}
}

func (r *rest) TimerCreate() (timer data.Timer, err error) {
	var response *http.Response
	var bytes []byte

	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerCreate)
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

//TimerRead
func (r *rest) TimerRead(id string) (timer data.Timer, err error) {
	var response *http.Response
	var bytes []byte
	var contract data.Contract

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerRead)
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
func (r *rest) TimerUpdate(t data.Timer) (timer data.Timer, err error) {
	var bytes []byte
	var contract data.Contract
	var response *http.Response

	//store id in contract
	contract.Timer = t
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerUpdate)
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
func (r *rest) TimerDelete(id string) (err error) {
	var bytes []byte
	var contract data.Contract

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerDelete)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func (r *rest) TimerStart(timerID string, startTime time.Time) (timer data.Timer, err error) {
	var contract data.Contract
	var response *http.Response
	var bytes []byte

	//store id in contract
	contract.ID = timerID
	contract.StartTime = startTime.UnixNano()
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerStart)
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
func (r *rest) TimerPause(timerID string, pauseTime time.Time) (timer data.Timer, err error) {
	var contract data.Contract
	var response *http.Response
	var bytes []byte

	//store id in contract
	contract.ID = timerID
	contract.PauseTime = pauseTime.UnixNano()
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerPause)
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
func (r *rest) TimerSubmit(timerID string, finishTime time.Time) (timer data.Timer, err error) {
	var contract data.Contract
	var response *http.Response
	var bytes []byte

	//store id in contract
	contract.ID = timerID
	contract.FinishTime = finishTime.UnixNano()
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimerSubmit)
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

func (r *rest) TimeSliceRead(id string) (timeSlice data.TimeSlice, err error) {
	var response *http.Response
	var bytes []byte
	var contract data.Contract

	//set timeslice id
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, r.address, r.port, data.RouteTimeSliceRead)
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
