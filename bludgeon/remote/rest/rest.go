package bludgeonrest

import (
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	api "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/api"
)

type remote struct {
	sync.Mutex
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
	r.Lock()
	defer r.Unlock()

	timer, err = api.TimerCreate(r.address, r.port)

	return
}

//
func (r *remote) TimerRead(id string) (timer bludgeon.Timer, err error) {
	r.Lock()
	defer r.Unlock()

	timer, err = api.TimerRead(r.address, r.port, id)

	return
}

//
func (r *remote) TimerUpdate(timer bludgeon.Timer) (err error) {
	r.Lock()
	defer r.Unlock()

	err = api.TimerUpdate(r.address, r.port, timer)

	return
}

//
func (r *remote) TimerDelete(id string) (err error) {
	r.Lock()
	defer r.Unlock()

	err = api.TimerDelete(r.address, r.port, id)

	return
}

//
func (r *remote) TimeSliceCreate(id string) (timeSlice bludgeon.TimeSlice, err error) {
	r.Lock()
	defer r.Unlock()

	timeSlice, err = api.TimeSliceCreate(r.address, r.port, id)

	return
}

//
func (r *remote) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	r.Lock()
	defer r.Unlock()

	timeSlice, err = api.TimeSliceRead(r.address, r.port, id)

	return
}

//
func (r *remote) TimeSliceUpdate(timeSlice bludgeon.TimeSlice) (err error) {
	r.Lock()
	defer r.Unlock()

	err = api.TimeSliceUpdate(r.address, r.port, timeSlice)

	return
}

//
func (r *remote) TimeSliceDelete(id string) (err error) {
	r.Lock()
	defer r.Unlock()

	err = api.TimeSliceDelete(r.address, r.port, id)

	return
}
