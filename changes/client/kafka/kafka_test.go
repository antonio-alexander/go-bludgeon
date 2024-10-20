package kafka_test

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	kafkaclient "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	restclient "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	tests "github.com/antonio-alexander/go-bludgeon/changes/client/tests"
	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/stretchr/testify/assert"
)

var (
	configKafka        = new(internal_kafka.Configuration)
	configChangesRest  = new(restclient.Configuration)
	configChangesKafka = new(kafkaclient.Configuration)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configKafka.Default()
	configKafka.FromEnv(envs)
	configKafka.Brokers = []string{"localhost:9092"}
	configKafka.ToSarama()
	configChangesRest.Default()
	configChangesRest.FromEnv(envs)
	configChangesRest.Rest.Port = "8080"
	configChangesKafka.Default()
}

type kafkaClientTest struct {
	changesHandler interface {
		common.Initializer
		common.Configurer
		client.Handler
	}
	changesClient interface {
		client.Client
		common.Initializer
		common.Configurer
	}
	logger common.Configurer
	*tests.Fixture
}

func newKafkaClientTest() *kafkaClientTest {
	logger := internal_logger.New()
	kafkaClient := internal_kafka.New()
	kafkaClient.SetUtilities(logger)
	changesClient := restclient.New()
	changesClient.SetUtilities(logger)
	changesHandler := kafkaclient.New()
	changesHandler.SetUtilities(logger)
	changesHandler.SetParameters(kafkaClient)
	return &kafkaClientTest{
		logger:         logger,
		changesHandler: changesHandler,
		changesClient:  changesClient,
		Fixture:        tests.NewFixture(changesHandler, changesClient),
	}
}

func (r *kafkaClientTest) Initialize(t *testing.T) {
	err := r.logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	assert.Nil(t, err)
	err = r.changesClient.Configure(configChangesRest)
	assert.Nil(t, err)
	err = r.changesHandler.Configure(configChangesKafka, configKafka)
	assert.Nil(t, err)
	//create topic
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesHandler.Initialize()
	assert.Nil(t, err)
}

func (r *kafkaClientTest) Shutdown(t *testing.T) {
	r.changesClient.Shutdown()
	r.changesHandler.Shutdown()
}

func TestChangesKafkaClient(t *testing.T) {
	r := newKafkaClientTest()

	r.Initialize(t)
	defer r.Shutdown(t)

	t.Run("Change Streaming", r.TestChangeStreaming)
	//TODO: test consumer group?
}
