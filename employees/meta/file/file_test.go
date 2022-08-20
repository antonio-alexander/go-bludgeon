package file_test

import (
	"os"
	"path"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"

	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_meta.json"

var config = new(internal_file.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.File = path.Join("../../tmp", filename)
	os.Remove(config.File)
}

func TestMetaFile(t *testing.T) {
	meta, logger := meta.New(), internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "employee_memory_test",
		Level:  internal_logger.Trace,
	})
	meta.SetUtilities(logger)
	err := meta.Configure(config)
	assert.Nil(t, err)
	err = meta.Initialize()
	assert.Nil(t, err)
	defer meta.Shutdown()

	t.Run("Employee CRUD", tests.TestEmployeeCRUD(meta))
}
