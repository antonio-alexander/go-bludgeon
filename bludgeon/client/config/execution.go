package bludgeonclientconfig

func Default() Configuration {
	return Configuration{
		//TODO: populate
	}
}

func Validate(config *Configuration) (err error) {
	return
}

func FromEnv(pwd string, envs map[string]string, config *Configuration) (err error) {
	return
}
