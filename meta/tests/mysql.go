package tests

//--------------------------------------------------------------------------------------------
// database_test.go contains all the tests to verify functionality of the bludgeon-database
// library, it contains all the unit and functions tests specific to the database
//--------------------------------------------------------------------------------------------

import (
	"fmt"
	"testing"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"

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

func TestIntInitializeShutdown(t *testing.T, meta common.MetaOwner, validConfig interface{}) {
	//Test: this unit test is meant to test whether or not the connect function works and to validate
	// certain use cases for that connect function
	//Notes:
	//Verification:

	err := meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerReadWrite(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
}, validConfig interface{}) {

	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	uuid, err := common.GenerateID()
	assert.Nil(t, err)
	timerWrite := common.Timer{
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
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := meta.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntDelete(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
}, validConfig interface{}) {
	//Test:
	//Notes:
	//Verification:

	uuid, err := common.GenerateID()
	tNow := time.Now()
	assert.Nil(t, err)
	timerWrite := common.Timer{
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
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := meta.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = meta.TimerDelete(uuid)
	assert.Nil(t, err)
	_, err = meta.TimerRead(uuid)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(common.ErrTimerNotFoundf, uuid), err.Error())
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceReadWrite(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
	common.MetaTimeSlice
}, validConfig interface{}) {

	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	sliceUUID, err := common.GenerateID()
	assert.Nil(t, err)
	timerUUID, err := common.GenerateID()
	assert.Nil(t, err)
	timerWrite := common.Timer{
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
	assert.Nil(t, err)
	sliceWrite := common.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		// Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = meta.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := meta.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceDelete(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
	common.MetaTimeSlice
}, validConfig interface{}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := common.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := common.GenerateID()
	assert.Nil(t, err)
	timerWrite := common.Timer{
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
	sliceWrite := common.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		// Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: tNow.Add(5*time.Second).UnixNano() - tNow.UnixNano(),
		Archived: true,
	}
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = meta.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := meta.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = meta.TimeSliceDelete(sliceUUID)
	assert.Nil(t, err)
	_, err = meta.TimerRead(sliceUUID)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(common.ErrTimerNotFoundf, sliceUUID), err.Error())
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceTimer(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
	common.MetaTimeSlice
}, validConfig interface{}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := common.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := common.GenerateID()
	assert.Nil(t, err)
	timerWrite := common.Timer{
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
	sliceWrite := common.TimeSlice{
		UUID:        sliceUUID,
		TimerUUID:   timerUUID,
		Start:       tNow.UnixNano(),
		Finish:      tNow.Add(5 * time.Second).UnixNano(),
		ElapsedTime: tNow.Add(5*time.Second).UnixNano() - tNow.UnixNano(),
		Archived:    true,
	}
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = meta.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := meta.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = meta.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerActiveSlice(t *testing.T, meta interface {
	common.MetaOwner
	common.MetaTimer
	common.MetaTimeSlice
}, validConfig interface{}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := common.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := common.GenerateID()
	assert.Nil(t, err)
	timerWrite := common.Timer{
		UUID: timerUUID,
		// ActiveSliceUUID: "",
		Comment:     "This is a test comment",
		Start:       tNow.UnixNano(),
		Finish:      tNow.Add(10 * time.Second).UnixNano(),
		ElapsedTime: tNow.Add(5*time.Second).UnixNano() - tNow.UnixNano(),
		Completed:   true,
		Archived:    false,
		Billed:      true,
		// EmployeeID:  0,
	}
	sliceWrite := common.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}
	err = meta.Initialize(validConfig)
	assert.Nil(t, err)
	err = meta.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = meta.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	timerWrite.ActiveSliceUUID = sliceUUID
	err = meta.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	timerRead, err := meta.TimerRead(timerUUID)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = meta.Shutdown()
	assert.Nil(t, err)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
//
