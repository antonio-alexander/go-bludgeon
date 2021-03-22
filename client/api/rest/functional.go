package restapi

import (
	"errors"
	"sync"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
)

type functional struct {
	sync.RWMutex
	address string
	port    string
}

func NewFunctional() interface {
	common.FunctionalOwner
	common.FunctionalTimer
	common.FunctionalTimeSlice
} {
	return &functional{}
}

var _ common.FunctionalOwner = &functional{}

//Initialize
func (f *functional) Initialize(element interface{}) (err error) {
	f.RLock()
	defer f.RUnlock()

	var config Configuration
	var ok bool

	if config, ok = element.(Configuration); !ok {
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

var _ common.FunctionalTimer = &functional{}

func (f *functional) TimerCreate() (timer common.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerCreate(f.address, f.port)

	return
}

func (f *functional) TimerRead(id string) (timer common.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerRead(f.address, f.port, id)

	return
}

func (f *functional) TimerUpdate(timerIn common.Timer) (timerOut common.Timer, err error) {
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

func (f *functional) TimerStart(id string, startTime time.Time) (timer common.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerStart(f.address, f.port, id, startTime)

	return
}

func (f *functional) TimerPause(id string, pauseTime time.Time) (timer common.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerPause(f.address, f.port, id, pauseTime)

	return
}

func (f *functional) TimerSubmit(timerID string, finishTime time.Time) (timer common.Timer, err error) {
	f.RLock()
	defer f.RUnlock()

	timer, err = TimerSubmit(f.address, f.port, timerID, finishTime)

	return
}

var _ common.FunctionalTimeSlice = &functional{}

func (f *functional) TimeSliceRead(id string) (timeSlice common.TimeSlice, err error) {
	f.RLock()
	defer f.RUnlock()

	timeSlice, err = TimeSliceRead(f.address, f.port, id)

	return
}
