package mysql_test

import (
	"os"
	"strings"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

	"github.com/stretchr/testify/assert"
)

var config *meta.Configuration

func init() {
	pwd, _ := os.Getwd()
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config = new(meta.Configuration)
	config.Default()
	config.FromEnv(pwd, envs)
}

func TestMetaMySQL(t *testing.T) {
	m := meta.New(
		logger.New(),
	)
	err := m.Initialize(config)
	assert.Nil(t, err)
	t.Run("Timer CRUD", tests.TestTimerCRUD(m))
	t.Run("Timers Read", tests.TestTimersRead(m))
	t.Run("Timer Logic", tests.TestTimerLogic(m))
}
