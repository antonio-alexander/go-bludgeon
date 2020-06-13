package bludgeonrestapi

import (
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type remote struct {
	sync.Mutex
	config Configuration
}

func NewRemote() interface {
	bludgeon.RemoteOwner
	bludgeon.RemoteTimer
	bludgeon.RemoteTimeSlice
} {

	//create remote
	return &remote{}
}

var _ bludgeon.RemoteOwner = &remote{}

//Initialize
func (r *remote) Initialize(element interface{}) (err error) {
	r.Lock()
	defer r.Unlock()

	var config Configuration

	//attempt to cast configuration
	if config, err = castConfiguration(element); err != nil {
		return
	}
	//store configuration
	//validate timeout
	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout
	}
	//store timeout for api
	ConfigTimeout = config.Timeout
	//store configuration
	r.config = config

	return
}

//Shutdown
func (r *remote) Shutdown() (err error) {
	r.Lock()
	defer r.Unlock()

	return
}

var _ bludgeon.RemoteTimer = &remote{}

//
func (r *remote) TimerCreate() (timer bludgeon.Timer, err error) {
	timer, err = TimerCreate(r.config.Address, r.config.Port)

	return
}

//
func (r *remote) TimerRead(id string) (timer bludgeon.Timer, err error) {
	timer, err = TimerRead(r.config.Address, r.config.Port, id)

	return
}

//
func (r *remote) TimerUpdate(t bludgeon.Timer) (timer bludgeon.Timer, err error) {
	timer, err = TimerUpdate(r.config.Address, r.config.Port, t)

	return
}

//
func (r *remote) TimerDelete(id string) (err error) {
	err = TimerDelete(r.config.Address, r.config.Port, id)

	return
}

//TimerStart
func (r *remote) TimerStart(timerID string, startTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerStart(r.config.Address, r.config.Port, timerID, startTime)

	return
}

//TimerPause
func (r *remote) TimerPause(timerID string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerPause(r.config.Address, r.config.Port, timerID, pauseTime)

	return
}

//TimerSubmit
func (r *remote) TimerSubmit(timerID string, finishTime time.Time) (timer bludgeon.Timer, err error) {
	timer, err = TimerSubmit(r.config.Address, r.config.Port, timerID, finishTime)

	return
}

var _ bludgeon.RemoteTimeSlice = &remote{}

func (r *remote) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	timeSlice, err = TimeSliceRead(r.config.Address, r.config.Port, id)

	return
}
