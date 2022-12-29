package mysql_test

import (
	"os"
	"strings"
	"testing"

	meta "github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/changes/meta/tests"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
)

var config *internal_mysql.Configuration

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	config = new(internal_mysql.Configuration)
	config.Default()
	config.FromEnv(envs)
	config.ParseTime = false
}

func TestMetaMysql(t *testing.T) {
	//create parameters
	logger := logger.New()
	m := meta.New()

	//set parameters
	m.SetParameters(logger)

	//configure
	err := m.Configure(nil, nil, config)
	assert.Nil(t, err)

	//initialize
	err = m.Initialize()
	assert.Nil(t, err)

	//execute tests
	t.Run("Change CRUD", tests.TestChangeCRUD(m))
	t.Run("Changes Search", tests.TestChangeSearch(m))
	t.Run("Registration CRUD", tests.TestRegistrationCRUD(m))
	t.Run("Change Registrations", tests.TestRegistrationChanges(m))
}
