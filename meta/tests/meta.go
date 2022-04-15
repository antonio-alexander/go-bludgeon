package tests

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/data"
	"github.com/antonio-alexander/go-bludgeon/meta"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestEmployeeCRUD(m meta.Employee) func(*testing.T) {
	return func(t *testing.T) {
		//create
		firstName := randomString(15)
		lastName := randomString(15)
		emailAddress := randomString(20) + "@foobar.duck"
		employee, err := m.EmployeeCreate(data.EmployeePartial{
			FirstName:    &firstName,
			LastName:     &lastName,
			EmailAddress: &emailAddress,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, employee.ID)
		assert.Equal(t, firstName, employee.FirstName)
		assert.Equal(t, lastName, employee.LastName)
		assert.Equal(t, emailAddress, employee.EmailAddress)
		//read
		employeeRead, err := m.EmployeeRead(employee.ID)
		assert.Nil(t, err)
		assert.Equal(t, employee, employeeRead)
		//read multiple
		employeesRead, err := m.EmployeesRead(data.EmployeeSearch{
			FirstName:    &firstName,
			LastName:     &lastName,
			EmailAddress: &emailAddress,
		})
		assert.Nil(t, err)
		assert.Len(t, employeesRead, 1)
		assert.Condition(t, func() bool {
			for _, employeeRead := range employeesRead {
				if reflect.DeepEqual(employeeRead, employee) {
					return true
				}
			}
			return false
		})
		assert.Contains(t, employeesRead, employee)
		*employee = *employeeRead
		//update
		//KIM: if we don't sleep, the tests below will fail for
		// last_updated
		firstName = randomString(25)
		time.Sleep(time.Second)
		employeeUpdated, err := m.EmployeeUpdate(employee.ID, data.EmployeePartial{
			FirstName: &firstName,
		})
		assert.Nil(t, err)
		assert.Equal(t, firstName, employeeUpdated.FirstName)
		assert.Greater(t, employeeUpdated.Version, employee.Version)
		assert.Greater(t, employeeUpdated.LastUpdated, employee.LastUpdated)
		//delete
		err = m.EmployeeDelete(employee.ID)
		assert.Nil(t, err)
		err = m.EmployeeDelete(employee.ID)
		assert.NotNil(t, err)
		_, err = m.EmployeeRead(employee.ID)
		assert.NotNil(t, err)
		//read
		employeesRead, err = m.EmployeesRead(data.EmployeeSearch{
			FirstName:    &firstName,
			LastName:     &lastName,
			EmailAddress: &emailAddress,
		})
		assert.Nil(t, err)
		assert.Len(t, employeesRead, 0)
	}
}

func TestEmployeesRead(m interface {
	meta.Employee
}) func(*testing.T) {
	return func(t *testing.T) {
		//TODO: create test
	}
}

func TestTimerCRUD(m interface {
	meta.Employee
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create employee
		firstName := randomString(15)
		lasttName := randomString(15)
		emailAddress := randomString(15) + "@foobar.duck"
		employee, err := m.EmployeeCreate(data.EmployeePartial{
			FirstName:    &firstName,
			LastName:     &lasttName,
			EmailAddress: &emailAddress,
		})
		assert.Nil(t, err)
		//create timer
		comment := randomString(25)
		timer, err := m.TimerCreate(data.TimerPartial{
			Comment:    &comment,
			EmployeeID: &employee.ID,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, timer.ID)
		assert.Equal(t, timer.Comment, comment)
		assert.Equal(t, timer.EmployeeID, employee.ID)
		//read
		timerRead, err := m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timer, timerRead)
		timers, err := m.TimersRead(data.TimerSearch{
			EmployeeID: &employee.ID,
		})
		assert.Nil(t, err)
		assert.Contains(t, timers, timer)
		//update
		updatedComment := randomString(25)
		timerUpdated, err := m.TimerUpdate(timer.ID, data.TimerPartial{
			Comment: &updatedComment,
		})
		assert.Nil(t, err)
		assert.Equal(t, updatedComment, timerUpdated.Comment)
		timer.Comment = timerUpdated.Comment
		timer.Audit = timerUpdated.Audit
		assert.Equal(t, timer, timerUpdated)
		//delete
		err = m.TimerDelete(timer.ID)
		assert.Nil(t, err)
		err = m.TimerDelete(timer.ID)
		assert.NotNil(t, err)
		//read
		timerRead, err = m.TimerRead(timer.ID)
		assert.NotNil(t, err)
		assert.Nil(t, timerRead)
	}
}

func TestTimersRead(m interface {
	meta.Employee
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//TODO: create test
	}
}

func TestTimerLogic(m interface {
	meta.Employee
	meta.Timer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create employee
		firstName := randomString(15)
		lasttName := randomString(15)
		emailAddress := randomString(15) + "@foobar.duck"
		employee, err := m.EmployeeCreate(data.EmployeePartial{
			FirstName:    &firstName,
			LastName:     &lasttName,
			EmailAddress: &emailAddress,
		})
		assert.Nil(t, err)
		//create timer
		comment := randomString(25)
		timer, err := m.TimerCreate(data.TimerPartial{
			Comment:    &comment,
			EmployeeID: &employee.ID,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, timer.ID)
		assert.Equal(t, timer.Comment, comment)
		assert.Equal(t, timer.EmployeeID, employee.ID)
		//start
		timerStarted, err := m.TimerStart(timer.ID)
		assert.Nil(t, err)
		assert.Greater(t, timerStarted.Start, int64(0))
		assert.Zero(t, timerStarted.Finish)
		time.Sleep(time.Second)
		//read
		timerRead, err := m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, timerRead.ElapsedTime, int64(1))
		//stop
		timerStopped, err := m.TimerStop(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timerStarted.Start, timerStopped.Start)
		//read
		timerRead, err = m.TimerRead(timer.ID)
		assert.Nil(t, err)
		assert.Equal(t, timerStopped, timerRead)
	}
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
