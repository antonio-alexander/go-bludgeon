package service

import "errors"

const DefaultTopic string = "changes"

const EnvNameChangesTopic string = "BLUDGEON_CHANGES_TOPIC"

type Configuration struct {
	Topic string
}

func (c *Configuration) Validate() error {
	if c.Topic == "" {
		return errors.New("topic is empty")
	}
	return nil
}

func (c *Configuration) Default() {
	c.Topic = DefaultTopic
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if topic := envs[EnvNameChangesTopic]; topic != "" {
		c.Topic = topic
	}
}
