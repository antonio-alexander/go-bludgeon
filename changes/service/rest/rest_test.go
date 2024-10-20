package rest_test

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	"github.com/antonio-alexander/go-bludgeon/changes/service/rest"
	"github.com/antonio-alexander/go-bludgeon/changes/service/rest/tests"
	"github.com/antonio-alexander/go-bludgeon/common"

	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/pkg/meta"
	internal_server "github.com/antonio-alexander/go-bludgeon/pkg/rest/server"

	"github.com/stretchr/testify/assert"
)

var (
	configService   = new(rest.Configuration)
	configServer    = new(internal_server.Configuration)
	configMetaMysql = new(mysql.Configuration)
	configMetaFile  = new(file.Configuration)
	configLogger    = new(internal_logger.Configuration)
)

type restServerTest struct {
	server interface {
		common.Initializer
		common.Configurer
	}
	service interface {
		common.Configurer
	}
	meta interface {
		common.Initializer
		common.Configurer
	}
	logic interface {
		common.Initializer
	}
	*tests.Fixture
}

func init() {
	rand.Seed(time.Now().UnixNano())
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configServer.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "9000"
	configServer.AllowedMethods = []string{http.MethodDelete, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodGet}
	configServer.ShutdownTimeout = 15 * time.Second
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configLogger.Default()
	configLogger.FromEnv(envs)
	configLogger.Level = internal_logger.Trace
	configLogger.Prefix = "test_logic"
	configService.Default()
	configService.FromEnv(envs)
}

func newRestServerTest(metaType internal_meta.Type) *restServerTest {
	var meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		common.Initializer
		common.Parameterizer
		common.Configurer
	}

	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	switch metaType {
	case internal_meta.TypeMemory:
		meta = memory.New()
		meta.SetParameters(logger)
	case internal_meta.TypeFile:
		meta = file.New()
		meta.SetParameters(logger)
	case internal_meta.TypeMySQL:
		meta = mysql.New()
		meta.SetParameters(logger)
	}
	changesLogic := logic.New()
	changesLogic.SetUtilities(logger)
	changesLogic.SetParameters(meta)
	changesService := rest.New()
	changesService.SetUtilities(logger)
	server := internal_server.New()
	server.SetUtilities(logger)
	changesService.SetParameters(changesLogic, server)
	server.SetParameters(changesService)
	return &restServerTest{
		server:  server,
		service: changesService,
		meta:    meta,
		logic:   changesLogic,
		Fixture: tests.NewFixture(configServer.Address, configServer.Port),
	}
}

func (r *restServerTest) Initialize(t *testing.T, metaType internal_meta.Type) {
	switch metaType {
	case internal_meta.TypeFile:
		err := r.meta.Configure(configMetaFile)
		if !assert.Nil(t, err) {
			assert.FailNow(t, "unable to configure meta file")
		}
	case internal_meta.TypeMySQL:
		err := r.meta.Configure(configMetaMysql)
		if !assert.Nil(t, err) {
			assert.FailNow(t, "unable to configure meta mysql")
		}
	}
	err := r.server.Configure(configServer)
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to configure server")
	}
	err = r.service.Configure(configService)
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to configure service")
	}
	err = r.meta.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize meta")
	}
	err = r.logic.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize logic")
	}
	err = r.server.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize server")
	}
	//KIM: we have to sleep here because the start for the rest
	// server isn't synchronous
	time.Sleep(2 * time.Second)
}

func (r *restServerTest) Shutdown(t *testing.T) {
	r.server.Shutdown()
	r.logic.Shutdown()
	r.meta.Shutdown()
}

func testChangesRestService(t *testing.T, metaType internal_meta.Type) {
	r := newRestServerTest(metaType)
	r.Initialize(t, metaType)
	defer r.Shutdown(t)

	// t.Run("Change Operations", r.TestChangeOperations)
	// t.Run("Change Streaming", r.TestChangeStreaming)
	t.Run("Change Registration", r.TestChangeRegistration)
}

func TestLogicMemory(t *testing.T) {
	testChangesRestService(t, internal_meta.TypeMemory)
}

func TestLogicFile(t *testing.T) {
	testChangesRestService(t, internal_meta.TypeFile)
}

func TestLogicMysql(t *testing.T) {
	testChangesRestService(t, internal_meta.TypeMySQL)
}
