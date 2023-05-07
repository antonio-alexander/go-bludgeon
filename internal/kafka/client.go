package kafka

import (
	"context"
	"encoding"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type kafka struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	logger.Printer
	*topicHandlers
	ctx           context.Context
	cancel        context.CancelFunc
	client        sarama.Client
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	consumer      sarama.Consumer
	initialized   bool
	configured    bool
	config        *Configuration
}

var _ sarama.ConsumerGroupHandler = &kafka{}

func New() interface {
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
	Client
} {
	nullLogger := logger.NewNullLogger()
	k := &kafka{
		topicHandlers: newTopicHandlers(),
		Logger:        nullLogger,
		Printer:       nullLogger,
	}
	return k
}

func (k *kafka) launchConsumer(topic string, stopper chan struct{}) error {
	var partitionConsumers []sarama.PartitionConsumer

	//REVIEW: what's the likelihood that new partitions are created?
	// should we continually check for new partitions for a given topic?
	partitions, err := k.consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, partition := range partitions {
		partitionConsumer, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			for _, partitionConsumer := range partitionConsumers {
				partitionConsumer.Close()
			}
			return err
		}
		partitionConsumers = append(partitionConsumers, partitionConsumer)
	}
	for i, partitionConsumer := range partitionConsumers {
		started := make(chan struct{})
		k.Add(1)
		go func(partition int32, partitionConsumer sarama.PartitionConsumer) {
			defer k.Done()
			defer func() {
				k.Trace(logAlias+"stopped partition consumer for topic \"%s\" (%d)", topic, partition)
			}()

			k.Trace(logAlias+"launched partition consumer for topic \"%s\" (%d)", topic, partition)
			chMessages := partitionConsumer.Messages()
			close(started)
			for {
				select {
				case <-k.ctx.Done():
					return
				case msg := <-chMessages:
					wg := new(sync.WaitGroup)
					for _, handler := range k.readHandlers(topic) {
						wg.Add(1)
						go func(handler HandleFx) {
							defer wg.Done()

							handler(topic, msg.Value)
						}(handler)
					}
					wg.Wait()
				}
			}
		}(partitions[i], partitionConsumer)
		<-started
	}
	return nil
}

func (k *kafka) launchConsumerGroup(topic string, stopper chan struct{}) error {
	chErr := make(chan error)
	ctx, cancel := context.WithCancel(k.ctx)
	started := make(chan struct{})
	k.Add(2)
	go func() {
		defer k.Done()

		select {
		case <-k.ctx.Done():
			cancel()
		case <-stopper:
			cancel()
		case <-ctx.Done():
		}
	}()
	go func(ctx context.Context) {
		defer k.Done()
		defer func() {
			k.Trace(logAlias+"stopped consumer group for topic \"%s\"", topic)
		}()

		topics := []string{topic}
		close(started)
		k.Trace(logAlias+"launched partition consumer group for topic \"%s\"", topic)
		//KIM: this blocks, so it needs to be in a go routine
		if err := k.consumerGroup.Consume(ctx, topics, k); err != nil {
			cancel()
			select {
			default:
				chErr <- err
			case <-chErr:
			}
		}
	}(ctx)
	<-started
	select {
	default:
		close(chErr)
		return nil
	case err := <-chErr:
		return err
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (k *kafka) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (k *kafka) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (k *kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	wg := new(sync.WaitGroup)
	topic := claim.Topic()
	handlers := k.readHandlers(topic)
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
		for _, handler := range handlers {
			wg.Add(1)
			go func(handler HandleFx, bytes []byte) {
				defer wg.Done()
				handler(topic, bytes)
			}(handler, msg.Value)
		}
	}
	wg.Wait()
	return nil
}

func (k *kafka) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (k *kafka) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case interface {
			logger.Logger
			logger.Printer
		}:
			k.Logger = p
			k.Printer = p
		case logger.Logger:
			k.Logger = p
		case logger.Printer:
			k.Printer = p
		}
	}
}

func (k *kafka) Configure(items ...interface{}) error {
	k.Lock()
	defer k.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			c = new(Configuration)
			c.Default()
			c.FromEnv(v)
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	k.config = c
	k.configured = true
	return nil
}

