package file_test

import (
	"os"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	metafile "github.com/antonio-alexander/go-bludgeon/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"

	"github.com/stretchr/testify/assert"
)

var (
	validConfig   *metafile.Configuration
	defaultConfig *metafile.Configuration
)

func init() {
	//TODO: setup variables from environment?
	pwd, _ := os.Getwd()
	defaultConfig = &metafile.Configuration{}
	defaultConfig.Default(pwd)
	validConfig = &metafile.Configuration{
		File: metafile.DefaultFile,
	}
}

func TestMetaMemory(t *testing.T) {
	m := metafile.New(
		logger.New(),
	)
	err := m.Initialize(validConfig)
	if assert.Nil(t, err) {
		t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
		t.Run("Timer CRUD", tests.TestTimerCRUD(m))
		t.Run("Timers Read", tests.TestTimersRead(m))
		t.Run("Timer Logic", tests.TestTimerLogic(m))
	}
	m.Shutdown()
}
