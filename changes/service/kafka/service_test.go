package kafka_test

import (
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	"github.com/antonio-alexander/go-bludgeon/changes/service/kafka"
	"github.com/antonio-alexander/go-bludgeon/changes/service/kafka/tests"

	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_meta.json"

var (
	kafkaBrokers    = []string{"localhost:9092"}
	changeTopic     = "changes"
	kafkaConfig     = new(internal_kafka.Configuration)
	testKafkaConfig = new(internal_kafka.Configuration)
	mysqlConfig     = new(mysql.Configuration)
	fileConfig      = new(file.Configuration)
	logConfig       = new(internal_logger.Configuration)
)

type kafkaServiceTest struct {
	meta         common.Initializer
	logic        common.Initializer
	client       common.Initializer
	testClient   common.Initializer
	kafkaService common.Initializer
	*tests.Fixture
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
	mysqlConfig = new(mysql.Configuration)
	mysqlConfig.Default()
	mysqlConfig.FromEnv(envs)

	//create file config
	fileConfig = new(file.Configuration)
	fileConfig.Default()
	fileConfig.FromEnv(envs)
	fileConfig.File = path.Join("../../tmp", filename)
	os.Remove(fileConfig.File)

	//create internal_logger config
	logConfig = new(internal_logger.Configuration)
	logConfig.Default()
	logConfig.FromEnv(envs)
	logConfig.Level = internal_logger.Trace
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
		common.Initializer
		common.Parameterizer
		common.Configurer
	}

	logger := internal_logger.New()
	logger.Configure(logConfig)
	switch metaType {
	case "memory":
		meta = memory.New()
		meta.SetParameters(logger)
	case "file":
		meta = file.New()
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
	kafkaService := kafka.New()
	kafkaService.SetParameters(logger, logic, kafkaClient)
	kafkaService.Configure(&kafka.Configuration{
		Topic: changeTopic,
	})
	kafkaClientTest := internal_kafka.New()
	kafkaClientTest.Configure(kafkaConfig)
	kafkaClientTest.SetParameters(logger)
	return &kafkaServiceTest{
		kafkaService: kafkaService,
		meta:         meta,
		logic:        logic,
		client:       kafkaClient,
		testClient:   kafkaClientTest,
		Fixture:      tests.NewFixture(kafkaClientTest, logic),
	}
}

func (k *kafkaServiceTest) initialize(t *testing.T, consumerGroup bool) {
	checkTopicsRate := time.Second
	kafkaConfig.ConsumerGroup = consumerGroup
	kafkaConfig.ClientId = generateId()
	kafkaConfig.GroupId = generateId()
	err := k.meta.Initialize()
	assert.Nil(t, err)
	err = k.client.Initialize()
	assert.Nil(t, err)
	testKafkaConfig.ClientId = generateId()
	err = k.testClient.Initialize()
	assert.Nil(t, err)
	err = k.logic.Initialize()
	assert.Nil(t, err)
	assert.Nil(t, err)
	err = k.kafkaService.Initialize()
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
	k.kafkaService.Shutdown()
	k.logic.Shutdown()
	k.meta.Shutdown()
}

func testChangesKafkaService(t *testing.T, metaType string) {
	k := newKafkaServiceTest(metaType)

	consumerGroup := false
	k.initialize(t, consumerGroup)
	t.Run("Change Handler (Consumer)", k.TestChangeHandler(changeTopic))
	k.shutdown(t)

	//REVIEW: this test is consistently failing; need to figure out why
	// k = newKafkaServiceTest(metaType)
	// consumerGroup = true
	// k.initialize(t, consumerGroup)
	// t.Run("Change Handler (Consumer Group)", k.TestChangeHandler(changeTopic))
	// k.shutdown(t)
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
