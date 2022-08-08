package mysql_test

import (
	"os"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

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

func TestMetaMysql(t *testing.T) {
	m := meta.New(logger.New())
	err := m.Initialize(config)
	assert.Nil(t, err)
	t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
}
