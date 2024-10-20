package logic_test

import (
	"os"
	"strings"
	"testing"

	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/changes/logic/tests"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	"github.com/antonio-alexander/go-bludgeon/common"

	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/pkg/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/pkg/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/pkg/meta/mysql"

	"github.com/stretchr/testify/assert"
)

var (
	mysqlConfig = new(internal_mysql.Configuration)
	fileConfig  = new(internal_file.Configuration)
	logConfig   = new(logger.Configuration)
)

type logicTest struct {
	meta interface {
		common.Initializer
		common.Configurer
	}
	logic common.Initializer
	*tests.Fixture
}

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	mysqlConfig = new(internal_mysql.Configuration)
	mysqlConfig.Default()
	mysqlConfig.FromEnv(envs)
	fileConfig = new(internal_file.Configuration)
	fileConfig.Default()
	fileConfig.FromEnv(envs)
	logConfig = new(logger.Configuration)
	logConfig.Default()
	logConfig.FromEnv(envs)
	logConfig.Level = logger.Trace
	logConfig.Prefix = "test_logic"
}

func newLogicTest(metaType internal_meta.Type) *logicTest {
	var meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		common.Initializer
		common.Parameterizer
		common.Configurer
	}

	logger := logger.New()
	logger.Configure(logConfig)
	switch metaType {
	case internal_meta.TypeMemory:
		meta = memory.New()
		meta.SetParameters(logger)
	case internal_meta.TypeFile:
		meta = file.New()
		meta.Configure()
		meta.SetParameters(logger)
	case internal_meta.TypeMySQL:
		meta = mysql.New()
		meta.SetParameters(logger)
	}
	logic := logic.New()
	logic.SetParameters(logger, meta)
	return &logicTest{
		meta:    meta,
		logic:   logic,
		Fixture: tests.NewFixture(logic),
	}
}

func (l *logicTest) Initialize(t *testing.T, metaType internal_meta.Type) {
	switch metaType {
	case internal_meta.TypeFile:
		err := l.meta.Configure(fileConfig)
		if !assert.Nil(t, err) {
			assert.FailNow(t, "unable to configure meta file")
		}
	case internal_meta.TypeMySQL:
		err := l.meta.Configure(mysqlConfig)
		if !assert.Nil(t, err) {
			assert.FailNow(t, "unable to configure meta mysql")
		}
	}
	err := l.meta.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize meta")
	}
	err = l.logic.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize logic")
	}
}

func (l *logicTest) Shutdown(t *testing.T) {
	l.logic.Shutdown()
	l.meta.Shutdown()
}

func testLogic(t *testing.T, metaType internal_meta.Type) {
	l := newLogicTest(internal_meta.TypeMemory)
	l.Initialize(t, metaType)
	defer l.Shutdown(t)

	t.Run("Change Registration", l.TestChangeRegistration)
	t.Run("Change Handlers", l.TestChangeHandlers)
}

func TestLogicMemory(t *testing.T) {
	testLogic(t, internal_meta.TypeMemory)
}

func TestLogicFile(t *testing.T) {
	testLogic(t, internal_meta.TypeFile)
}

func TestLogicMysql(t *testing.T) {
	testLogic(t, internal_meta.TypeMySQL)
}
