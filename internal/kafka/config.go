package kafka

import (
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	EnvNameKafkaBrokers       string = "BLUDGEON_KAFKA_BROKERS"
	EnvNameKafkaClientId      string = "BLUDGEON_KAFKA_CLIENT_ID"
	EnvNameKafkaGroupId       string = "BLUDGEON_KAFKA_GROUP_ID"
	EnvNameKafkaConsumerGroup string = "BLUDGEON_KAFKA_CONSUMER_GROUP"
	EnvNameKafkaEnableLog     string = "BLUDGEON_KAFKA_ENABLE_LOG"
)

const (
	NoBrokersConfigured  string = "no brokers configured"
	NoClientIdConfigured string = "no client id configured"
	NoGroupIdConfigured  string = "no group id configured"
)

var (
	ErrNoBrokersConfigured  = errors.New(NoBrokersConfigured)
	ErrNoClientIdConfigured = errors.New(NoClientIdConfigured)
	ErrNoGroupIdConfigured  = errors.New(NoGroupIdConfigured)
)

type Configuration struct {
	Brokers       []string `json:"brokers"`
	ClientId      string   `json:"client_id"`
	GroupId       string   `json:"group_id"`
	EnableLog     bool     `json:"enable_log"`
	ConsumerGroup bool     `json:"consumer_group"`
}

func (c *Configuration) Default() {
	c.ClientId = uuid.Must(uuid.NewRandom()).String()
	c.GroupId = uuid.Must(uuid.NewRandom()).String()
	c.EnableLog = false
	c.ConsumerGroup = false
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if brokers, ok := envs[EnvNameKafkaBrokers]; ok && brokers != "" {
		c.Brokers = strings.Split(brokers, ",")
	}
	if clientId, ok := envs[EnvNameKafkaClientId]; ok && clientId != "" {
		c.ClientId = clientId
	}
	if groupId, ok := envs[EnvNameKafkaGroupId]; ok && groupId != "" {
		c.GroupId = groupId
	}
	if s, ok := envs[EnvNameKafkaConsumerGroup]; ok && s != "" {
		c.ConsumerGroup, _ = strconv.ParseBool(s)
	}
	if s, ok := envs[EnvNameKafkaEnableLog]; ok && s != "" {
		c.EnableLog, _ = strconv.ParseBool(s)
	}
}

func (c *Configuration) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrNoBrokersConfigured
	}
	if c.ConsumerGroup && c.GroupId == "" {
		return ErrNoGroupIdConfigured
	}
	if c.ClientId == "" {
		return ErrNoClientIdConfigured
	}
	return nil
}

func (c *Configuration) ToSarama() ([]string, *sarama.Config) {
	config := sarama.NewConfig()
	config.ClientID = c.ClientId
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ChannelBufferSize = 1024
	config.Consumer.Return.Errors = true
	return c.Brokers, config
}
