package bludgeonrestapi

import (
	"errors"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/config"
)

type functional struct {
	sync.RWMutex
	address string
	port    string
}

func NewFunctional() interface {
	bludgeon.FunctionalOwner
	bludgeon.FunctionalTimer
	bludgeon.FunctionalTimeSlice
} {
	return &functional{}
}

var _ bludgeon.FunctionalOwner = &functional{}

//Initialize
func (f *functional) Initialize(element interface{}) (err error) {
	f.RLock()
	defer f.RUnlock()

	var config rest.Configuration
	var ok bool

	if config, ok = element.(rest.Configuration); !ok {
		err = errors.New("Unable to cast into configuration")

		return
	}
	f.address, f.port = config.Address, config.Port

	return
}

//Shutdown
func (f *functional) Shutdown() (err error) {
	f.RLock()
	defer f.RUnlock()

	f.address, f.port = "", ""

	return
}

var _ bludgeon.FunctionalTimer = &functional{}

func (f *functional) TimerCreate() (timer bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerCreate(f.address, f.port)

	return
}

func (f *functional) TimerRead(id string) (timer bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerRead(f.address, f.port, id)

	return
}

func (f *functional) TimerUpdate(timerIn bludgeon.Timer) (timerOut bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timerOut, err = TimerUpdate(f.address, f.port, timerIn)

	return
}

func (f *functional) TimerDelete(id string) (err error) {
	f.RLock()
	defer f.RUnlock()

	err = TimerDelete(f.address, f.port, id)

	return
}

func (f *functional) TimerStart(id string, startTime time.Time) (timer bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerStart(f.address, f.port, id, startTime)

	return
}

func (f *functional) TimerPause(id string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerPause(f.address, f.port, id, pauseTime)

	return
}

func (f *functional) TimerSubmit(timerID string, finishTime time.Time) (timer bludgeon.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerSubmit(f.address, f.port, timerID, finishTime)

	return
}

var _ bludgeon.FunctionalTimeSlice = &functional{}

func (f *functional) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	f.RLock()
	defer f.RUnlock()

	timeSlice, err = TimeSliceRead(f.address, f.port, id)

	return
}
