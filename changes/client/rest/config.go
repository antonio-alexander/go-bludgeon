package rest

import (
	"errors"
	"strconv"
	"time"

	internal_cache "github.com/antonio-alexander/go-bludgeon/changes/internal/cache"
	internal_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/client"
	internal_websocket "github.com/antonio-alexander/go-bludgeon/pkg/websocket/client"
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
	DefaultUpsertQueueRate time.Duration = time.Second
	DefaultAutoAcknowledge bool          = true
	DefaultDisableQueue    bool          = false
	DefaultDisableCache    bool          = false
)

type Configuration struct {
	Rest            internal_rest.Configuration
	Websocket       internal_websocket.Configuration
	Cache           internal_cache.Configuration
	UpsertQueueRate time.Duration
	DisableQueue    bool
	DisableCache    bool
}

func (c *Configuration) Default() {
	c.UpsertQueueRate = DefaultUpsertQueueRate
	c.Rest.Address = DefaultAddress
	c.Rest.Port = DefaultPort
	c.Rest.Timeout = internal_rest.DefaultTimeout
	c.Websocket.ReadTimeout = internal_websocket.DefaultReadTimeout
	c.Websocket.WriteTimeout = internal_websocket.DefaultWriteTimeout
	c.DisableCache = DefaultDisableCache
	c.DisableQueue = DefaultDisableQueue
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
		return errors.New(internal_rest.ErrAddressEmpty)
	}
	if c.Rest.Port == "" {
		return errors.New(internal_rest.ErrPortEmpty)
	}
	return nil
}
