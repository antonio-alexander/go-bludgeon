package restclient_test

import (
	"context"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	restclient "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var config = new(restclient.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config.Default()
	config.FromEnv(envs)
	rand.Seed(time.Now().UnixNano())
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

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

type restClientTest struct {
	client interface {
		client.Client
		client.Handler
		internal.Configurer
		internal.Initializer
	}
}

func newRestclientTest() *restClientTest {
	logger := logger.New()
	// logger.Configure(&logger.Configuration{
	// 	Level:  logger.Trace,
	// 	Prefix: "bludgeon_rest_server_test",
	// })
	client := restclient.New()
	client.SetUtilities(logger)
	return &restClientTest{
		client: client,
	}
}

func (r *restClientTest) Initialize(t *testing.T) {
	err := r.client.Configure(config)
	assert.Nil(t, err)
	err = r.client.Initialize()
	assert.Nil(t, err)
}

func (r *restClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
}

func (r *restClientTest) TestChangeOperations(t *testing.T) {
	ctx := context.TODO()

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//upsert change
	changeCreated, err := r.client.ChangeUpsert(ctx, data.ChangePartial{
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
	changeRead, err := r.client.ChangeRead(ctx, changeId)
	assert.Nil(t, err)
	assert.NotNil(t, changeRead)
	assert.Equal(t, changeCreated, changeRead)

	//delete change
	err = r.client.ChangeDelete(ctx, changeId)
	assert.Nil(t, err)

	//read change
	changeRead, err = r.client.ChangeRead(ctx, changeId)
	assert.NotNil(t, err)
	assert.Nil(t, changeRead)
}

func (r *restClientTest) TestChangeStreaming(t *testing.T) {
	var change *data.Change

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	changeReceived := make(chan struct{})

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//register handler
	handlerId, err := r.client.HandlerCreate(func(changes ...*data.Change) error {
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
	change, err = r.client.ChangeUpsert(ctx, data.ChangePartial{
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
	err = r.client.HandlerDelete(handlerId)
	assert.Nil(t, err)
}

func TestEmployeesRestClient(t *testing.T) {
	r := newRestclientTest()

	r.Initialize(t)
	defer r.Shutdown(t)

	t.Run("Test Change Operations", r.TestChangeOperations)
	//KIM: this test is disabled because reading via websockets
	// is Janky
	// t.Run("Test Change Streaming", r.TestChangeStreaming)
}
