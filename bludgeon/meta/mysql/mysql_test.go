package bludgeonmetamysql_test

//--------------------------------------------------------------------------------------------
// database_test.go contains all the tests to verify functionality of the bludgeon-database
// library, it contains all the unit and functions tests specific to the database
//--------------------------------------------------------------------------------------------

import (
	"fmt"
	"testing"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"

	"github.com/stretchr/testify/assert"
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
	validConfig   *mysql.Configuration
	defaultConfig *mysql.Configuration
)

func init() {
	//TODO: setup variables from environment?
	defaultConfig = &mysql.Configuration{}
	defaultConfig.Default()
	validConfig = &mysql.Configuration{
		Hostname:       mysql.DefaultHostname,
		Port:           mysql.DefaultPort,
		Username:       mysql.DefaultUsername,
		Password:       mysql.DefaultPassword,
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

	db := mysql.NewMetaMySQL()
	err := db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	uuid, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	timerWrite := bludgeon.Timer{
		UUID: uuid,
		// ActiveSliceUUID: "",
		Comment: "This is a test comment",
		Start:   tNow.UnixNano(),
		Finish:  tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Completed: true,
		Archived:  false,
		Billed:    true,
		// EmployeeID:  0,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := db.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	uuid, err := bludgeon.GenerateID()
	tNow := time.Now()
	assert.Nil(t, err)
	timerWrite := bludgeon.Timer{
		UUID: uuid,
		// ActiveSliceUUID: "",
		Comment: "This is a test comment",
		Start:   tNow.UnixNano(),
		Finish:  tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Completed: true,
		Archived:  false,
		Billed:    true,
		// EmployeeID:  0,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := db.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = db.TimerDelete(uuid)
	assert.Nil(t, err)
	_, err = db.TimerRead(uuid)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(mysql.ErrTimerNotFoundf, uuid), err.Error())
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	uuid, err := bludgeon.GenerateID()
	tNow := time.Now()
	assert.Nil(t, err)
	sliceWrite := bludgeon.TimeSlice{
		UUID: uuid,
		// TimerUUID:   uuid,
		Start:  tNow.UnixNano(),
		Finish: tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimeSliceWrite(uuid, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := db.TimeSliceRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	uuid, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	sliceWrite := bludgeon.TimeSlice{
		UUID: uuid,
		// TimerUUID:   uuid,
		Start:  tNow.UnixNano(),
		Finish: tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: false,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimeSliceWrite(uuid, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := db.TimeSliceRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = db.TimeSliceDelete(uuid)
	assert.Nil(t, err)
	_, err = db.TimerRead(uuid)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(mysql.ErrTimerNotFoundf, uuid), err.Error())
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	timerWrite := bludgeon.Timer{
		UUID: timerUUID,
		// ActiveSliceUUID: "",
		Comment: "This is a test comment",
		Start:   tNow.UnixNano(),
		Finish:  tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Completed: true,
		Archived:  false,
		Billed:    true,
		// EmployeeID:  0,
	}
	sliceWrite := bludgeon.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = db.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := db.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = db.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := bludgeon.GenerateID()
	assert.Nil(t, err)
	timerWrite := bludgeon.Timer{
		UUID: timerUUID,
		// ActiveSliceUUID: "",
		Comment: "This is a test comment",
		Start:   tNow.UnixNano(),
		Finish:  tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Completed: true,
		Archived:  false,
		Billed:    true,
		// EmployeeID:  0,
	}
	sliceWrite := bludgeon.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}
	db := mysql.NewMetaMySQL()
	err = db.Initialize(validConfig)
	assert.Nil(t, err)
	err = db.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = db.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	timerWrite.ActiveSliceUUID = sliceUUID
	err = db.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	timerRead, err := db.TimerRead(timerUUID)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = db.Shutdown()
	assert.Nil(t, err)
}
