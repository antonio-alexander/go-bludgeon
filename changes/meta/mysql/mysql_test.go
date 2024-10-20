package mysql_test

import (
	"os"
	"strings"
	"testing"

	mysql "github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/changes/meta/tests"
	common "github.com/antonio-alexander/go-bludgeon/common"

	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/stretchr/testify/assert"
)

var config = new(mysql.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.ParseTime = false
}

type fileMetaTest struct {
	metaMysql interface {
		common.Initializer
		common.Configurer
	}
	*tests.Fixture
}

func newMysqlMetaTest() *fileMetaTest {
	logger := logger.New()
	metaMysql := mysql.New()
	metaMysql.SetUtilities(logger)
	return &fileMetaTest{
		metaMysql: metaMysql,
		Fixture:   tests.NewFixture(metaMysql),
	}
}

func (m *fileMetaTest) Initialize(t *testing.T) {
	err := m.metaMysql.Configure(config)
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to configure meta file")
	}
	err = m.metaMysql.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize persistence")
	}
}

func (m *fileMetaTest) Shutdown(t *testing.T) {
	m.metaMysql.Shutdown()
}

func testMetaMysql(t *testing.T) {
	m := newMysqlMetaTest()
	m.Initialize(t)
	defer m.Shutdown(t)

	t.Run("Change CRUD", m.TestChangeCRUD)
	t.Run("Changes Read", m.TestChangeSearch)
	t.Run("Registration CRUD", m.TestRegistrationCRUD)
	t.Run("Change Registrations", m.TestRegistrationChanges)
}

func TestMetaMysql(t *testing.T) {
	testMetaMysql(t)
}
