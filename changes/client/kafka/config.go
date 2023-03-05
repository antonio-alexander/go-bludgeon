package kafkaclient

import "errors"

var DefaultTopic = "changes"

type Configuration struct {
	Topic string
}

func (c *Configuration) Default() {
	c.Topic = DefaultTopic
}

func (c *Configuration) Validate() error {
	if c.Topic == "" {
		return errors.New("topic is empty")
	}
	return nil
}
