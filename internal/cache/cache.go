package cache

import (
	"errors"
	"sync"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	stash "github.com/antonio-alexander/go-stash"
	stashmemory "github.com/antonio-alexander/go-stash/memory"
)

type cache struct {
	sync.RWMutex
	logger.Logger
	stash interface {
		stash.Stasher
		stashmemory.Memory
	}
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
	return &cache{
		stash:  stashmemory.New(),
		Logger: logger.NewNullLogger(),
	}
}

func (c *cache) Printf(format string, a ...any) (int, error) {
	c.Trace(format, a...)
	return 0, nil
}

func (c *cache) Configure(items ...interface{}) error {
	c.Lock()
	defer c.Unlock()

	var cc *Configuration

	for _, item := range items {
		switch item := item.(type) {
		case Configuration:
			cc = &item
		case *Configuration:
			cc = item
		}
	}
	if cc != nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	if err := c.stash.Configure(cc.Configuration); err != nil {
		return err
	}
	c.config = cc
	c.configured = true
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
	c.stash.SetParameters(c)
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
	if err := c.stash.Initialize(); err != nil {
		return err
	}
	c.initialized = true
	return nil
}

func (c *cache) Shutdown() {
	c.Lock()
	defer c.Unlock()

	if !c.initialized {
		return
	}
	if err := c.stash.Shutdown(); err != nil {
		c.Error("error while shutting down stash: %s", err)
	}
	c.initialized, c.configured = false, false
}

func (c *cache) Write(id string, v Cacheable) {
	if !c.initialized {
		return
	}
	if v == nil {
		return
	}
	replaced, err := c.stash.Write(id, v)
	if err != nil {
		c.Error("error while writing cached item: %s", err)
	}
	if replaced {
		c.Trace("cached item (%s) replaced", id)
	}
}

func (c *cache) Delete(id string) {
	if !c.initialized {
		return
	}
	if err := c.stash.Delete(id); err != nil {
		c.Error("error while deleting cached item (%s): %s", id, err)
	}
}

func (c *cache) Read(id string, v Cacheable) error {
	if !c.initialized {
		return nil
	}
	if err := c.stash.Read(id, v); err != nil {
		c.Error("error while reading cached item (%s): %s", id, err)
		return err
	}
	return nil
}
