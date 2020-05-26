package bludgeonremotemock

import (
	"errors"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	remote "github.com/antonio-alexander/go-bludgeon/bludgeon/remote"
)

type remoteMock interface{}

func NewRemoteMock() interface {
	remote.Remote
} {
	return &remoteMock{}
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
func (m *remoteMock) TimerUpdate(id string, timer bludgeon.Timer) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimerDelete(id string) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceCreate(timerid string) (id string, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceRead(id string) (timer bludgeon.TimeSlice, err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceUpdate(id string, timer bludgeon.TimeSlice) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}

//
func (m *remoteMock) TimeSliceDelete(id string) (err error) {
	err = errors.New(ErrNotImplemented)

	return
}
