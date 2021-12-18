package tests

import (
	"fmt"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"

	"github.com/stretchr/testify/assert"
)

func TestIntTimerReadWrite(t *testing.T, m interface {
	meta.Owner
	meta.Timer
}) {

	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	uuid, err := data.GenerateID()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	err = m.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := m.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = m.Shutdown()
	assert.Nil(t, err)
}

func TestIntDelete(t *testing.T, m interface {
	meta.Owner
	meta.Timer
}) {
	//Test:
	//Notes:
	//Verification:

	uuid, err := data.GenerateID()
	tNow := time.Now()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	err = m.TimerWrite(uuid, timerWrite)
	assert.Nil(t, err)
	timerRead, err := m.TimerRead(uuid)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = m.TimerDelete(uuid)
	assert.Nil(t, err)
	_, err = m.TimerRead(uuid)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(meta.ErrTimerNotFoundf, uuid), err.Error())
	err = m.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceReadWrite(t *testing.T, m interface {
	meta.Owner
	meta.Timer
	meta.TimeSlice
}) {

	//Test:
	//Notes:
	//Verification:
	tNow := time.Now()
	sliceUUID, err := data.GenerateID()
	assert.Nil(t, err)
	timerUUID, err := data.GenerateID()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	sliceWrite := data.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		// Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}

	assert.Nil(t, err)
	err = m.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = m.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := m.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = m.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceDelete(t *testing.T, m interface {
	meta.Owner
	meta.Timer
	meta.TimeSlice
}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := data.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := data.GenerateID()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	sliceWrite := data.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		// Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: tNow.Add(5*time.Second).UnixNano() - tNow.UnixNano(),
		Archived: true,
	}

	assert.Nil(t, err)
	err = m.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = m.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := m.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = m.TimeSliceDelete(sliceUUID)
	assert.Nil(t, err)
	_, err = m.TimerRead(sliceUUID)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(meta.ErrTimerNotFoundf, sliceUUID), err.Error())
	err = m.Shutdown()
	assert.Nil(t, err)
}

func TestIntSliceTimer(t *testing.T, m interface {
	meta.Owner
	meta.Timer
	meta.TimeSlice
}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := data.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := data.GenerateID()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	sliceWrite := data.TimeSlice{
		UUID:        sliceUUID,
		TimerUUID:   timerUUID,
		Start:       tNow.UnixNano(),
		Finish:      tNow.Add(5 * time.Second).UnixNano(),
		ElapsedTime: tNow.Add(5*time.Second).UnixNano() - tNow.UnixNano(),
		Archived:    true,
	}

	assert.Nil(t, err)
	err = m.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = m.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	sliceRead, err := m.TimeSliceRead(sliceUUID)
	assert.Nil(t, err)
	assert.Equal(t, sliceWrite, sliceRead)
	err = m.Shutdown()
	assert.Nil(t, err)
}

func TestIntTimerActiveSlice(t *testing.T, m interface {
	meta.Owner
	meta.Timer
	meta.TimeSlice
}) {
	//Test:
	//Notes:
	//Verification:

	tNow := time.Now()
	timerUUID, err := data.GenerateID()
	assert.Nil(t, err)
	sliceUUID, err := data.GenerateID()
	assert.Nil(t, err)
	timerWrite := data.Timer{
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
	sliceWrite := data.TimeSlice{
		UUID:      sliceUUID,
		TimerUUID: timerUUID,
		Start:     tNow.UnixNano(),
		Finish:    tNow.Add(5 * time.Second).UnixNano(),
		// ElapsedTime: 0,
		Archived: true,
	}

	assert.Nil(t, err)
	err = m.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	err = m.TimeSliceWrite(sliceUUID, sliceWrite)
	assert.Nil(t, err)
	timerWrite.ActiveSliceUUID = sliceUUID
	err = m.TimerWrite(timerUUID, timerWrite)
	assert.Nil(t, err)
	timerRead, err := m.TimerRead(timerUUID)
	assert.Nil(t, err)
	assert.Equal(t, timerWrite, timerRead)
	err = m.Shutdown()
	assert.Nil(t, err)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
