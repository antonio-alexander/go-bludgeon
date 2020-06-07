package bludgeonrestapi

import (
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type remote struct {
	address string
	port    string
	timeout time.Duration
}

func NewRemote(address, port string, timeout time.Duration) interface {
	bludgeon.Remote
} {

	//validate timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	//store timeout
	ConfigTimeout = timeout
	//create remote
	return &remote{
		address: address,
		port:    port,
		timeout: timeout,
	}
}

var _ bludgeon.Remote = &remote{}

//
func (r *remote) TimerCreate() (timer bludgeon.Timer, err error) {
	timer, err = TimerCreate(r.address, r.port)

	return
}

//
func (r *remote) TimerRead(id string) (timer bludgeon.Timer, err error) {
	timer, err = TimerRead(r.address, r.port, id)

	return
}

//
func (r *remote) TimerUpdate(t bludgeon.Timer) (timer bludgeon.Timer, err error) {
	timer, err = TimerUpdate(r.address, r.port, t)

	return
}

//
func (r *remote) TimerDelete(id string) (err error) {
	err = TimerDelete(r.address, r.port, id)

	return
}

//TimerStart
func (r *remote) TimerStart(timerID string, startTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerStart(r.address, r.port, timerID, startTime)

	return
}

//TimerPause
func (r *remote) TimerPause(timerID string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerPause(r.address, r.port, timerID, pauseTime)

	return
}

//TimerSubmit
func (r *remote) TimerSubmit(timerID string, finishTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerSubmit(r.address, r.port, timerID, finishTime)

	return
}

func (r *remote) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	timeSlice, err = TimeSliceRead(r.address, r.port, id)

	return
}
