package metamysql_test

import (
	"testing"
	"time"

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

var (
	validConfig   *metamysql.Configuration
	defaultConfig *metamysql.Configuration
)

func init() {
	//TODO: setup variables from environment?
	defaultConfig = &metamysql.Configuration{}
	defaultConfig.Default()
	validConfig = &metamysql.Configuration{
		Hostname:       metamysql.DefaultHostname,
		Port:           metamysql.DefaultPort,
		Username:       metamysql.DefaultUsername,
		Password:       metamysql.DefaultPassword,
		Database:       TestDatabaseName,
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ParseTime:      false,
	}
}

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntTimerReadWrite(t, meta)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntDelete(t, meta)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntSliceReadWrite(t, meta)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntSliceDelete(t, meta)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntSliceTimer(t, meta)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	meta := metamysql.New()
	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	tests.TestIntTimerActiveSlice(t, meta)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
