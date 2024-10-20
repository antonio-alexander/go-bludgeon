package tests

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/client"
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/internal/cache"

	internal_mock "github.com/antonio-alexander/go-bludgeon/pkg/rest/client/mock"

	goqueue "github.com/antonio-alexander/go-queue"
	uuid "github.com/google/uuid"
	assert "github.com/stretchr/testify/assert"
)

type Fixture struct {
	client.Client
	client.Handler
	cache.Cache
	internal_mock.Mock
	goqueue.Peeker
	goqueue.Length
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

// REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randomString(nLetters ...int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	nLetter := 20
	if len(nLetters) > 0 {
		nLetter = nLetters[0]
	}
	b := make([]rune, nLetter)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewFixture(items ...interface{}) *Fixture {
	f := &Fixture{}
	for _, item := range items {
		switch item := item.(type) {
		case interface {
			client.Client
			client.Handler
		}:
			f.Client = item
			f.Handler = item
		case cache.Cache:
			f.Cache = item
		case client.Client:
			f.Client = item
		case client.Handler:
			f.Handler = item
		case internal_mock.Mock:
			f.Mock = item
		case interface {
			goqueue.Peeker
		}:
			f.Peeker = item
		}
	}
	return f
}

func (f *Fixture) TestChangeOperations(t *testing.T) {
	ctx := context.TODO()

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//upsert change
	changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	assert.NotEmpty(t, changeCreated.Id)
	assert.Equal(t, dataId, changeCreated.DataId)
	assert.Equal(t, dataVersion, changeCreated.DataVersion)
	assert.Equal(t, dataType, changeCreated.DataType)
	// assert.Equal(t, whenChanged, changeCreated.WhenChanged)
	assert.Equal(t, serviceName, changeCreated.DataServiceName)
	changeId := changeCreated.Id

	//read change
	changeRead, err := f.ChangeRead(ctx, changeId)
	assert.Nil(t, err)
	assert.NotNil(t, changeRead)
	assert.Equal(t, changeCreated, changeRead)

	//delete change
	err = f.ChangeDelete(ctx, changeId)
	assert.Nil(t, err)

	//read change
	changeRead, err = f.ChangeRead(ctx, changeId)
	assert.NotNil(t, err)
	assert.Nil(t, changeRead)
}

func (f *Fixture) TestChangeStreaming(t *testing.T) {
	var change *data.Change

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	changeReceived := make(chan struct{})

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//register handler
	handlerId, err := f.HandlerCreate(func(changes ...*data.Change) error {
		for _, c := range changes {
			if reflect.DeepEqual(change, c) {
				select {
				default:
					close(changeReceived)
				case <-changeReceived:
				}
			}
		}
		return nil
	})
	assert.Nil(t, err)

	//wait for handler to connect
	time.Sleep(10 * time.Second)

	//upsert change
	change, err = f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)

	//wait for change to be received
	select {
	case <-changeReceived:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm change received")
	}
	cancel()

	//unregister handler
	err = f.HandlerDelete(handlerId)
	assert.Nil(t, err)
}

func (f *Fixture) TestRegistrationOperations(t *testing.T) {
	ctx := context.TODO()

	// generate dynamic constants
	registrationId := randomString()
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	// create registration
	err := f.RegistrationUpsert(ctx, registrationId)
	assert.Nil(t, err)

	// upsert change
	changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)

	// read registration changes
	changes, err := f.RegistrationChangesRead(ctx, registrationId)
	assert.Nil(t, err)
	assert.NotNil(t, changes)
	assert.Contains(t, changes, changeCreated)

	// acknowledge change
	err = f.RegistrationChangeAcknowledge(ctx, registrationId, changeCreated.Id)
	assert.Nil(t, err)

	// read registration changes
	changes, err = f.RegistrationChangesRead(ctx, registrationId)
	assert.Nil(t, err)
	assert.Empty(t, changes)

	// delete registration
	err = f.RegistrationDelete(ctx, registrationId)
	assert.Nil(t, err)
}

// cache disabled, queue enabled
func (f *Fixture) TestChangeUpsertQueue(t *testing.T) {
	ctx := context.TODO()

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	// set the output mock values to generate an error
	f.MockDoRequest(nil, http.StatusInternalServerError,
		errors.New("mock error"))

	//upsert change
	changePartial := data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	}
	changeCreated, err := f.ChangeUpsert(ctx, changePartial)
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	assert.Empty(t, changeCreated.Id)

	// wait for two cycles for additional code coverage
	time.Sleep(2 * time.Second)

	// validate item in the queue
	item, underflow := f.PeekHead()
	assert.False(t, underflow)
	assert.NotNil(t, item)
	assert.Equal(t, item, changePartial)

	// set the output mock values to generate an error
	f.MockDoRequest(nil, http.StatusOK, nil)

	//wait for the go routine to start
	time.Sleep(time.Second)

	//validate that item is dequeued
	item, underflow = f.PeekHead()
	assert.True(t, underflow)
	assert.Nil(t, item)
}

func (f *Fixture) TestChangeCache(t *testing.T) {
	ctx := context.TODO()

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//upsert change
	changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	assert.NotEmpty(t, changeCreated.Id)
	changeId := changeCreated.Id

	//validate change cached
	changeCached := f.Read(changeId)
	assert.Equal(t, changeCreated, changeCached)

	//delete change
	err = f.ChangeDelete(ctx, changeId)
	assert.Nil(t, err)

	//validate change removed
	changeCached = f.Read(changeId)
	assert.Nil(t, changeCached)
}
