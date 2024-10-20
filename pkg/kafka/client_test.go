package kafka_test

import (
	"math/rand"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
	kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	kafkaBrokers = []string{"localhost:9092"}
	kafkaConfig  *kafka.Configuration
	saramaConfig *sarama.Config
)

type kafkaClientTest struct {
	kafkaClient interface {
		kafka.Client
		common.Initializer
		common.Configurer
	}
	logger interface {
		common.Configurer
	}
	saramaClient sarama.Client
}

func init() {
	//get environment
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}

	//create kafka config
	kafkaConfig = new(kafka.Configuration)
	kafkaConfig.Default()
	kafkaConfig.Brokers = kafkaBrokers
	kafkaConfig.FromEnv(envs)

	//create sarama config
	saramaConfig = sarama.NewConfig()
	saramaConfig.ClientID = generateId()
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ChannelBufferSize = 1024
	saramaConfig.Consumer.Return.Errors = true

	//seed
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

func newKafkaClientTest() *kafkaClientTest {
	logger := logger.New()
	kafkaClient := kafka.New()
	kafkaClient.SetUtilities(logger)
	return &kafkaClientTest{
		kafkaClient: kafkaClient,
		logger:      logger,
	}
}

func (k *kafkaClientTest) initialize(t *testing.T, consumerGroup bool) {
	err := k.logger.Configure(&logger.Configuration{
		Level:  logger.Trace,
		Prefix: "bludgeon_kafka_client_test",
	})
	assert.Nil(t, err)
	kafkaConfig.ConsumerGroup = consumerGroup
	// kafkaConfig.EnableLog = true
	err = k.kafkaClient.Configure(nil, nil, kafkaConfig)
	assert.Nil(t, err)
	err = k.kafkaClient.Initialize()
	assert.Nil(t, err)
	saramaClient, err := sarama.NewClient(kafkaConfig.Brokers, saramaConfig)
	assert.Nil(t, err)
	k.saramaClient = saramaClient
}

func (k *kafkaClientTest) shutdown(t *testing.T) {
	if err := k.saramaClient.Close(); err != nil {
		t.Logf("error whiel closing sarama client: %s", err)
	}
	k.kafkaClient.Shutdown()
}

func (k *kafkaClientTest) TestPublishConsumer(t *testing.T) {
	var wg sync.WaitGroup

	//generate dynamic constants
	testTopic := "test.kafka-client"
	testBytes := []byte(randomString(25))
	start, stopper := make(chan struct{}), make(chan struct{})
	messageReceived := make(chan struct{})

	//create sarama consumer and subscribe to topic/partitions
	consumer, err := sarama.NewConsumerFromClient(k.saramaClient)
	assert.Nil(t, err)
	defer func() {
		if err := consumer.Close(); err != nil {
			t.Logf("error while closing consumer: %s", err)
		}
	}()
	partitions, err := consumer.Partitions(testTopic)
	assert.Nil(t, err)
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()
			//consume the partition and messages
			partitionConsumer, err := consumer.ConsumePartition(testTopic, partition, sarama.OffsetNewest)
			assert.Nil(t, err)
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					t.Logf("error while closing partition consumer: %s", err)
				}
			}()
			<-start
			for {
				select {
				case <-messageReceived:
					return
				case <-stopper:
					return
				case m := <-partitionConsumer.Messages():
					if reflect.DeepEqual(m.Value, testBytes) {
						select {
						default:
							close(messageReceived)
						case <-messageReceived:
						}
					}
				}
			}
		}(partition)
	}

	//publish messsage continuously
	wg.Add(1)
	go func() {
		defer wg.Done()

		tPublish := time.NewTicker(time.Second)
		defer tPublish.Stop()
		<-start
		for {
			select {
			case <-stopper:
				return
			case <-tPublish.C:
				err := k.kafkaClient.Publish(testTopic, testBytes)
				assert.Nil(t, err)
			}
		}
	}()

	//confirm receipt of message
	close(start)
	select {
	case <-messageReceived:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm message received")
	}

	//clean up resources
	close(stopper)
	wg.Wait()
}

func (k *kafkaClientTest) TestSubscribeConsumer(t *testing.T) {
	var wg sync.WaitGroup

	//generate dynamic constants
	testTopic := "test.kafka-client"
	start, stopper := make(chan struct{}), make(chan struct{})
	testBytes := []byte(randomString(25))
	messageReceived := make(chan struct{})

	//subscribe
	handlerId, err := k.kafkaClient.Subscribe(testTopic, func(topic string, bytes []byte) {
		if topic != testTopic {
			return
		}
		if reflect.DeepEqual(bytes, testBytes) {
			select {
			default:
				close(messageReceived)
			case <-messageReceived:
			}
		}
	})
	assert.Nil(t, err)
	defer func() { k.kafkaClient.Unsubscribe(testTopic, handlerId) }()

	//periodically publish data
	wg.Add(1)
	go func() {
		defer wg.Done()

		tPublish := time.NewTicker(time.Second)
		defer tPublish.Stop()
		<-start
		for {
			select {
			case <-stopper:
				return
			case <-messageReceived:
				return
			case <-tPublish.C:
				err := k.kafkaClient.Publish(testTopic, testBytes)
				assert.Nil(t, err)
			}
		}
	}()

	//start the go routines
	close(start)
	select {
	case <-messageReceived:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm message received")
	}

	//clean up
	close(stopper)
	wg.Wait()
}

func TestKafkaClient(t *testing.T) {
	k := newKafkaClientTest()

	consumerGroup := false
	k.initialize(t, consumerGroup)
	t.Run("Test Publish (Consumer)", k.TestPublishConsumer)
	t.Run("Test Subscribe (Consumer)", k.TestSubscribeConsumer)
	k.shutdown(t)

	consumerGroup = true
	k.initialize(t, consumerGroup)
	t.Run("Test Subscribe (Consumer Group)", k.TestSubscribeConsumer)
	k.shutdown(t)
}
