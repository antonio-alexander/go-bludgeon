package service

const logAlias string = "[kafka_service] "

var DefaultTopics = []string{"changes"}

type Configuration struct {
	Topics []string
}

func (c *Configuration) Default() {
	c.Topics = DefaultTopics
}
