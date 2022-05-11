package file

import (
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

const (
	FileEmpty     string = "file is empty"
	LockFileEmpty string = "lock file is empty"
)

const (
	EnvNameFile        string = "BLUDGEON_META_FILE"
	EnvNameLockFile    string = "BLUDGEON_META_LOCK_FILE"
	EnvNameFileLocking string = "BLUDGEON_META_FILE_LOCKING"
)

const (
	DefaultFile        string = "data/bludgeon.json"
	DefaultLockFile    string = "data/bludgeon.lock"
	DefaultFileLocking bool   = false
)

var (
	ErrFileEmpty     = errors.New(FileEmpty)
	ErrLockFileEmpty = errors.New(LockFileEmpty)
)

type Configuration struct {
	File        string `json:"file"`
	FileLocking bool   `json:"file_locking"`
	LockFile    string `json:"lock_file"`
}

func (c *Configuration) Default(pwd string) {
	c.File = filepath.Join(pwd, DefaultFile)
	c.FileLocking = DefaultFileLocking
}

func (c *Configuration) Validate() (err error) {
	if c.File == "" {
		return ErrFileEmpty
	}
	if c.FileLocking {
		if c.LockFile == "" {
			return ErrLockFileEmpty
		}
	}
	return
}

func (c *Configuration) FromEnv(pwd string, envs map[string]string) {
	if file, ok := envs[EnvNameFile]; ok && file != "" {
		c.File = file
	}
	if lockFile, ok := envs[EnvNameLockFile]; ok && lockFile != "" {
		c.LockFile = lockFile
	}
	if s, ok := envs[EnvNameFileLocking]; ok && s != "" {
		if fileLocking, err := strconv.ParseBool(s); err == nil {
			c.FileLocking = fileLocking
		}
	}
}
