package restclient

import (
	"errors"
	"strconv"
	"time"

	internal_cache "github.com/antonio-alexander/go-bludgeon/changes/internal/cache"
	internal_rest_client "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	internal_websocket_client "github.com/antonio-alexander/go-bludgeon/internal/websocket/client"
)

const (
	EnvNameAddress         string = "BLUDGEON_CHANGES_REST_ADDRESS"
	EnvNamePort            string = "BLUDGEON_CHANGES_REST_PORT"
	EnvNameDisableQueue    string = "BLUDGEON_CHANGES_DISABLE_QUEUE"
	EnvNameDisableCache    string = "BLUDGEON_CHANGES_DISABLE_CACHE"
	EnvNameAutoAcknowledge string = "BLUDGEON_CHANGES_AUTO_ACKNOWLEDGE"
)

const (
	DefaultAddress         string        = "localhost"
	DefaultPort            string        = "8014"
	DefaultAutoAcknowledge bool          = true
	DefaultUpsertQueueRate time.Duration = time.Second
	DefaultDisableQueue    bool          = false
	DefaultDisableCache    bool          = false
	DefaultQueueSize       int           = 100
)

type Configuration struct {
	Rest            *internal_rest_client.Configuration
	Websocket       *internal_websocket_client.Configuration
	Cache           *internal_cache.Configuration
	UpsertQueueRate time.Duration
	QueueSize       int
	DisableQueue    bool
	DisableCache    bool
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Rest:            new(internal_rest_client.Configuration),
		Websocket:       new(internal_websocket_client.Configuration),
		Cache:           internal_cache.NewConfiguration(),
		UpsertQueueRate: DefaultUpsertQueueRate,
		QueueSize:       DefaultQueueSize,
		DisableQueue:    DefaultDisableQueue,
		DisableCache:    DefaultDisableCache,
	}
}

func (c *Configuration) Default() {
	c.UpsertQueueRate = DefaultUpsertQueueRate
	c.Rest.Address = DefaultAddress
	c.Rest.Port = DefaultPort
	c.Rest.Timeout = internal_rest_client.DefaultTimeout
	c.Websocket.ReadTimeout = internal_websocket_client.DefaultReadTimeout
	c.Websocket.WriteTimeout = internal_websocket_client.DefaultWriteTimeout
	c.DisableCache = DefaultDisableCache
	c.DisableQueue = DefaultDisableQueue
	c.QueueSize = DefaultQueueSize
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if address, ok := envs[EnvNameAddress]; ok {
		c.Rest.Address = address
	}
	if port, ok := envs[EnvNamePort]; ok {
		c.Rest.Port = port
	}
	if s, ok := envs[EnvNameDisableCache]; ok {
		c.DisableQueue, _ = strconv.ParseBool(s)
	}
	if s, ok := envs[EnvNameDisableQueue]; ok {
		c.DisableCache, _ = strconv.ParseBool(s)
	}
}

func (c *Configuration) Validate() error {
	if c.Rest.Address == "" {
		return errors.New(internal_rest_client.ErrAddressEmpty)
	}
	if c.Rest.Port == "" {
		return errors.New(internal_rest_client.ErrPortEmpty)
	}
	return nil
}
