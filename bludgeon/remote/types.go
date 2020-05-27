package bludgeonremote

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type Remote interface {
	//
	TimerCreate() (timer bludgeon.Timer, err error)
	//
	TimerRead(id string) (timer bludgeon.Timer, err error)

	//
	TimerUpdate(id string, timer bludgeon.Timer) (err error)

	//
	TimerDelete(id string) (err error)

	//
	TimeSliceCreate(timerid string) (id string, err error)

	//
	TimeSliceRead(id string) (timer bludgeon.TimeSlice, err error)

	//
	TimeSliceUpdate(id string, timer bludgeon.TimeSlice) (err error)

	//
	TimeSliceDelete(id string) (err error)
}
