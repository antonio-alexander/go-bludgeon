package cache

import (
	data "github.com/antonio-alexander/go-bludgeon/changes/data"
)

type Cache interface {
	//CacheWrite can be used to write a change
	// to an in-memory store
	Write(*data.Change)

	//CacheRead can be used to read a change from
	// an in-memory store
	Read(changeId string) *data.Change

	//CacheDelete can be used to remove a change from
	// an in-memory store
	Delete(changeId string)
}
