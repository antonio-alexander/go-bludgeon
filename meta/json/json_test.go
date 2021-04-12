package metajson_test

//--------------------------------------------------------------------------------------------
// database_test.go contains all the tests to verify functionality of the bludgeon-database
// library, it contains all the unit and functions tests specific to the database
//--------------------------------------------------------------------------------------------

import (
	"os"
	"testing"

	metajson "github.com/antonio-alexander/go-bludgeon/meta/json"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"
)

//--------------------------------------------------------------------------------------------------
//
//
// Normal Use Cases:
//
// Edge Cases:
//
//--------------------------------------------------------------------------------------------------

const ()

var (
	validConfig   *metajson.Configuration
	defaultConfig *metajson.Configuration
)

func init() {
	//TODO: setup variables from environment?
	pwd, _ := os.Getwd()
	defaultConfig = &metajson.Configuration{}
	defaultConfig.Default(pwd)
	validConfig = &metajson.Configuration{
		File: metajson.DefaultFile,
	}
}

//--------------------------------------------------------------------------------------------------
// UNIT TESTS
// Purpose: Unit Tests can only check the input and output of exported functions. For cases, inputs
// can be prefixed with an 'i' and outputs with an 'o. Use a map that uses a string and an anonymous
// struct. The string is the case description and the struct is a collection of inputs and outputs
//
// Function Prefix: TestUnit
//--------------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------------
// FUNCTION TESTS
//
// Purpose: Function Tests check the use of multiple package functions that do not rely on an
// external source
// Function Prefix: TestFunc
//
// Progression:
// 1. Level 1
// 		a. Level 2
// 			(1) Level 3
//--------------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------------
// INTEGRATION TESTS
//
// Purpose: Integration tests check the use of multiple package functions that rely on one or more
// external source
// Function Prefix: TestInt
//
// Progression:
// 1. Level 1
// 		a. Level 2
// 			(1) Level 3
//--------------------------------------------------------------------------------------------------

func TestIntInitializeShutdown(t *testing.T) {
	//Test: this unit test is meant to test whether or not the connect function works and to validate
	// certain use cases for that connect function
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntInitializeShutdown(t, json, validConfig)
}

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntTimerReadWrite(t, json, validConfig)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntDelete(t, json, validConfig)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntSliceReadWrite(t, json, validConfig)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntSliceDelete(t, json, validConfig)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntSliceTimer(t, json, validConfig)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := metajson.NewMetaJSON()
	tests.TestIntTimerActiveSlice(t, json, validConfig)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
