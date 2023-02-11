package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/logic"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"
	"github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	"github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	service "github.com/antonio-alexander/go-bludgeon/employees/service/rest"

	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_logic.json"

var (
	letterRunes         = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	configMetaMysql     = new(internal_mysql.Configuration)
	configMetaFile      = new(internal_file.Configuration)
	configLogger        = new(internal_logger.Configuration)
	configServer        = new(internal_server.Configuration)
	configChangesClient = new(changesclientrest.Configuration)
	configLogic         = new(logic.Configuration)
)

type restServiceTest struct {
	server interface {
		internal.Initializer
		internal.Configurer
	}
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
	}
	client *http.Client
}

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configLogger.Default()
	configLogger.FromEnv(envs)
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configLogic.Default()
	configLogic.FromEnv(envs)
	configChangesClient.Default()
	configChangesClient.FromEnv(envs)
	configServer.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "8081"
	configServer.ShutdownTimeout = 10 * time.Second
}

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func newRestServiceTest(metaType string) *restServiceTest {
	var employeeMeta interface {
		meta.Employee
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}

	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "bludgeon_rest_server_test",
		Level:  internal_logger.Trace,
	})
	switch metaType {
	default:
		employeeMeta = memory.New()
	case "mysql":
		employeeMeta = mysql.New()
	case "file":
		employeeMeta = file.New()
	}
	employeeMeta.SetUtilities(logger)
	changesClient := changesclientrest.New()
	changesClient.SetUtilities(logger)
	employeeLogic := logic.New()
	employeeLogic.SetParameters(employeeMeta, changesClient)
	employeeLogic.SetUtilities(logger)
	employeeLogic.Configure(configLogic)
	employeeService := service.New()
	employeeService.SetUtilities(logger)
	employeeService.SetParameters(employeeLogic)
	server := internal_server.New()
	server.SetUtilities(logger)
	server.SetParameters(employeeService)
	return &restServiceTest{
		server:        server,
		meta:          employeeMeta,
		changesClient: changesClient,
		client:        &http.Client{},
	}
}

func (r *restServiceTest) doRequest(route, method string, data []byte) ([]byte, int, error) {
	uri := fmt.Sprintf("http://%s:%s%s", configServer.Address, configServer.Port, route)
	request, err := http.NewRequest(method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, -1, err
	}
	response, err := r.client.Do(request)
	if err != nil {
		return nil, -1, err
	}
	bytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	return bytes, response.StatusCode, err
}

func (r *restServiceTest) initialize(t *testing.T, metaType string) {
	switch metaType {
	case "mysql":
		err := r.meta.Configure(configMetaMysql)
		assert.Nil(t, err)
	case "file":
		err := r.meta.Configure(configMetaFile)
		assert.Nil(t, err)
	}
	err := r.meta.Initialize()
	assert.Nil(t, err)
	err = r.changesClient.Configure(configChangesClient)
	assert.Nil(t, err)
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Initialize()
	assert.Nil(t, err)
}

func (r *restServiceTest) shutdown(t *testing.T) {
	r.server.Shutdown()
	r.changesClient.Shutdown()
	r.meta.Shutdown()
}

func (r *restServiceTest) TestEmployeeOperations(t *testing.T) {
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
	employeeId := employeeCreated.ID
	defer func() {
		r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeId), http.MethodDelete, nil)
	}()

	//read created employee
	bytes, statusCode, err = r.doRequest(fmt.Sprintf(data.RouteEmployeesIDf, employeeCreated.ID), http.MethodGet, bytes)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	employeeRead := &data.Employee{}
	err = json.Unmarshal(bytes, employeeRead)
	assert.Nil(t, err)
	assert.Equal(t, employeeCreated, employeeRead)
	//read all employees
	search := &data.EmployeeSearch{IDs: []string{employeeId}}
	bytes, statusCode, err = r.doRequest(data.RouteEmployeesSearch+search.ToParams(), http.MethodGet, nil)
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

func testEmployeesRestService(t *testing.T, metaType string) {
	r := newRestServiceTest(metaType)

	r.initialize(t, metaType)
	defer r.shutdown(t)

	t.Run("Test Employee Operations", r.TestEmployeeOperations)
}

func TestEmployeesRestServiceMemory(t *testing.T) {
	testEmployeesRestService(t, "memory")
}

func TestEmployeesRestServiceFile(t *testing.T) {
	testEmployeesRestService(t, "file")
}

func TestEmployeesRestServiceMysql(t *testing.T) {
	testEmployeesRestService(t, "mysql")
}
