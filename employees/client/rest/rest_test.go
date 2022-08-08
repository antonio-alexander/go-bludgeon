package rest_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	restclient "github.com/antonio-alexander/go-bludgeon/employees/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/stretchr/testify/assert"
)

var (
	config      *restclient.Configuration
	letterRunes []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config = new(restclient.Configuration)
	config.Default()
	config.FromEnv(pwd, envs)
	rand.Seed(time.Now().UnixNano())
}

type restClientTest struct {
	client restclient.Client
}

func newRestclientTest() *restClientTest {
	logger := logger.New("bludgeon_rest_server_test")
	client := restclient.New(logger)
	return &restClientTest{
		client: client,
	}
}

func (r *restClientTest) Initialize(t *testing.T) {
	err := r.client.Initialize(config)
	assert.Nil(t, err)
}

func (r *restClientTest) TestEmployeeOperations(t *testing.T) {
	ctx := context.TODO()
	firstName, lastName := randomString(25), randomString(25)
	emailAddress := fmt.Sprintf("%s@foobar.duck", randomString(25))

	//create employee
	employeeCreated, err := r.client.EmployeeCreate(ctx, data.EmployeePartial{
		FirstName:    &firstName,
		LastName:     &lastName,
		EmailAddress: &emailAddress,
	})
	assert.Nil(t, err)
	assert.Equal(t, firstName, employeeCreated.FirstName)
	assert.Equal(t, lastName, employeeCreated.LastName)
	assert.Equal(t, emailAddress, employeeCreated.EmailAddress)
	employeeID := employeeCreated.ID

	//read created employee
	employeeRead, err := r.client.EmployeeRead(ctx, employeeID)
	assert.Nil(t, err)
	assert.Equal(t, employeeCreated, employeeRead)

	//read all employees
	employees, err := r.client.EmployeesRead(ctx, data.EmployeeSearch{
		IDs: []string{employeeID},
	})
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Len(t, employees, 1)
	assert.Condition(t, func() bool {
		for _, employee := range employees {
			if reflect.DeepEqual(employee, employeeRead) {
				return true
			}
		}
		return false
	})

	//update employee
	updatedFirstName := randomString(25)
	employeeUpdated, err := r.client.EmployeeUpdate(ctx, employeeID, data.EmployeePartial{
		FirstName: &updatedFirstName,
	})
	assert.Nil(t, err)
	assert.Equal(t, updatedFirstName, employeeUpdated.FirstName)

	//read updated employee
	employeeRead, err = r.client.EmployeeRead(ctx, employeeID)
	assert.Nil(t, err)
	assert.Equal(t, employeeUpdated, employeeRead)

	//delete employee
	err = r.client.EmployeeDelete(ctx, employeeID)
	assert.Nil(t, err)

	//delete employee again
	err = r.client.EmployeeDelete(ctx, employeeID)
	assert.NotNil(t, err)

	//attempt to read deleted employee
	employeeRead, err = r.client.EmployeeRead(ctx, employeeID)
	assert.Nil(t, employeeRead)
	assert.NotNil(t, err)
}

func TestEmployeesRestClient(t *testing.T) {
	r := newRestclientTest()
	r.Initialize(t)
	t.Run("Test Employee Operations", r.TestEmployeeOperations)
}
