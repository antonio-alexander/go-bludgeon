package bludgeonmetamysqlconfig

import (
	"fmt"
)

func Default() Configuration {
	return Configuration{
		//TODO: populate
	}
}

//Validate is used to ensure that the values being configured make sense
// it's not necessarily to prevent a misconfiguration, but to use default values in the
// event a value doesn't exist
func Validate(config *Configuration) (err error) {
	switch config.Driver {
	case "mysql", "postgres":
		switch config.Driver {
		case "mysql":
			if config.Port == "" {
				config.Port = DefaultMysqlPort
			}
		case "postgres":
			if config.Port == "" {
				config.Port = DefaultPostgresPort
			}
		}
		if config.Database == "" {
			config.Database = DefaultDatabase
		}
		if config.Username == "" {
			config.Username = DefaultUsername
		}
		if config.Password == "" {
			config.Password = DefaultPassword
		}
		if config.Hostname == "" {
			config.Hostname = DefaultHostname
		}
	case "sqlite":
		if config.FilePath == "" {
			config.FilePath = DefaultDatabasePath
		}
	default:
		err = fmt.Errorf(ErrDriverUnsupported, config.Driver)
	}

	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout
	}

	return
}

func FromEnv(pwd string, envs map[string]string, config *Configuration) (err error) {
	return
}
