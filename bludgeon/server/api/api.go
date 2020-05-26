package bludgeonrestapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/server/common"
)

func TimerCreate(address, port string) (timer bludgeon.Timer, err error) {
	var response *http.Response
	var bytes []byte

	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerCreate)
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
	var contract common.ContractServerIn

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerRead)
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
func TimerUpdate(address, port string, t bludgeon.Timer) (timer bludgeon.Timer, err error) {
	var bytes []byte
	var contract common.ContractServerIn
	var response *http.Response

	//store id in contract
	contract.Timer = t
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerUpdate)
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
func TimerDelete(address, port, id string) (err error) {
	var bytes []byte
	var contract common.ContractServerIn

	//store id in contract
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerDelete)
	//execute request and get response
	_, err = doRequest(uri, POST, bytes)

	return
}

//
func TimerStart(address, port string, timerID string, startTime time.Time) (timer bludgeon.Timer, err error) {
	var contract common.ContractServerIn
	var response *http.Response
	var bytes []byte

	//store id in contract
	contract.ID = timerID
	contract.StartTime = startTime.UnixNano()
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerStart)
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
func TimerPause(address, port, timerID string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	var contract common.ContractServerIn
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
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerPause)
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
func TimerSubmit(address, port, timerID string, finishTime time.Time) (timer bludgeon.Timer, err error) {
	var contract common.ContractServerIn
	var response *http.Response
	var bytes []byte

	//store id in contract
	contract.ID = timerID
	contract.FinishTime = finishTime.UnixNano()
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimerSubmit)
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

func TimeSliceRead(address, port, id string) (timeSlice bludgeon.TimeSlice, err error) {
	var response *http.Response
	var bytes []byte
	var contract common.ContractServerIn

	//set timeslice id
	contract.ID = id
	//marshal contract
	if bytes, err = json.Marshal(&contract); err != nil {
		return
	}
	//create uri
	uri := fmt.Sprintf(URIf, address, port, common.RouteTimeSliceRead)
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
