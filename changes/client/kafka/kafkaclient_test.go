package kafkaclient_test

import (
	"context"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	kafkaclient "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	restclient "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/google/uuid"
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

type kafkaClientTest struct {
	changesHandler interface {
		internal.Initializer
		internal.Configurer
		client.Handler
	}
	changesClient interface {
		client.Client
		internal.Initializer
		internal.Configurer
	}
	logger internal.Configurer
}

func newKafkaClientTest() *kafkaClientTest {
	logger := internal_logger.New()
	kafkaClient := internal_kafka.New()
	kafkaClient.SetUtilities(logger)

	//
	changesClient := restclient.New()
	changesClient.SetUtilities(logger)

	//
	changesHandler := kafkaclient.New()
	changesHandler.SetUtilities(logger)
	changesHandler.SetParameters(kafkaClient)
	return &kafkaClientTest{
		logger:         logger,
		changesHandler: changesHandler,
		changesClient:  changesClient,
	}
}

func (r *kafkaClientTest) Initialize(t *testing.T) {
	err := r.logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	assert.Nil(t, err)

	//create topic
	err = r.changesClient.Configure(configChangesRest)
	assert.Nil(t, err)
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.changesHandler.Configure(configChangesKafka, configKafka)
	assert.Nil(t, err)
	err = r.changesHandler.Initialize()
	assert.Nil(t, err)
}

func (r *kafkaClientTest) Shutdown(t *testing.T) {
	r.changesClient.Shutdown()
	r.changesHandler.Shutdown()
}

func (r *kafkaClientTest) TestChangeStreaming(t *testing.T) {
	var change *data.Change

	ctx := context.TODO()
	changeReceived := make(chan struct{})

	//generate dynamic constants
	dataId, dataVersion := generateId(), rand.Intn(1000)
	whenChanged, serviceName := time.Now().UnixNano(), randomString()
	dataType := "test"

	//register handler
	handlerId, err := r.changesHandler.HandlerCreate(func(changes ...*data.Change) error {
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

	time.Sleep(10 * time.Second)

	//upsert change
	change, err = r.changesClient.ChangeUpsert(ctx, data.ChangePartial{
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

	//unregister handler
	err = r.changesHandler.HandlerDelete(handlerId)
	assert.Nil(t, err)
}

func TestChangesKafkaClient(t *testing.T) {
	r := newKafkaClientTest()

	r.Initialize(t)
	defer r.Shutdown(t)

	t.Run("Change Streaming", r.TestChangeStreaming)
	//TODO: test consumer group?
}
