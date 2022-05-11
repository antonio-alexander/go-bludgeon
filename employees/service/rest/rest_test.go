package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/logic"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"
	"github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/employees/service/rest"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/stretchr/testify/assert"
)

var (
	address         string        = "localhost"
	port            string        = "8081"
	shutdownTimeout time.Duration = 15 * time.Second
	letterRunes     []rune        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type restServerTest struct {
	server interface {
		internal_server.Owner
		internal_server.Router
	}
	meta interface {
		meta.Owner
		meta.Serializer
		meta.Employee
	}
	logic interface {
		logic.Logic
	}
	client *http.Client
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	if _, ok := envs["BLUDGEON_REST_ADDRESS"]; ok {
		address = envs["BLUDGEON_REST_ADDRESS"]
	}
	if _, ok := envs["BLUDGEON_REST_PORT"]; ok {
		port = envs["BLUDGEON_REST_PORT"]
	}
	if _, ok := envs["BLUDGEON_REST_TIMEOUT"]; ok {
		if i, err := strconv.Atoi(envs["BLUDGEON_REST_SHUTDOWN_TIMEOUT"]); err != nil {
			shutdownTimeout = time.Duration(i) * time.Second
		}
	}
}

func new() *restServerTest {
	logger := internal_logger.New("bludgeon_rest_server_test")
	server := internal_server.New(logger)
	employeeMeta := memory.New()
	employeeLogic := logic.New(logger, employeeMeta)
	rest.New(logger, server, employeeLogic)
	return &restServerTest{
		server: server,
		meta:   employeeMeta,
		logic:  employeeLogic,
		client: &http.Client{},
	}
}

func (r *restServerTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s%s", address, port, route)
	request, err := http.NewRequest(method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, -1, err
	}
	response, err := r.client.Do(request)
	if err != nil {
		return nil, -1, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return bytes, response.StatusCode, err
}

func (r *restServerTest) initialize(t *testing.T) {
	err := r.server.Start(&internal_server.Configuration{
		Address:         address,
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
	})
	assert.Nil(t, err)
}

func (r *restServerTest) shutdown(t *testing.T) {
	r.meta.Shutdown()
	r.server.Stop()
}

func (r *restServerTest) TestEmployeeOperations(t *testing.T) {
	//create employee
	firstName, lastName := randomString(25), randomString(25)
	emailAddress := fmt.Sprintf("%s@foobar.duck", randomString(25))
	bytes, err := json.Marshal(&data.EmployeePartial{
		FirstName:    &firstName,
		LastName:     &lastName,
		EmailAddress: &emailAddress,
	})
	assert.Nil(t, err)
	bytes, statusCode, err := r.doRequest(data.RouteEmployees, http.MethodPost, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employeeCreated := &data.Employee{}
	err = json.Unmarshal(bytes, employeeCreated)
	assert.Nil(t, err)
	//read created employee
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeCreated.ID), http.MethodGet, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employeeRead := &data.Employee{}
	err = json.Unmarshal(bytes, employeeRead)
	assert.Nil(t, err)
	assert.Equal(t, employeeCreated, employeeRead)
	//read all employees
	bytes, statusCode, err = r.doRequest(data.RouteEmployeesSearch, http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employees := []data.Employee{}
	err = json.Unmarshal(bytes, &employees)
	assert.Nil(t, err)
	assert.Len(t, employees, 1)
	assert.Contains(t, employees, *employeeCreated)
	//update employee
	updatedFirstName := randomString(25)
	bytes, err = json.Marshal(&data.EmployeePartial{
		FirstName: &updatedFirstName,
	})
	assert.Nil(t, err)
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeCreated.ID), http.MethodPut, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employeeUpdated := &data.Employee{}
	err = json.Unmarshal(bytes, employeeUpdated)
	assert.Nil(t, err)
	assert.Equal(t, updatedFirstName, employeeUpdated.FirstName)
	//read updated employee
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeUpdated.ID), http.MethodGet, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employeeRead = &data.Employee{}
	err = json.Unmarshal(bytes, employeeRead)
	assert.Nil(t, err)
	assert.Equal(t, employeeUpdated, employeeRead)
	//delete employee
	_, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeCreated.ID), http.MethodDelete, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, statusCode)
	//delete employee again
	_, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeCreated.ID), http.MethodDelete, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, statusCode)
	//attempt to read deleted employee
	_, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeUpdated.ID), http.MethodGet, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestEmployeesRestService(t *testing.T) {
	r := new()
	r.initialize(t)
	t.Run("Test Employee Operations", r.TestEmployeeOperations)
	r.shutdown(t)
}
