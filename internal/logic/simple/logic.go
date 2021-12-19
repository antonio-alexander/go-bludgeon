package simple

import (
	"fmt"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logic "github.com/antonio-alexander/go-bludgeon/internal/logic"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
)

//REVIEW: as coded, this would be most efficient if there was no mutex, needing a mutex
// indicates that we're doing logic here that should be done within meta (using it's mutex)
// it may or may not be preferable in the long run to do it that way.

type logicSimple struct {
	sync.WaitGroup
	sync.RWMutex
	logger.Logger
	meta.Timer
	meta.TimeSlice
	stopper chan struct{}
	chCache <-chan data.CacheData
	started bool
}

func New(parameters ...interface{}) interface {
	logic.Logic
	logic.Functional
} {
	l := &logicSimple{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			meta.Timer
			meta.TimeSlice
		}:
			l.Timer = p
			l.TimeSlice = p
		case meta.Timer:
			l.Timer = p
		case meta.TimeSlice:
			l.TimeSlice = p
		case logger.Logger:
			l.Logger = p
		}
	}
	return l
}

func (l *logicSimple) Start() (err error) {
	l.Lock()
	defer l.Unlock()

	if l.started {
		return
	}
	l.stopper = make(chan struct{})
	l.chCache = make(<-chan data.CacheData)
	l.launchCache()
	l.started = true

	return
}

func (l *logicSimple) Stop() (err error) {
	l.Lock()
	defer l.Unlock()

	if !l.started {
		return
	}
	close(l.stopper)
	l.Wait()
	l.started = false

	return
}

func (l *logicSimple) launchCache() {
	//create channel to block until go routine enters business logic
	started := make(chan struct{})
	//launch go Routine
	l.Add(1)
	go func() {
		defer l.Done()

		//REVIEW: could just store the function call in a slice and have the business
		// logic just attempt to run and if successful, delete from the slice

		//TODO: create ticker to periodically do business logic
		//create signal or way to receive cached operations that
		// need to be synchronized
		//close started to indicate that business logic has started
		close(started)
		//start business logic
		for {
			select {
			//periodically attempt to execute business logic
			// signal to store cached operations
			case <-l.chCache:
			case <-l.stopper:
				return
			}
		}
	}()
	//wait for goRoutine to enter business logic
	<-started
}

func (l *logicSimple) TimerCreate() (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	if timer.UUID, err = data.GenerateID(); err != nil {
		return
	}
	if err = l.Timer.TimerWrite(timer.UUID, timer); err != nil {
		return
	}
	l.Debug("created timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimerRead(id string) (timer data.Timer, err error) {
	l.RLock()
	defer l.RUnlock()

	var timeSlice data.TimeSlice

	if timer, err = l.Timer.TimerRead(id); err != nil {
		return
	}
	if timer.ActiveSliceUUID != "" {
		if timeSlice, err = l.TimeSlice.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
		timer.ElapsedTime += (time.Now().UnixNano() - timeSlice.Start)
	}
	l.Debug("read timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimerUpdate(t data.Timer) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	//TODO: this function should be deprecated and replaced with
	// functions that edit properties instead
	if timer, err = l.Timer.TimerRead(t.UUID); err != nil {
		return
	}
	if t.Comment != "" {
		timer.Comment = t.Comment
	}
	if err = l.Timer.TimerWrite(timer.UUID, timer); err != nil {
		return
	}
	l.Debug("updated timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimerDelete(timerID string) (err error) {
	l.Lock()
	defer l.Unlock()

	if err = l.Timer.TimerDelete(timerID); err != nil {
		return
	}
	l.Debug("deleted timer with ID %s", timerID)

	return
}

func (l *logicSimple) TimerStart(id string, startTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//to start the timer, we need to first read it
	// and then set the active time slice, then update
	// the timer
	if timer, err = l.Timer.TimerRead(id); err != nil {
		return
	}
	if timer.Archived {
		return data.Timer{}, fmt.Errorf(data.ErrTimerIsArchivedf, timer.UUID)
	}
	if timer.Start == 0 {
		timer.Start = startTime.UnixNano()
	}
	if timer.ActiveSliceUUID != "" {
		if timeSlice, err = l.TimeSlice.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
	} else {
		if timeSlice.UUID, err = data.GenerateID(); err != nil {
			return
		}
		timeSlice.TimerUUID = timer.UUID
		if err = l.TimeSlice.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
			return
		}
		timeSlice.Start = startTime.UnixNano()
	}
	timer.ActiveSliceUUID = timeSlice.UUID
	if err = l.TimeSlice.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	if err = l.Timer.TimerWrite(timer.UUID, timer); err != nil {
		return
	}
	//REVIEW: if we create a GUID, we will need some way to overwrite the cached operation
	l.Debug("started timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimerPause(timerID string, pauseTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//if the given timer exists, when we pause the timer, we will
	// grab the active time slice, set the finish time to now
	// set the active slice value to -1 update the elapsed time,
	// add all the timers that aren't archived together
	if timer, err = l.Timer.TimerRead(timerID); err != nil {
		return
	}
	if timer.ActiveSliceUUID == "" {
		//REVIEW: would it make more sense to output an error that says
		// unable to pause a timer that isn't in progress?
		err = fmt.Errorf(data.ErrNoActiveTimeSlicef, timer.UUID)

		return
	}
	if timeSlice, err = l.TimeSlice.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
		return
	}
	timeSlice.Finish = pauseTime.UnixNano()
	timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
	timeSlice.Archived = true
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	if err = l.TimeSlice.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	timer.ActiveSliceUUID = ""
	if err = l.Timer.TimerWrite(timer.UUID, timer); err != nil {
		return
	}
	l.Debug("paused timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimerSubmit(timerID string, submitTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submitted" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)
	if timer, err = l.Timer.TimerRead(timerID); err != nil {
		return
	}
	if timer.ActiveSliceUUID != "" {
		if timeSlice, err = l.TimeSlice.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
		timeSlice.Finish = submitTime.UnixNano()
		timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
		timeSlice.Archived = true
		timer.ActiveSliceUUID = ""
	}
	timer.Finish = submitTime.UnixNano()
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	timer.Completed = true
	if err = l.TimeSlice.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	if err = l.Timer.TimerWrite(timer.UUID, timer); err != nil {
		return
	}
	l.Debug("submitted timer with ID %s", timer.UUID)

	return
}

func (l *logicSimple) TimeSliceRead(timeSliceID string) (timeSlice data.TimeSlice, err error) {
	l.RLock()
	defer l.RUnlock()

	if timeSlice, err = l.TimeSlice.TimeSliceRead(timeSliceID); err != nil {
		return
	}
	l.Debug("read time slice with ID %s", timeSliceID)

	return
}
