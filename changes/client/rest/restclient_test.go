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

	internal_cache "github.com/antonio-alexander/go-bludgeon/changes/internal/cache"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	queueSize int    = 10
	restPort  string = "8080"
)

var (
	config       = restclient.NewConfiguration()
	configCache  = internal_cache.NewConfiguration()
	configLogger = new(logger.Configuration)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.Rest.Port = restPort
	configLogger.Level = logger.Trace
	configCache.Default()
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
	logger interface {
		logger.Logger
		logger.Printer
		internal.Configurer
	}
	client interface {
		client.Client
		client.Handler
		internal.Configurer
		internal.Initializer
	}
	cache interface {
		internal_cache.Cache
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
	}
	queue interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.EnqueueInFronter
		goqueue.Length
		goqueue.Event
		goqueue.Peeker
		finite.EnqueueLossy
		finite.Resizer
		finite.Capacity
	}
}

func newRestclientTest() *restClientTest {
	logger := logger.New()
	queue := finite.New(queueSize)
	cache := internal_cache.New()
	cache.SetUtilities(logger)
	client := restclient.New()
	client.SetUtilities(logger)
	client.SetParameters(cache, queue)
	return &restClientTest{
		client: client,
		cache:  cache,
		queue:  queue,
		logger: logger,
	}
}

func (r *restClientTest) assertChange(t *testing.T, ctx context.Context, change *data.Change) func() bool {
	return func() bool {
		dataId, dataType := change.DataId, change.DataType
		dataAction, dataServiceName := change.DataAction, change.DataServiceName
		tStop := time.After(10 * time.Second)
		tCheck := time.NewTicker(time.Second)
		defer tCheck.Stop()
		checkChange := func() bool {
			changes, err := r.client.ChangesRead(ctx, data.ChangeSearch{
				DataIds:      []string{dataId},
				Types:        []string{dataType},
				Actions:      []string{dataAction},
				ServiceNames: []string{dataServiceName},
			})
			if err != nil {
				return false
			}
			if len(changes) == 0 {
				return false
			}
			for _, change := range changes {
				switch {
				case change.DataId != dataId:
					t.Logf("change: %s, dataId(%s) doesn't match", change.Id, dataId)
					continue
				case change.DataType != dataType:
					t.Logf("change: %s, dataType(%s) doesn't match", change.Id, dataType)
					continue
				case change.DataAction != dataAction:
					t.Logf("change: %s, dataAction(%s) doesn't match", change.Id, dataAction)
					continue
				case change.DataServiceName != dataServiceName:
					t.Logf("change: %s, serviceName(%s) doesn't match", change.Id, dataServiceName)
					continue
				}
				return true
			}
			return false
		}
		for {
			select {
			case <-tStop:
				return checkChange()
			case <-tCheck.C:
				if checkChange() {
					return true
				}
			}
		}
	}
}

func (r *restClientTest) Initialize(t *testing.T) {
	//configure
	err := r.logger.Configure(configLogger)
	assert.Nil(t, err)
	err = r.client.Configure(config)
	assert.Nil(t, err)
	err = r.cache.Configure(configCache)
	assert.Nil(t, err)
	//initialize
	err = r.cache.Initialize()
	assert.Nil(t, err)
	err = r.client.Initialize()
	assert.Nil(t, err)
}

func (r *restClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
	r.cache.Shutdown()
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

func (r *restClientTest) TestRegistrationOperations(t *testing.T) {
	ctx := context.TODO()

	// generate dynamic constants
	registrationId := randomString()
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	// create registration
	err := r.client.RegistrationUpsert(ctx, registrationId)
	assert.Nil(t, err)

	// upsert change
	changeCreated, err := r.client.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)

	// read registration changes
	changes, err := r.client.RegistrationChangesRead(ctx, registrationId)
	assert.Nil(t, err)
	assert.NotNil(t, changes)
	assert.Contains(t, changes, changeCreated)

	// acknowledge change
	err = r.client.RegistrationChangeAcknowledge(ctx, registrationId, changeCreated.Id)
	assert.Nil(t, err)

	// read registration changes
	changes, err = r.client.RegistrationChangesRead(ctx, registrationId)
	assert.Nil(t, err)
	assert.Empty(t, changes)

	// delete registration
	err = r.client.RegistrationDelete(ctx, registrationId)
	assert.Nil(t, err)
}

func (r *restClientTest) TestChangeUpsertQueue(t *testing.T) {
	ctx := context.TODO()

	//initialize with the cache disabled and the queue enabled
	config.DisableCache, config.DisableQueue = true, false
	config.Rest.Port = "7999" //this is an invalid port to force use of the queue
	r.Initialize(t)
	defer r.Shutdown(t)

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//upsert change
	changePartial := data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	}
	changeCreated, err := r.client.ChangeUpsert(ctx, changePartial)
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	assert.Empty(t, changeCreated.Id)

	// wait for two cycles for additional code coverage
	time.Sleep(2 * config.UpsertQueueRate)

	//shutdown
	r.Shutdown(t)

	// validate item in the queue
	item, overflow := r.queue.PeekHead()
	assert.False(t, overflow)
	assert.NotNil(t, item)
	assert.Equal(t, item, changePartial)

	//re-initialize with a valid port
	config.DisableCache, config.DisableQueue = true, false
	config.Rest.Port = restPort //this is an invalid port to force use of the queue
	r.Initialize(t)

	//wait for the go routine to start
	time.Sleep(time.Second)

	//wait for the change to be upserted by underlying code/queue
	assert.Condition(t, r.assertChange(t, ctx, changeCreated))
}

func (r *restClientTest) TestChangeCache(t *testing.T) {
	ctx := context.TODO()

	//initialize
	config.DisableCache, config.DisableQueue = false, true
	r.Initialize(t)
	defer r.Shutdown(t)

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
	changeId := changeCreated.Id

	//validate change cached
	changeCached := r.cache.Read(changeId)
	assert.Equal(t, changeCreated, changeCached)

	//delete change
	err = r.client.ChangeDelete(ctx, changeId)
	assert.Nil(t, err)

	//validate change removed
	changeCached = r.cache.Read(changeId)
	assert.Nil(t, changeCached)
}

func testEmployeesRestClient(t *testing.T) {
	r := newRestclientTest()

	config.DisableCache, config.DisableQueue = true, true
	r.Initialize(t)
	t.Run("Test Change Operations", r.TestChangeOperations)
	t.Run("Test Registration Operations", r.TestRegistrationOperations)
	//KIM: this test is disabled because reading via websockets
	// is Janky
	// t.Run("Test Change Streaming", r.TestChangeStreaming)
	r.Shutdown(t)

	//these internally initialize/shutdown
	t.Run("Test Change Upsert Queue", r.TestChangeUpsertQueue)
	t.Run("Test Change Cache", r.TestChangeCache)
}

func TestEmployeesRestClient(t *testing.T) {
	testEmployeesRestClient(t)
}
