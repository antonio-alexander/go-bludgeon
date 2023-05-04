package rest_test

import (
	"os"
	"strings"
	"testing"
	"time"

	restclient "github.com/antonio-alexander/go-bludgeon/timers/client/rest"
	tests "github.com/antonio-alexander/go-bludgeon/timers/client/tests"
	cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	employeesclientrest "github.com/antonio-alexander/go-bludgeon/employees/client/rest"
	kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"

	"github.com/stretchr/testify/assert"
)

var (
	config               = new(restclient.Configuration)
	configCache          = new(cache.Configuration)
	configKafka          = new(kafka.Configuration)
	configEmployeeClient = new(employeesclientrest.Configuration)
	configChangesClient  = changesclientrest.NewConfiguration()
	configChangesHandler = new(changesclientkafka.Configuration)
	configLogger         = &logger.Configuration{
		Prefix: "bludgeon_rest_client_test",
		Level:  logger.Trace,
	}
)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.Port = "8080"
	config.ChangeRateRead = time.Second
	config.ChangeRateRegistration = time.Second
	config.ChangesRegistrationId = "timers_client_rest"
	configEmployeeClient.Default()
	configEmployeeClient.FromEnv(envs)
	configEmployeeClient.Port = "9000"
	configChangesClient.Default()
	configChangesClient.FromEnv(envs)
	configChangesClient.Rest.Timeout = 30 * time.Second
	configChangesClient.Rest.Port = "9020"
	configChangesHandler.Default()
	configChangesHandler.FromEnv(envs)
	configCache.Default()
	configCache.FromEnvs(envs)
	configCache.Debug = true
	configKafka.Default()
	configKafka.FromEnv(envs)
	configKafka.Brokers = []string{"localhost:9092"}
	configKafka.ConsumerGroup = false
}

type restClientTest struct {
	client interface {
		internal.Initializer
		internal.Configurer
	}
	cache interface {
		internal.Initializer
		internal.Configurer
	}
	employeesClient interface {
		internal.Initializer
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
	}
	changesHandler interface {
		internal.Initializer
		internal.Configurer
	}
	kafkaClient interface {
		internal.Initializer
		internal.Configurer
	}
	*tests.TestFixture
}

func newRestClientTest() *restClientTest {
	logger := logger.New()
	logger.Configure(configLogger)
	client := restclient.New()
	cache := cache.New()
	kafkaClient := kafka.New()
	employeesClient := employeesclientrest.New()
	changesClient, changesHandler := changesclientrest.New(), changesclientkafka.New()
	kafkaClient.SetUtilities(logger)
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	client.SetUtilities(logger)
	cache.SetUtilities(logger)
	employeesClient.SetUtilities(logger)
	client.SetParameters(cache, changesClient, changesHandler)
	changesHandler.SetParameters(kafkaClient)
	return &restClientTest{
		client:          client,
		cache:           cache,
		employeesClient: employeesClient,
		changesClient:   changesClient,
		changesHandler:  changesHandler,
		kafkaClient:     kafkaClient,
		TestFixture:     tests.NewTestFixture(client, cache, employeesClient, changesClient),
	}
}

func (r *restClientTest) Initialize(t *testing.T) {
	err := r.client.Configure(config)
	assert.Nil(t, err)
	err = r.employeesClient.Configure(configEmployeeClient)
	assert.Nil(t, err)
	err = r.changesClient.Configure(configChangesClient)
	assert.Nil(t, err)
	err = r.changesHandler.Configure(configChangesHandler, configKafka)
	assert.Nil(t, err)
	err = r.kafkaClient.Configure(configKafka)
	assert.Nil(t, err)
	err = r.cache.Configure(configCache)
	assert.Nil(t, err)
	err = r.client.Initialize()
	assert.Nil(t, err)
	err = r.employeesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesHandler.Initialize()
	assert.Nil(t, err)
	//REVIEW: would be nice to have some way to block
	// until the registration is upserted
	// err = r.kafkaClient.Initialize()
	// assert.Nil(t, err)
	err = r.cache.Initialize()
	assert.Nil(t, err)
}

func (r *restClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
	r.employeesClient.Shutdown()
	r.changesClient.Shutdown()
	r.changesHandler.Shutdown()
	r.kafkaClient.Shutdown()
	r.cache.Shutdown()
}

func testTimersRestClient(t *testing.T) {
	r := newRestClientTest()

	config.DisableCache = true
	r.Initialize(t)
	t.Run("Test Timer Operations", r.TestTimers)
	r.Shutdown(t)

	config.DisableCache = false
	r.Initialize(t)
	t.Run("Test Timer Cache", r.TestTimerCache)
	r.Shutdown(t)
}

func TestTimersRestClient(t *testing.T) {
	testTimersRestClient(t)
}
