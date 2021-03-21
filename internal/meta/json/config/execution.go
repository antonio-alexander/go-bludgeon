package bludgeonmetajsonconfig

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func Default(directory string) Configuration {
	return Configuration{
		File: filepath.Join(directory, DefaultFile),
	}
}

func Validate(config *Configuration) (err error) {
	//check if the file is empty
	if config.File == "" {
		err = errors.New(ErrFileEmpty)
	}

	return
}

func FromEnv(pwd string, envs map[string]string, config *Configuration) (err error) {
	//get the file from the file environmental variable
	if file, ok := envs[EnvNameFile]; ok {
		config.File = file
	}

	return
}
