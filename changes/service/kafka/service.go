package service

import (
	"context"
	"errors"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/internal"

	"github.com/antonio-alexander/go-bludgeon/internal/kafka"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
)

type kafkaService struct {
	sync.RWMutex
	sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	logger.Logger
	kafka.Client
	logic       logic.Logic
	muHandlers  sync.RWMutex
	handlers    map[string]string
	initialized bool
	configured  bool
	config      *Configuration
}

func New() interface {
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	return &kafkaService{
		handlers: make(map[string]string),
		Logger:   logger.NewNullLogger(),
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
		return nil
	}
}

func (k *kafkaService) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			k.Logger = p
		}
	}
}

func (k *kafkaService) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			k.logic = p
		case kafka.Client:
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
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		c = new(Configuration)
		c.Default()
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
	for _, topic := range k.config.Topics {
		handlerId, err := k.logic.HandlerCreate(k.ctx, k.handleFx(topic))
		if err != nil {
			return err
		}
		k.writeHandler(topic, handlerId)
	}
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
