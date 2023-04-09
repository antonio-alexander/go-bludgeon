package cache

import (
	"errors"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
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
	stash := stashmemory.New()
	return &cache{
		stash:  stash,
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
		if err := c.stash.Configure(config.Configuration); err != nil {
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

func (c *cache) Write(change *data.Change) {
	c.Lock()
	defer c.Unlock()

	if !c.initialized {
		return
	}
	if change == nil {
		return
	}
	replaced, err := c.stash.Write(change.Id, change)
	if err != nil {
		c.Error("error while writing cached item: %s", err)
	}
	if replaced {
		c.Trace("cached item (%s) replaced", change.Id)
	}
}

func (c *cache) Delete(changeId string) {
	c.Lock()
	defer c.Unlock()

	if !c.initialized {
		return
	}
	if err := c.stash.Delete(changeId); err != nil {
		c.Error("error while deleting cached item (%s): %s", changeId, err)
	}
}

func (c *cache) Read(changeId string) *data.Change {
	c.RLock()
	defer c.RUnlock()

	if !c.initialized {
		return nil
	}
	change := new(data.Change)
	if err := c.stash.Read(changeId, change); err != nil {
		c.Error("error while reading cached item (%s): %s", changeId, err)
		return nil
	}
	return change
}
