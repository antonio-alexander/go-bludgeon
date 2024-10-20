package cache

import "encoding"

type Cacheable interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type Cache interface {
	//CacheWrite can be used to write a change
	// to an in-memory store
	Write(id string, v Cacheable)

	//CacheRead can be used to read a change from
	// an in-memory store
	Read(id string, v Cacheable) error

	//CacheDelete can be used to remove a change from
	// an in-memory store
	Delete(id string)
}
