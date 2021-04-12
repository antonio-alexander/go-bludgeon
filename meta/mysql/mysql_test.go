package metamysql_test

//--------------------------------------------------------------------------------------------
// database_test.go contains all the tests to verify functionality of the bludgeon-database
// library, it contains all the unit and functions tests specific to the database
//--------------------------------------------------------------------------------------------

import (
	"testing"
	"time"

	metamysql "github.com/antonio-alexander/go-bludgeon/meta/mysql"
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

const (
	TestDatabaseName string = "bludgeon"
	rootUsername     string = "root"
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

	db := metamysql.NewMetaMySQL()
	tests.TestIntInitializeShutdown(t, db, validConfig)
}

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntTimerReadWrite(t, db, validConfig)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntDelete(t, db, validConfig)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntSliceReadWrite(t, db, validConfig)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntSliceDelete(t, db, validConfig)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntSliceTimer(t, db, validConfig)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	db := metamysql.NewMetaMySQL()
	tests.TestIntTimerActiveSlice(t, db, validConfig)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
//
