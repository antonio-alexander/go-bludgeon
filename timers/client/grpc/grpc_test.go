package grpc_test

import (
	"os"
	"strings"
	"testing"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	client "github.com/antonio-alexander/go-bludgeon/timers/client"
	grpcclient "github.com/antonio-alexander/go-bludgeon/timers/client/grpc"
	tests "github.com/antonio-alexander/go-bludgeon/timers/client/tests"
	cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	employeesclientrest "github.com/antonio-alexander/go-bludgeon/employees/client/rest"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_cache "github.com/antonio-alexander/go-bludgeon/timers/internal/cache"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	config               = new(grpcclient.Configuration)
	configCache          = new(internal_cache.Configuration)
	configKafka          = new(kafka.Configuration)
	configEmployeeClient = new(employeesclientrest.Configuration)
	configChangesClient  = changesclientrest.NewConfiguration()
	configChangesHandler = new(changesclientkafka.Configuration)
	configLogger         = &logger.Configuration{
		Prefix: "bludgeon_grpc_client_test",
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
	config.Options = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	config.Port = "8081"
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
	configKafka.Default()
	configKafka.FromEnv(envs)
	configKafka.Brokers = []string{"localhost:9092"}
}

type grpcClientTest struct {
	client interface {
		client.Client
		internal.Parameterizer
		internal.Configurer
		internal.Initializer
	}
	cache interface {
		internal_cache.Cache
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
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
	// kafkaClient interface {
	// 	internal.Initializer
	// 	internal.Configurer
	// }
	*tests.TestFixture
}

func newGrpcClientTest() *grpcClientTest {
	logger := logger.New()
	logger.Configure(configLogger)
	client := grpcclient.New()
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
	// changesHandler.SetParameters(kafkaClient)
	return &grpcClientTest{
		client:          client,
		cache:           cache,
		employeesClient: employeesClient,
		changesClient:   changesClient,
		changesHandler:  changesHandler,
		// kafkaClient:     kafkaClient,
		TestFixture: tests.NewTestFixture(client, cache, employeesClient, changesClient),
	}
}

func (r *grpcClientTest) Initialize(t *testing.T) {
	err := r.client.Configure(config)
	assert.Nil(t, err)
	err = r.employeesClient.Configure(configEmployeeClient)
	assert.Nil(t, err)
	err = r.changesClient.Configure(configChangesClient)
	assert.Nil(t, err)
	err = r.changesHandler.Configure(configChangesHandler, configKafka)
	assert.Nil(t, err)
	// err = r.kafkaClient.Configure(configKafka)
	// assert.Nil(t, err)
	err = r.client.Initialize()
	assert.Nil(t, err)
	err = r.employeesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesHandler.Initialize()
	assert.Nil(t, err)
	// err = r.kafkaClient.Initialize()
	// assert.Nil(t, err)
}

func (r *grpcClientTest) Shutdown(t *testing.T) {
	r.client.Shutdown()
	r.employeesClient.Shutdown()
	r.changesClient.Shutdown()
	r.changesHandler.Shutdown()
	// r.kafkaClient.Shutdown()
}

func testTimersGrpcClient(t *testing.T) {
	r := newGrpcClientTest()

	config.DisableCache = true
	r.Initialize(t)
	t.Run("Test Timer Operations", r.TestTimers)
	r.Shutdown(t)

	config.DisableCache = false
	r.Initialize(t)
	t.Run("Test Timer Cache", r.TestTimerCache)
	r.Shutdown(t)
}

func TestTimersGrpcClient(t *testing.T) {
	testTimersGrpcClient(t)
}
