package kafka

import (
	"regexp"

	"github.com/pkg/errors"
)

const (
	logAlias  string = "[kafka_client] "
	ConfigNil string = "config is nil"
)

var ErrConfigNil = errors.New(ConfigNil)

type HandleFx func(topic string, bytes []byte)

type Client interface {
	Publish(topic string, item interface{}) (err error)
	Subscribe(topic string, handler HandleFx) (handlerId string, err error)
	Unsubscribe(topic string, handlerIds ...string)
	Topics(regEx *regexp.Regexp) ([]string, error)
}
