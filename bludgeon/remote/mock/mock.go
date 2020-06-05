package bludgeonremotemock

import (
	"errors"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type remoteMock struct{}

func NewRemoteMock() interface {
	bludgeon.Remote
} {
	return &remoteMock{
		//
	}
}

//
func (m *remoteMock) TimerCreate() (timer bludgeon.Timer, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimerRead(id string) (timer bludgeon.Timer, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimerUpdate(timer bludgeon.Timer) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimerDelete(id string) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//TimerStart
func (m *remoteMock) TimerStart(timerID string, startTime time.Time) (timer bludgeon.Timer, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//TimerPause
func (m *remoteMock) TimerPause(timerID string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//TimerSubmit
func (m *remoteMock) TimerSubmit(timerID string, finishTime time.Time) (timer bludgeon.Timer, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceCreate(timerid string) (timeSlice bludgeon.TimeSlice, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceRead(id string) (timer bludgeon.TimeSlice, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceUpdate(timeSlice bludgeon.TimeSlice) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceDelete(id string) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}
