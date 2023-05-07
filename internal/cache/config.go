package cache

import (
	"time"

	stash "github.com/antonio-alexander/go-stash"
	stashmemory "github.com/antonio-alexander/go-stash/memory"
)

const (
	EnvCacheMaxSize        string = "BLUDGEON_CACHE_MAX_SIZE"
	EnvCacheTimeToLive     string = "BLUDGEON_CACHE_TIME_TO_LIVE"
	EnvCacheDebug          string = "BLUDGEON_CACHE_DEBUG"
	EnvCacheEvictionPolicy string = "BLUDGEON_CACHE_EVICTION_POLICY"
)

const (
	DefaultMaxSize        int                  = 50 * 1024 * 1024 //50MB
	DefaultTimeToLive     time.Duration        = time.Hour
	DefaultDebug          bool                 = true
	DefaultEvictionPolicy stash.EvictionPolicy = stash.LeastFrequentlyUsed
)

type Configuration struct {
	stashmemory.Configuration
}
