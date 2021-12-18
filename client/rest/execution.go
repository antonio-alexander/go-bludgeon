package restapi

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//doRequest
func doRequest(uri, method string, data []byte) (response *http.Response, err error) {
	var request *http.Request

	//create client
	//create request
	//execute request and parse response
	//check to see if response was unsuccessful
	//attempt to read body and get bytes from response
	//check if length of bytes greater than 0
	client := &http.Client{
		Timeout: ConfigTimeout,
	}
	if request, err = http.NewRequest(method, uri, bytes.NewBuffer(data)); err != nil {
		return
	}
	if response, err = client.Do(request); err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		data, err = ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return
		}
		if len(data) > 0 {
			err = errors.New(string(data))
		} else {
			err = fmt.Errorf("failure: %d", response.StatusCode)
		}
	}

	return
}
