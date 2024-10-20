package kafka

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/client"
	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/common"

	internal_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/google/uuid"
)

type kafkaClient struct {
	sync.RWMutex
	sync.WaitGroup
	internal_logger.Logger
	internal_kafka.Client
	handlers    map[string]client.HandlerFx
	config      *Configuration
	subscribeId string
	initialized bool
	configured  bool
}

func New() interface {
	client.Handler
	common.Initializer
	common.Parameterizer
	common.Configurer
} {
	kafka := internal_kafka.New()
	return &kafkaClient{
		Client:   kafka,
		handlers: make(map[string]client.HandlerFx),
		Logger:   internal_logger.NewNullLogger(),
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
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			internal_kafka.Client
		}:
			k.Client = p
		}
	}
	switch {
	case k.Client == nil:
		panic("kafka client not set")
	}
}

func (k *kafkaClient) SetUtilities(parameters ...interface{}) {
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

func (k *kafkaClient) Initialize() error {
	k.Lock()
	defer k.Unlock()

	subscribeId, err := k.Subscribe(k.config.Topic, k.subscribeFx)
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
	k.Unsubscribe(k.subscribeId)
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
