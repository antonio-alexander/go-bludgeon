package rest_test

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	restclient "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	tests "github.com/antonio-alexander/go-bludgeon/changes/client/tests"

	internal_cache "github.com/antonio-alexander/go-bludgeon/changes/internal/cache"
	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	internal_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/client"
	internal_mock "github.com/antonio-alexander/go-bludgeon/pkg/rest/client/mock"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"

	"github.com/stretchr/testify/assert"
)

const (
	queueSize int    = 10
	restPort  string = "8080"
)

var (
	config       = new(restclient.Configuration)
	configCache  = new(internal_cache.Configuration)
	configLogger = new(internal_logger.Configuration)
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
	configLogger.Level = internal_logger.Trace
	configCache.Default()
}

type restClientTest struct {
	logger interface {
		internal_logger.Logger
		internal_logger.Printer
		common.Configurer
	}
	client interface {
		client.Client
		client.Handler
		common.Configurer
		common.Initializer
		common.Parameterizer
	}
	cache interface {
		internal_cache.Cache
		common.Configurer
		common.Initializer
		common.Parameterizer
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
	restClient interface {
		internal_rest.Client
		common.Parameterizer
		common.Configurer
	}
	mockRestclient interface {
		internal_rest.Client
		internal_mock.Mock
	}
	*tests.Fixture
}

func newRestclientTest() *restClientTest {
	logger := internal_logger.New()
	queue := finite.New(queueSize)
	cache := internal_cache.New()
	cache.SetUtilities(logger)
	client := restclient.New()
	client.SetUtilities(logger)
	client.SetParameters(cache, queue)
	restClient := internal_rest.New()
	mockRestclient := internal_mock.New()
	return &restClientTest{
		client:         client,
		cache:          cache,
		queue:          queue,
		logger:         logger,
		restClient:     restClient,
		mockRestclient: mockRestclient,
		Fixture:        tests.NewFixture(client, cache, queue, mockRestclient),
	}
}

func (r *restClientTest) Initialize(t *testing.T, mock bool) {
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

	if mock {
		r.client.SetParameters(r.mockRestclient)
	} else {
		r.client.SetParameters(r.restClient)
	}
}

func (r *restClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
	r.cache.Shutdown()
}

func testEmployeesRestClient(t *testing.T) {
	r := newRestclientTest()

	config.DisableCache, config.DisableQueue = true, true
	r.Initialize(t, false)
	t.Run("Test Change Operations", r.TestChangeOperations)
	t.Run("Test Registration Operations", r.TestRegistrationOperations)
	//KIM: this test is disabled because reading via websockets
	// is Janky
	// t.Run("Test Change Streaming", r.TestChangeStreaming)
	r.Shutdown(t)

	config.DisableCache, config.DisableQueue = false, true
	r.Initialize(t, true)
	t.Run("Test Change Upsert Queue", r.TestChangeUpsertQueue)
	t.Run("Test Change Cache", r.TestChangeCache)
	r.Shutdown(t)
}

func TestEmployeesRestClient(t *testing.T) {
	testEmployeesRestClient(t)
}
