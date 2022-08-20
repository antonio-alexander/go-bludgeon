package kafkaclient

var DefaultTopic = "changes"

type Configuration struct {
	Topic string
}

func (c *Configuration) Default() {
	c.Topic = DefaultTopic
}