func (k *kafka) Initialize() error {
	k.Lock()
	defer k.Unlock()

	if k.initialized {
		return errors.New("already initialized")
	}
	if !k.configured {
		return errors.New("not configured")
	}
	if k.config.EnableLog {
		sarama.Logger = k
	}
	client, err := sarama.NewClient(k.config.ToSarama())
	if err != nil {
		return err
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return err
	}
	if k.config.ConsumerGroup {
		consumerGroup, err := sarama.NewConsumerGroupFromClient(k.config.GroupId, client)
		if err != nil {
			return err
		}
		k.consumerGroup = consumerGroup
	} else {
		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			return err
		}
		k.consumer = consumer
	}
	k.client, k.producer = client, producer
	k.ctx, k.cancel = context.WithCancel(context.Background())
	k.initialized = true
	return nil
}

func (k *kafka) Shutdown() {
	k.Lock()
	defer k.Unlock()

	if !k.initialized {
		return
	}
	k.cancel()
	k.deleteTopics()
	k.Wait()
	if k.config.ConsumerGroup {
		if err := k.consumerGroup.Close(); err != nil {
			k.Error(logAlias+"error while closing consumer group: %s", err)
		}
	} else {
		if err := k.consumer.Close(); err != nil {
			k.Error(logAlias+"error while closing consumer: %s", err)
		}
	}
	if err := k.producer.Close(); err != nil {
		k.Error(logAlias+"error while closing producer: %s", err)
	}
	if err := k.client.Close(); err != nil {
		k.Error(logAlias+"error while closing client: %s", err)
	}
	k.producer, k.client = nil, nil
	k.consumerGroup, k.consumer = nil, nil
	k.initialized, k.configured = false, false
}

func (k *kafka) Publish(topic string, item interface{}) error {
	var byteEncoder sarama.Encoder

	if !k.initialized {
		return errors.New("not initialized")
	}

	switch v := item.(type) {
	default:
		return errors.New("unsupported type")
	case encoding.BinaryMarshaler:
		bytes, err := v.MarshalBinary()
		if err != nil {
			return err
		}
		byteEncoder = sarama.ByteEncoder(bytes)
	case []byte:
		byteEncoder = sarama.ByteEncoder(v)
	}
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.ByteEncoder{},
		Headers:   []sarama.RecordHeader{},
		Timestamp: time.Now(),
		Offset:    sarama.OffsetNewest,
		Partition: 0,
		Value:     byteEncoder,
	})
	return err
}

func (k *kafka) Subscribe(topic string, handler HandleFx) (string, error) {
	if !k.initialized {
		return "", errors.New("not initialized")
	}
	stopper, newTopic := k.upsertTopic(topic)
	if newTopic {
		k.Trace(logAlias+"subscribed to topic \"%s\"", topic)
		switch {
		default:
			if err := k.launchConsumer(topic, stopper); err != nil {
				return "", err
			}
		case k.config.ConsumerGroup:
			if err := k.launchConsumerGroup(topic, stopper); err != nil {
				return "", err
			}
		}
	}
	handlerId := k.writeHandler(topic, handler)
	k.Trace(logAlias+"subscribed to topic \"%s\", handlerId \"%s\"", topic, handlerId)
	return handlerId, nil
}

func (k *kafka) Unsubscribe(topic string, handlerIds ...string) {
	if !k.initialized {
		return
	}
	k.deleteHandlers(topic, handlerIds...)
	if len(handlerIds) > 0 {
		k.Trace(logAlias+"unsubscribed from topic \"%s\", handlerIds \"%s\"", topic, strings.Join(handlerIds, ","))
	} else {
		k.Trace(logAlias+"unsubscribed from topic \"%s\"", topic)
	}
}

func (k *kafka) Topics(regEx *regexp.Regexp) ([]string, error) {
	if !k.initialized {
		return nil, errors.New("not initialized")
	}
	if regEx != nil {
		var matchedtopics []string

		topics, err := k.client.Topics()
		if err != nil {
			return nil, err
		}
		for _, topic := range topics {
			if matches := regEx.FindAllString(topic, 1); len(matches) > 0 {
				matchedtopics = append(matchedtopics, topic)
			}
		}
		return matchedtopics, nil
	}
	return k.client.Topics()
}
