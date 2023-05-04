package tests

import (
	"context"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	employeesclient "github.com/antonio-alexander/go-bludgeon/employees/client"
	employeesdata "github.com/antonio-alexander/go-bludgeon/employees/data"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestFixture struct {
	client          client.Client
	cache           cache.Cache
	employeesClient employeesclient.Client
	changesClient   changesclient.Client
}

func NewTestFixture(items ...interface{}) *TestFixture {
	t := &TestFixture{}
	for _, item := range items {
		switch item := item.(type) {
		case client.Client:
			t.client = item
		case cache.Cache:
			t.cache = item
		case employeesclient.Client:
			t.employeesClient = item
		case changesclient.Client:
			t.changesClient = item
		}
	}
	return t
}

func (f *TestFixture) TestTimers(t *testing.T) {
	ctx := context.TODO()

	//create a timer
	timer, err := f.client.TimerCreate(ctx, data.TimerPartial{})
	assert.Nil(t, err)
	assert.NotEmpty(t, timer.ID)
	assert.Empty(t, timer.ActiveTimeSliceID)
	assert.False(t, timer.Completed)
	assert.False(t, timer.Archived)
	timerID := timer.ID

	//read the timer
	timerRead, err := f.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timer, timerRead)

	//start the timer
	timerStarted, err := f.client.TimerStart(ctx, timerID)
	assert.Nil(t, err)
	assert.NotEmpty(t, timerStarted.ActiveTimeSliceID)
	assert.NotZero(t, timerStarted.Start)

	//wait for a second
	time.Sleep(time.Second)

	//read the timer
	timerRead, err = f.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, timerRead.ElapsedTime, int64(time.Second))

	//stop the timer
	timerStopped, err := f.client.TimerStop(ctx, timerID)
	assert.Nil(t, err)

	//read the timer
	timerRead, err = f.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)

	//wait one second
	time.Sleep(time.Second)

	//read the timer
	timerRead, err = f.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerStopped, timerRead)

	//submit the timer
	tNow := time.Now()
	timerSubmitted, err := f.client.TimerSubmit(ctx, timerID, tNow.UnixNano())
	assert.Nil(t, err)

	//read the timer
	timerRead, err = f.client.TimerRead(ctx, timerID)
	assert.Nil(t, err)
	assert.Equal(t, timerSubmitted, timerRead)

	//delete the timer
	err = f.client.TimerDelete(ctx, timerID)
	assert.Nil(t, err)
}

func (f *TestFixture) generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func (f *TestFixture) TestTimerCache(t *testing.T) {
	ctx := context.TODO()

	//generate dynamic constants
	firstName, lastName := f.generateId(), f.generateId()
	emailAddress := f.generateId() + "@foobaf.duck"

	//create employee
	employeeCreated, err := f.employeesClient.EmployeeCreate(ctx, employeesdata.EmployeePartial{
		FirstName:    &firstName,
		LastName:     &lastName,
		EmailAddress: &emailAddress,
	})
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to create employee")
	}
	assert.NotNil(t, employeeCreated)
	employeeId := employeeCreated.ID

	//create timer
	timerCreated, err := f.client.TimerCreate(ctx, data.TimerPartial{
		EmployeeID: &employeeId,
	})
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to create timer")
	}
	assert.NotNil(t, timerCreated)
	timerId := timerCreated.ID
	t.Logf("created timer: %s", timerId)

	//validate timer in cache
	timerCached := new(data.Timer)
	err = f.cache.Read(timerId, timerCached)
	assert.Nil(t, err)
	assert.NotNil(t, timerCached)
	assert.Equal(t, timerCreated, timerCached)

	//delete employee
	err = f.employeesClient.EmployeeDelete(ctx, employeeId)
	assert.Nil(t, err)

	//wait for change to propagate
	time.Sleep(time.Second)

	//validate that the timer is removed from cache
	validateTimerDeletedFx := func() bool {
		timerCached = new(data.Timer)
		err = f.cache.Read(timerId, timerCached)
		return err != nil
	}
	tCheck := time.NewTicker(time.Second)
	defer tCheck.Stop()
	tStop := time.After(30 * time.Second)
	for stop := false; !stop; {
		select {
		case <-tCheck.C:
			if validateTimerDeletedFx() {
				stop = true
			}
		case <-tStop:
			assert.Fail(t, "unable to confirm timer deleted")
			stop = true
		}
	}
}
