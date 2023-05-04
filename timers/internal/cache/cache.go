package cache

import (
	"errors"
	"sync"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/timers/data"

	stash "github.com/antonio-alexander/go-stash"
	stashmemory "github.com/antonio-alexander/go-stash/memory"
)

type cache struct {
	sync.RWMutex
	logger.Logger
	stash.Stasher
	stashmemory.Memory
	config      *Configuration
	initialized bool
	configured  bool
}

func New() interface {
	Cache
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	stash := stashmemory.New()
	return &cache{
		Stasher: stash,
		Memory:  stash,
		Logger:  logger.NewNullLogger(),
	}
}

func (c *cache) Printf(format string, a ...any) (int, error) {
	c.Trace(format, a...)
	return 0, nil
}

func (c *cache) Configure(items ...interface{}) error {
	c.Lock()
	defer c.Unlock()

	var config *Configuration

	for _, item := range items {
		switch item := item.(type) {
		case Configuration:
			config = &item
		case *Configuration:
			config = item
		}
	}
	if config != nil {
		if err := c.Memory.Configure(config.Configuration); err != nil {
			return err
		}
		c.config = config
		c.configured = true
	}
	return nil
}

func (c *cache) SetParameters(parameters ...interface{}) {}

func (c *cache) SetUtilities(parameters ...interface{}) {
	c.Lock()
	defer c.Unlock()

	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			c.Logger = p
		}
	}
}

func (c *cache) Initialize() error {
	c.Lock()
	defer c.Unlock()

	if !c.configured {
		return errors.New("not configured")
	}
	if c.initialized {
		return errors.New("already initialized")
	}
	if err := c.Memory.Initialize(); err != nil {
		return err
	}
	c.initialized = true
	return nil
}

func (c *cache) Shutdown() {
	c.Lock()
	defer c.Unlock()

	if err := c.Memory.Shutdown(); err != nil {
		c.Error("error while shutting down stash: %s", err)
	}
	c.initialized, c.configured = false, false
}

func (c *cache) Write(key string, item interface{}) error {
	var value stash.Cacheable

	switch item := item.(type) {
	default:
		return errors.New("unsupported data type")
	case *data.Timer:
		value = item
	case *data.TimeSlice:
		value = item
	}
	_, err := c.Stasher.Write(key, value)
	return err
}

func (c *cache) Read(key string, item interface{}) error {
	var value stash.Cacheable

	switch item := item.(type) {
	default:
		return ErrUnsupportedDataType
	case *data.Timer:
		value = item
	case *data.TimeSlice:
		value = item
	}
	return c.Stasher.Read(key, value)
}

func (c *cache) Delete(key string) error {
	return c.Stasher.Delete(key)
}
