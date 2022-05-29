package file_test

import (
	"os"
	"strings"
	"testing"

	metafile "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"

	"github.com/stretchr/testify/assert"
)

var config *internal_file.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config = &internal_file.Configuration{}
	config.Default(pwd)
	config.FromEnv(pwd, envs)
}

func TestMetaFile(t *testing.T) {
	m := metafile.New(
		logger.New(),
	)
	err := m.Initialize(config)
	assert.Nil(t, err)
	t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
	m.Shutdown()
}
