package mysql_test

import (
	"context"
	"os"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

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

func TestMetaMySql(t *testing.T) {
	ctx, logger := context.TODO(), internal_logger.New()
	m := meta.New()

	logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "test_meta_mysql",
	})
	err := m.Configure(config)
	assert.Nil(t, err)
	err = m.Initialize()
	assert.Nil(t, err)
	defer m.Shutdown()

	t.Run("Timer CRUD", tests.TestTimerCRUD(ctx, m))
	t.Run("Timers Read", tests.TestTimersRead(ctx, m))
	t.Run("Timer Logic", tests.TestTimerLogic(ctx, m))
}
