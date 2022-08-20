package service_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	logic "github.com/antonio-alexander/go-bludgeon/changes/logic"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	file "github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	service_kafka "github.com/antonio-alexander/go-bludgeon/changes/service/kafka"
	internal "github.com/antonio-alexander/go-bludgeon/internal"

	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_meta.json"

var (
	kafkaBrokers    = []string{"localhost:9092"}
	serviceName     = generateId()
	changeTopic     = "change." + serviceName
	kafkaConfig     = new(internal_kafka.Configuration)
	testKafkaConfig = new(internal_kafka.Configuration)
	mysqlConfig     = new(internal_mysql.Configuration)
	fileConfig      = new(internal_file.Configuration)
	logConfig       = new(logger.Configuration)
)

type kafkaServiceTest struct {
	meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		internal.Shutdowner
	}
	logic interface {
		logic.Logic
		internal.Initializer
	}
	client interface {
		internal_kafka.Client
		internal.Initializer
	}
	testClient interface {
		internal_kafka.Client
		internal.Initializer
	}
	logger interface {
		logger.Logger
		logger.Printer
	}
	service interface {
		internal.Initializer
	}
}

func init() {
	//get environment
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}

	//create mysql config
	mysqlConfig = new(internal_mysql.Configuration)
	mysqlConfig.Default()
	mysqlConfig.FromEnv(envs)

	//create file config
	fileConfig = new(internal_file.Configuration)
	fileConfig.Default()
	fileConfig.FromEnv(envs)
	fileConfig.File = path.Join("../../tmp", filename)
	os.Remove(fileConfig.File)

	//create logger config
	logConfig = new(logger.Configuration)
	logConfig.Default()
	logConfig.FromEnv(envs)
	logConfig.Level = logger.Trace
	logConfig.Prefix = "test_service"

	//create kafka config
	kafkaConfig.Default()
	kafkaConfig.Brokers = kafkaBrokers
	kafkaConfig.FromEnv(envs)

	//create test kafka config
	testKafkaConfig.Default()
	testKafkaConfig.Brokers = kafkaBrokers
	testKafkaConfig.ConsumerGroup = false
	testKafkaConfig.FromEnv(envs)

	//seed
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func newKafkaServiceTest(metaType string) *kafkaServiceTest {
	var meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		internal.Initializer
		internal.Parameterizer
		internal.Configurer
	}

	logger := logger.New()
	logger.Configure(logConfig)
	switch metaType {
	case "memory":
		meta = memory.New()
		meta.SetParameters(logger)
	case "file":
		meta = file.New()
		meta.Configure()
		meta.SetParameters(logger)
		meta.Configure(fileConfig)
	case "mysql":
		meta = mysql.New()
		meta.SetParameters(logger)
		meta.Configure(mysqlConfig)
	}
	logic := logic.New()
	logic.SetParameters(logger, meta)
	kafkaClient := internal_kafka.New()
	kafkaClient.SetParameters(logger)
	kafkaClient.Configure(kafkaConfig)
	service := service_kafka.New()
	service.SetParameters(logger, logic, kafkaClient)
	service.Configure(&service_kafka.Configuration{
		Topics: kafkaBrokers,
	})
	kafkaClientTest := internal_kafka.New()
	kafkaClientTest.Configure(kafkaConfig)
	kafkaClientTest.SetParameters(logger)
	return &kafkaServiceTest{
		service:    service,
		meta:       meta,
		logic:      logic,
		client:     kafkaClient,
		testClient: kafkaClientTest,
		logger:     logger,
	}
}

func (k *kafkaServiceTest) initialize(t *testing.T, consumerGroup bool) {
	checkTopicsRate := time.Second
	kafkaConfig.ConsumerGroup = consumerGroup
	kafkaConfig.ClientId = generateId()
	kafkaConfig.GroupId = generateId()
	err := k.client.Initialize()
	assert.Nil(t, err)
	testKafkaConfig.ClientId = generateId()
	err = k.testClient.Initialize()
	assert.Nil(t, err)
	err = k.logic.Initialize()
	assert.Nil(t, err)
	assert.Nil(t, err)
	err = k.service.Initialize()
	assert.Nil(t, err)
	k.createTopic(t)
	//KIM: by sleeping two times the topics rate we ensure
	// that we're subscribed prior to
	time.Sleep(2 * checkTopicsRate)
}

