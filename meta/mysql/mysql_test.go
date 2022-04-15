package metamysql_test

import (
	"testing"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"

	"github.com/stretchr/testify/assert"
)

const (
	TestDatabaseName string = "bludgeon"
	rootUsername     string = "bludgeon"
	bludgeonUsername string = "bludgeon"
	testCaseMap      string = "Test case: %s"
)

var config *metamysql.Configuration

func init() {
	//TODO: setup variables from environment?
	config = &metamysql.Configuration{
		Hostname:       "127.0.0.1", //metamysql.DefaultHostname,
		Port:           "3306",      //metamysql.DefaultPort,
		Username:       "bludgeon",  //metamysql.DefaultUsername,
		Password:       "bludgeon",  //metamysql.DefaultPassword,
		Database:       "bludgeon",  //TestDatabaseName,
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ParseTime:      true,
	}
}

func TestMetaMysql(t *testing.T) {
	m := metamysql.New(
		logger.New(),
	)
	err := m.Initialize(config)
	if assert.Nil(t, err) {
		t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
		t.Run("Timer CRUD", tests.TestTimerCRUD(m))
		// t.Run("Timers Read", tests.TestTimersRead(m))
		t.Run("Timer Logic", tests.TestTimerLogic(m))
	}
}
