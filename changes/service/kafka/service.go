package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/common"

	internal_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
)

type kafkaService struct {
	sync.RWMutex
	sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	internal_logger.Logger
	internal_kafka.Client
	logic       logic.Logic
	muHandlers  sync.RWMutex
	handlers    map[string]string
	initialized bool
	configured  bool
	config      *Configuration
}

func New() interface {
	common.Configurer
	common.Initializer
	common.Parameterizer
} {
	return &kafkaService{
		handlers: make(map[string]string),
		Logger:   internal_logger.NewNullLogger(),
	}
}

func (k *kafkaService) readTopics() []string {
	k.muHandlers.RLock()
	defer k.muHandlers.RUnlock()

	var topics []string

	for topic := range k.handlers {
		topics = append(topics, topic)
	}
	return topics
}

func (k *kafkaService) readHandler(topic string) (string, bool) {
	k.muHandlers.RLock()
	defer k.muHandlers.RUnlock()

	topic, ok := k.handlers[topic]
	return topic, ok
}

func (k *kafkaService) writeHandler(topic, handlerId string) {
	k.muHandlers.Lock()
	defer k.muHandlers.Unlock()

	if handlerId != "" {
		k.handlers[topic] = handlerId
	}
}

func (k *kafkaService) deleteHandler(topic string) {
	k.muHandlers.Lock()
	defer k.muHandlers.Unlock()

	delete(k.handlers, topic)
}

func (k *kafkaService) handleFx(topic string) logic.HandlerFx {
	return func(ctx context.Context, handlerId string, changes []*data.Change) error {
		if err := k.Client.Publish(topic, data.ToWrapper(
			&data.ChangeDigest{Changes: changes},
		)); err != nil {
			k.Error(logAlias, "error while publishing to topic \"%s\", handler \"%s\": %s", topic, handlerId, err)
			return err
		}
		k.Trace("published change(s) to topic \"%s\", handler \"%s\"", topic, handlerId)
		return nil
	}
}

func (k *kafkaService) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case internal_logger.Logger:
			k.Logger = p
		}
	}
}

func (k *kafkaService) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			k.logic = p
		case internal_kafka.Client:
			k.Client = p
		}
	}
	switch {
	case k.logic == nil:
		panic("logic not set")
	case k.Client == nil:
		panic("kafka client not set")
	}
}

func (k *kafkaService) Configure(items ...interface{}) error {
	k.Lock()
	defer k.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		default:
			c = new(Configuration)
			if err := internal_config.Get(item, configKey, c); err != nil {
				return err
			}
		case internal_config.Envs:
			c = new(Configuration)
			c.Default()
			c.FromEnv(v)
		case *Configuration:
			c = v
		case Configuration:
			c = &v
		}
	}
	if c == nil {
		return errors.New(internal_config.ErrConfigurationNotFound)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	k.config = c
	k.configured = true
	return nil
}

func (k *kafkaService) Initialize() error {
	k.Lock()
	defer k.Unlock()
	if k.initialized {
		return errors.New("already initialized")
	}
	if !k.configured {
		return errors.New("not configured")
	}
	k.ctx, k.cancel = context.WithCancel(context.Background())
	topic := k.config.Topic
	handlerId, err := k.logic.HandlerCreate(k.ctx, k.handleFx(topic))
	if err != nil {
		return err
	}
	k.writeHandler(topic, handlerId)
	k.Info("created handler \"%s\" for topic \"%s\"", handlerId, topic)
	k.initialized = true
	k.Info(logAlias + "initialized")
	return nil
}

func (k *kafkaService) Shutdown() {
	k.Lock()
	defer k.Unlock()
	if !k.initialized {
		return
	}
	k.cancel()
	for _, topic := range k.readTopics() {
		handlerId, ok := k.readHandler(topic)
		if !ok {
			continue
		}
		k.Unsubscribe(topic, handlerId)
		k.deleteHandler(topic)
	}
	k.Wait()
	k.initialized = false
	k.Info(logAlias + "shutdown")
}
