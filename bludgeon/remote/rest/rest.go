package bludgeonrest

import (
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	api "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/api"
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
		timeout = api.DefaultTimeout
	}
	//store timeout
	api.ConfigTimeout = timeout
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
	timer, err = api.TimerCreate(r.address, r.port)

	return
}

//
func (r *remote) TimerRead(id string) (timer bludgeon.Timer, err error) {
	timer, err = api.TimerRead(r.address, r.port, id)

	return
}

//
func (r *remote) TimerUpdate(timer bludgeon.Timer) (err error) {
	err = api.TimerUpdate(r.address, r.port, timer)

	return
}

//
func (r *remote) TimerDelete(id string) (err error) {
	err = api.TimerDelete(r.address, r.port, id)

	return
}

//TimerStart
func (r *remote) TimerStart(timerID string, startTime time.Time) (err error) {
	err = api.TimerStart(r.address, r.port, timerID, startTime)

	return
}

//TimerPause
func (r *remote) TimerPause(timerID string, pauseTime time.Time) (err error) {
	err = api.TimerPause(r.address, r.port, timerID, pauseTime)

	return
}

//TimerSubmit
func (r *remote) TimerSubmit(timerID string, finishTime time.Time) (err error) {
	err = api.TimerSubmit(r.address, r.port, timerID, finishTime)

	return
}

//
func (r *remote) TimeSliceCreate(id string) (timeSlice bludgeon.TimeSlice, err error) {
	timeSlice, err = api.TimeSliceCreate(r.address, r.port, id)

	return
}

//
func (r *remote) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	timeSlice, err = api.TimeSliceRead(r.address, r.port, id)

	return
}

//
func (r *remote) TimeSliceUpdate(timeSlice bludgeon.TimeSlice) (err error) {
	err = api.TimeSliceUpdate(r.address, r.port, timeSlice)

	return
}

//
func (r *remote) TimeSliceDelete(id string) (err error) {
	err = api.TimeSliceDelete(r.address, r.port, id)

	return
}
