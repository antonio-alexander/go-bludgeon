package metajson

import (
	"path/filepath"

	"github.com/pkg/errors"
)

//error constants
const (
	ErrFileEmpty string = "File is empty"
)

//environmental variables
const (
	EnvNameFile string = "BLUDGEON_META_JSON_FILE"
)

//defaults
const (
	DefaultFile string = "data/bludgeon.json"
)

type Configuration struct {
	File string
}

func (c *Configuration) Default(pwd string) {
	c.File = filepath.Join(pwd, DefaultFile)
}

func (c *Configuration) Validate() (err error) {
	//check if the file is empty
	if c.File == "" {
		err = errors.New(ErrFileEmpty)

		return
	}

	return
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	//get the file from the file environmental variable
	if file, ok := envs[EnvNameFile]; ok {
		c.File = file
	}
}
