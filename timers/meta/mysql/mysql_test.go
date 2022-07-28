package mysql_test

import (
	"context"
	"os"
	"strings"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

	"github.com/stretchr/testify/assert"
)

var config *mysql.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config = new(mysql.Configuration)
	config.Default()
	config.FromEnv(pwd, envs)
}

func TestMetaMySQL(t *testing.T) {
	ctx := context.TODO()
	m := meta.New(
		logger.New(),
	)
	err := m.Initialize(config)
	assert.Nil(t, err)
	t.Run("Timer CRUD", tests.TestTimerCRUD(ctx, m))
	t.Run("Timers Read", tests.TestTimersRead(ctx, m))
	t.Run("Timer Logic", tests.TestTimerLogic(ctx, m))
}
