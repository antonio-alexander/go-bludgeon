package cache

import (
	"time"

	stash "github.com/antonio-alexander/go-stash"
	stashmemory "github.com/antonio-alexander/go-stash/memory"
)

const (
	DefaultMaxSize        int                  = 50 * 1024 * 1024 //50MB
	DefaultTimeToLive     time.Duration        = time.Hour
	DefaultDebug          bool                 = true
	DefaultEvictionPolicy stash.EvictionPolicy = stash.LeastFrequentlyUsed
)

type Configuration struct {
	*stashmemory.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: new(stashmemory.Configuration),
	}
}

func (c *Configuration) Default() {
	c.EvictionPolicy = DefaultEvictionPolicy
	c.TimeToLive = DefaultTimeToLive
	c.MaxSize = DefaultMaxSize
	c.Debug = DefaultDebug
}

func (c *Configuration) FromEnvs(envs map[string]string) {
	//
}
