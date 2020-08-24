package bludgeonrestapi

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//doRequest
func doRequest(uri, method string, dataIn []byte) (response *http.Response, err error) {
	var request *http.Request

	//create client
	client := &http.Client{
		Timeout: ConfigTimeout,
	}
	//create request
	if request, err = http.NewRequest(method, uri, bytes.NewBuffer(dataIn)); err != nil {
		return
	}
	//execute request and parse response
	if response, err = client.Do(request); err != nil {
		return
	}
	//check to see if response was unsuccessful
	if response.StatusCode != http.StatusOK {
		var bytes []byte

		//attempt to read body and get bytes from response
		bytes, _ = ioutil.ReadAll(response.Body)
		response.Body.Close()
		//check if length of bytes greater than 0
		if len(bytes) > 0 {
			err = errors.New(string(bytes))
		} else {
			err = fmt.Errorf("Failure: %d", response.StatusCode)
		}
	}

	return
}
