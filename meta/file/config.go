package file

import (
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	ErrFileEmpty string = "file is empty"
	EnvNameFile  string = "BLUDGEON_META_JSON_FILE"
	DefaultFile  string = "data/bludgeon.json"
)

type Configuration struct {
	File string
}

func (c *Configuration) Default(pwd string) {
	c.File = filepath.Join(pwd, DefaultFile)
}

func (c *Configuration) Validate() (err error) {
	if c.File == "" {
		return errors.New(ErrFileEmpty)
	}

	return
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if file, ok := envs[EnvNameFile]; ok {
		c.File = file
	}
}
