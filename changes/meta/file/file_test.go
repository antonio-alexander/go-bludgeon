package file_test

import (
	"os"
	"path"
	"strings"
	"testing"

	file "github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/changes/meta/tests"

	common "github.com/antonio-alexander/go-bludgeon/common"
	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_meta.json"

var config = new(file.Configuration)

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

type fileMetaTest struct {
	metaFile interface {
		common.Initializer
		common.Configurer
	}
	*tests.Fixture
}

func newFileMetaTest() *fileMetaTest {
	logger := logger.New()
	metaFile := file.New()
	metaFile.SetParameters()
	metaFile.SetUtilities(logger)
	return &fileMetaTest{
		metaFile: metaFile,
		Fixture:  tests.NewFixture(metaFile),
	}
}

func (m *fileMetaTest) Initialize(t *testing.T) {
	err := m.metaFile.Configure(config)
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to configure meta file")
	}
	err = m.metaFile.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize persistence")
	}
}

func (m *fileMetaTest) Shutdown(t *testing.T) {
	m.metaFile.Shutdown()
}

func testMetaFile(t *testing.T) {
	m := newFileMetaTest()
	m.Initialize(t)
	defer m.Shutdown(t)

	t.Run("Change CRUD", m.TestChangeCRUD)
	t.Run("Changes Read", m.TestChangeSearch)
	t.Run("Registration CRUD", m.TestRegistrationCRUD)
	t.Run("Change Registrations", m.TestRegistrationChanges)
}

func TestMetaFile(t *testing.T) {
	testMetaFile(t)
}
