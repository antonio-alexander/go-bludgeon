package cache

import (
	"strconv"
	"time"

	internal_cache "github.com/antonio-alexander/go-bludgeon/pkg/cache"

	"github.com/antonio-alexander/go-stash"
)

type Configuration struct {
	internal_cache.Configuration
}

func (c *Configuration) Default() {
	c.EvictionPolicy = internal_cache.DefaultEvictionPolicy
	c.TimeToLive = internal_cache.DefaultTimeToLive
	c.MaxSize = internal_cache.DefaultMaxSize
	c.Debug = internal_cache.DefaultDebug
}

func (c *Configuration) FromEnvs(envs map[string]string) {
	if s := envs["BLUDGEON_CACHE_MAX_SIZE"]; s != "" {
		c.MaxSize, _ = strconv.Atoi(s)
	}
	if s := envs["BLUDGEON_CACHE_TIME_TO_LIVE"]; s != "" {
		i, _ := strconv.Atoi(s)
		c.TimeToLive = time.Duration(time.Duration(i) * time.Second)
	}
	if s := envs["BLUDGEON_CACHE_DEBUG"]; s != "" {
		c.Debug, _ = strconv.ParseBool(s)
	}
	if s := envs["BLUDGEON_CACHE_EVICTION_POLICY"]; s != "" {
		c.EvictionPolicy = stash.EvictionPolicy(s)
	}
}