func (k *kafkaServiceTest) createTopic(t *testing.T) {
	kafkaConfig.ClientId = uuid.Must(uuid.NewRandom()).String()
	admin, err := sarama.NewClusterAdmin(kafkaConfig.ToSarama())
	assert.Nil(t, err)
	defer func() {
		if err := admin.Close(); err != nil {
			t.Logf("error while closing admin: %s", err)
		}
	}()
	topics, err := admin.ListTopics()
	assert.Nil(t, err)
	topicFound := false
	for topicName := range topics {
		if topicName == changeTopic {
			topicFound = true
			break
		}
	}
	if topicFound {
		return
	}
	err = admin.CreateTopic(changeTopic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	assert.Nil(t, err)
}

func (k *kafkaServiceTest) shutdown(t *testing.T) {
	k.client.Shutdown()
	k.testClient.Shutdown()
	k.service.Shutdown()
	k.logic.Shutdown()
	k.meta.Shutdown()
}

func (k *kafkaServiceTest) testChangeHandler(t *testing.T) {
	var changeId string

	//generate dynamic constants
	ctx := context.TODO()
	changeReceived := make(chan struct{})
	serviceName := serviceName

	//subscribe for change
	handlerId, err := k.testClient.Subscribe(changeTopic, func(topic string, bytes []byte) {
		if len(bytes) == 0 {
			t.Logf("no bytes received")
			return
		}
		wrapper := &data.Wrapper{}
		if err := json.Unmarshal(bytes, wrapper); err != nil {
			t.Logf("error while unmarshalling json: %s", err)
			return
		}
		item, err := data.FromWrapper(wrapper)
		if err != nil {
			t.Logf("error during  FromWrapper: %s", err)
			return
		}
		switch v := item.(type) {
		case *data.ResponseChange:
			if changeId == v.Change.Id {
				select {
				default:
					close(changeReceived)
					changeId = v.Change.Id
					return
				case <-changeReceived:
				}
			}
		case *data.ChangeDigest:
			for _, change := range v.Changes {
				select {
				default:
					close(changeReceived)
					changeId = change.Id
					return
				case <-changeReceived:
				}
			}
		}
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, handlerId)

	//upsert change and validate change received
	dataId, version := generateId(), rand.Int()
	dataType, whenChanged := "employee", time.Now().UnixNano()
	changeCreated, err := k.logic.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &version,
		DataType:        &dataType,
		DataServiceName: &serviceName,
		WhenChanged:     &whenChanged,
	})
	assert.Nil(t, err)
	changeId = changeCreated.Id
	defer func() {
		k.logic.ChangesDelete(ctx, changeId)
	}()

	//validate that change received
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm change received")
	case <-changeReceived:
	}
}

func testChangesKafkaService(t *testing.T, metaType string) {
	k := newKafkaServiceTest(metaType)

	consumerGroup := false
	k.initialize(t, consumerGroup)
	t.Run("Change Handler (Consumer)", k.testChangeHandler)
	k.shutdown(t)

	k = newKafkaServiceTest(metaType)
	consumerGroup = true
	k.initialize(t, consumerGroup)
	t.Run("Change Handler (Consumer Group)", k.testChangeHandler)
	k.shutdown(t)
}

func TestChangesKafkaServiceMemory(t *testing.T) {
	testChangesKafkaService(t, "memory")
}

func TestChangesKafkaServiceFile(t *testing.T) {
	testChangesKafkaService(t, "file")
}

func TestChangesKafkaServiceMysql(t *testing.T) {
	testChangesKafkaService(t, "mysql")
}
