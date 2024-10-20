package cache

import (
	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_cache "github.com/antonio-alexander/go-bludgeon/pkg/cache"
)

type cache struct {
	internal_cache.Cache
	common.Configurer
	common.Parameterizer
	common.Initializer
}

func New() interface {
	Cache
	common.Configurer
	common.Initializer
	common.Parameterizer
} {
	c := internal_cache.New()
	return &cache{
		Cache:         c,
		Configurer:    c,
		Parameterizer: c,
	}
}

func (c *cache) Write(change *data.Change) {
	if change == nil {
		return
	}
	c.Cache.Write(change.Id, change)
}

func (c *cache) Read(changeId string) *data.Change {
	change := new(data.Change)
	if err := c.Cache.Read(changeId, change); err != nil {
		return nil
	}
	return change
}
