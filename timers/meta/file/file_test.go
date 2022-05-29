package file_test

import (
	"os"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	metafile "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"

	"github.com/stretchr/testify/assert"
)

var (
	validConfig   *internal_file.Configuration
	defaultConfig *internal_file.Configuration
)

func init() {
	pwd, _ := os.Getwd()
	defaultConfig = &internal_file.Configuration{}
	defaultConfig.Default(pwd)
	validConfig = &internal_file.Configuration{
		File:        internal_file.DefaultFile,
		FileLocking: internal_file.DefaultFileLocking,
		LockFile:    internal_file.DefaultLockFile,
	}
}

func TestMetaFile(t *testing.T) {
	m := metafile.New(
		logger.New(),
	)
	err := m.Initialize(validConfig)
	assert.Nil(t, err)
	t.Run("Timer CRUD", tests.TestTimerCRUD(m))
	t.Run("Timers Read", tests.TestTimersRead(m))
	t.Run("Timer Logic", tests.TestTimerLogic(m))
	m.Shutdown()
}
