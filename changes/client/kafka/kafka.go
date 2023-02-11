package kafkaclient

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/client"
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/internal"

	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/google/uuid"
)

type kafkaClient struct {
	sync.RWMutex
	sync.WaitGroup
	internal_logger.Logger
	handlers    map[string]client.HandlerFx
	subscribeId string
	initialized bool
	configured  bool
	config      *Configuration
	kafkaClient interface {
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
		internal_kafka.Client
	}
}

func New() interface {
	client.Handler
	internal.Initializer
	internal.Parameterizer
	internal.Configurer
} {
	return &kafkaClient{
		handlers:    make(map[string]client.HandlerFx),
		kafkaClient: internal_kafka.New(),
		Logger:      logger.NewNullLogger(),
	}
}

func (k *kafkaClient) handleFx(changes ...*data.Change) {
	k.RLock()
	defer k.RUnlock()

	var wg sync.WaitGroup

	for _, handleFx := range k.handlers {
		wg.Add(1)
		go func(handleFx client.HandlerFx) {
			defer wg.Done()

			handleFx(changes...)
		}(handleFx)
	}
	wg.Wait()
}

func (k *kafkaClient) subscribeFx(topic string, bytes []byte) {
	if len(bytes) == 0 {
		k.Trace(logAlias + "no bytes received")
		return
	}
	wrapper := &data.Wrapper{}
	if err := json.Unmarshal(bytes, wrapper); err != nil {
		k.Error(logAlias+"error while unmarshalling json: %s", err)
		return
	}
	item, err := data.FromWrapper(wrapper)
	if err != nil {
		k.Error(logAlias+"error during  FromWrapper: %s", err)
		return
	}
	switch v := item.(type) {
	default:
		k.Trace(logAlias+"received unsupported type: %T", v)
	case *data.Change:
		k.handleFx(v)
	case *data.ChangeDigest:
		k.handleFx(v.Changes...)
	}
}

func (k *kafkaClient) SetParameters(parameters ...interface{}) {
	k.kafkaClient.SetParameters(parameters...)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			internal.Configurer
			internal.Initializer
			internal.Parameterizer
			internal_kafka.Client
		}:
			k.kafkaClient = p
		}
	}
	switch {
	case k.kafkaClient == nil:
		panic("kafka client not set")
	}
}

func (k *kafkaClient) SetUtilities(parameters ...interface{}) {
	k.kafkaClient.SetUtilities(parameters...)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case internal_logger.Logger:
			k.Logger = p
		}
	}
}

func (k *kafkaClient) Configure(items ...interface{}) error {
	k.Lock()
	defer k.Unlock()

	var c *Configuration

	if err := k.kafkaClient.Configure(items...); err != nil {
		return err
	}
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

func (k *kafkaClient) Initialize() error {
	k.Lock()
	defer k.Unlock()

	if err := k.kafkaClient.Initialize(); err != nil {
		return err
	}
	subscribeId, err := k.kafkaClient.Subscribe(k.config.Topic, k.subscribeFx)
	if err != nil {
		return err
	}
	k.initialized, k.subscribeId = true, subscribeId
	return nil
}

func (k *kafkaClient) Shutdown() {
	k.Lock()
	defer k.Unlock()

	if !k.initialized {
		return
	}
	k.kafkaClient.Unsubscribe(k.subscribeId)
	k.initialized, k.subscribeId = false, ""
	k.Info(logAlias + "shutdown")
}

func (k *kafkaClient) HandlerCreate(handlerFx client.HandlerFx) (string, error) {
	k.Lock()
	defer k.Unlock()

	handlerId := uuid.Must(uuid.NewRandom()).String()
	k.handlers[handlerId] = handlerFx
	return handlerId, nil
}

func (k *kafkaClient) HandlerConnected(handlerId string) (bool, error) {
	return k.subscribeId != "", nil
}

func (k *kafkaClient) HandlerDelete(handlerId string) error {
	k.Lock()
	defer k.Unlock()
	_, ok := k.handlers[handlerId]
	if !ok {
		return errors.New("handler not found")
	}
	delete(k.handlers, handlerId)
	return nil
}
