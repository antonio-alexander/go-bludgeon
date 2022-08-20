package file_test

import (
	"os"
	"path"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/changes/meta/tests"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

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
	//create meta
	logger := logger.New()
	m := meta.New()

	//set parameters
	m.SetParameters(logger)

	// initialize
	err := m.Configure(config)
	assert.Nil(t, err)
	err = m.Initialize()
	assert.Nil(t, err)
	defer func() {
		m.Shutdown()
	}()

	//test meta
	t.Run("Change CRUD", tests.TestChangeCRUD(m))
	t.Run("Changes Read", tests.TestChangeSearch(m))
	t.Run("Registration CRUD", tests.TestRegistrationCRUD(m))
	t.Run("Change Registrations", tests.TestRegistrationChanges(m))
}
