package cache

import "errors"

const UnsupportedDataType string = "unsupported data type"

var ErrUnsupportedDataType = errors.New(UnsupportedDataType)

type Cache interface {
	Write(key string, item interface{}) error
	Read(key string, item interface{}) error
	Delete(key string) error
}
