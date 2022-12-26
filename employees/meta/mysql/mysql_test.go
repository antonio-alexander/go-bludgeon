package mysql_test

import (
	"os"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
)

var config = new(internal_mysql.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config.Default()
	config.FromEnv(envs)
	config.ParseTime = true
}

func TestMetaMysql(t *testing.T) {
	meta, logger := meta.New(), internal_logger.New()

	//set parameters
	meta.SetUtilities(logger)

	//configure
	err := logger.Configure(internal_logger.Configuration{
		Prefix: "test_meta_mysql",
		Level:  internal_logger.Trace,
	})
	assert.Nil(t, err)
	err = meta.Configure(config)
	assert.Nil(t, err)

	//initialize
	err = meta.Initialize()
	assert.Nil(t, err)
	defer meta.Shutdown()

	//test
	t.Run("Employee CRUD", tests.TestEmployeeCRUD(meta))
}
