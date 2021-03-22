package server

type Configuration struct {
	TokenWait int64 //how long a token is valid (seconds)
}

func (c *Configuration) Default() Configuration {
	return Configuration{
		//TODO: populate
	}
}

func (c *Configuration) Validate(config *Configuration) (err error) {
	return
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string, config *Configuration) (err error) {
	return
}
