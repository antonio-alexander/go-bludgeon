package file_test

import (
	"os"
	"testing"

	file "github.com/antonio-alexander/go-bludgeon/meta/file"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"

	"github.com/stretchr/testify/assert"
)

var (
	validConfig   *file.Configuration
	defaultConfig *file.Configuration
)

func init() {
	//TODO: setup variables from environment?
	pwd, _ := os.Getwd()
	defaultConfig = &file.Configuration{}
	defaultConfig.Default(pwd)
	validConfig = &file.Configuration{
		File: file.DefaultFile,
	}
}

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntTimerReadWrite(t, json)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntDelete(t, json)
	err := json.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntSliceReadWrite(t, json)
	err := json.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntSliceDelete(t, json)
	err := json.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntSliceTimer(t, json)
	err := json.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := file.New()
	json.Initialize(validConfig)
	tests.TestIntTimerActiveSlice(t, json)
	err := json.Shutdown()
	assert.Nil(t, err)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
