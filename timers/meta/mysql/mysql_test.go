package metamysql_test

import (
	"os"
	"strings"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	metamysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
)

var config *internal_mysql.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config = new(internal_mysql.Configuration)
	config.Default()
	config.FromEnv(pwd, envs)
}

func TestMetaMySQL(t *testing.T) {
	m := metamysql.New(
		logger.New(),
	)
	err := m.Initialize(config)
	assert.Nil(t, err)
	t.Run("Timer CRUD", tests.TestTimerCRUD(m))
	t.Run("Timers Read", tests.TestTimersRead(m))
	t.Run("Timer Logic", tests.TestTimerLogic(m))
}
