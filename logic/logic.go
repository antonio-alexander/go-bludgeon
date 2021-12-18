package logic

import (
	"fmt"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
)

//REVIEW: as coded, this would be most efficient if there was no mutex, needing a mutex
// indicates that we're doing logic here that should be done within meta (using it's mutex)
// it may or may not be preferable in the long run to do it that way.

type logic struct {
	sync.WaitGroup
	sync.RWMutex
	data.Logger
	meta interface {
		meta.Timer
		meta.TimeSlice
	}
	stopper chan struct{}
	chCache <-chan data.CacheData
	started bool
}

func New(logger data.Logger,
	meta interface {
		meta.Timer
		meta.TimeSlice
	}) interface {
	Logic
	Functional
} {
	return &logic{
		meta:   meta,
		Logger: logger,
	}
}

func (l *logic) Start() (err error) {
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
func (l *logic) Stop() (err error) {
	l.Lock()
	defer l.Unlock()

	close(l.stopper)
	l.Wait()
	l.started = false

	return
}

//LaunchCache
func (l *logic) launchCache() {
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
			//signal to store cached operations
			case <-l.chCache:
			case <-l.stopper:
				return
			}
		}
	}()
	//wait for goRoutine to enter business logic
	<-started
}

func (l *logic) TimerCreate() (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	if timer.UUID, err = data.GenerateID(); err != nil {
		return
	}
	err = l.meta.TimerWrite(timer.UUID, timer)

	return
}

func (l *logic) TimerRead(id string) (timer data.Timer, err error) {
	l.RLock()
	defer l.RUnlock()

	var timeSlice data.TimeSlice

	//attempt to read the timer
	if timer, err = l.meta.TimerRead(id); err != nil {
		return
	}
	//attempt to get the activeSlice if it exists
	if timer.ActiveSliceUUID != "" {
		//read the time slice to store in meta if remote
		if timeSlice, err = l.meta.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
		//update elapsed time in realtime if there is an active timeslice
		timer.ElapsedTime += (time.Now().UnixNano() - timeSlice.Start)
	}

	return
}

func (l *logic) TimerUpdate(t data.Timer) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	//get the current timer
	if timer, err = l.meta.TimerRead(t.UUID); err != nil {
		return
	}
	//update only the items that have changed
	// update the comment
	if t.Comment != "" {
		timer.Comment = t.Comment
	}
	//store the timeslice
	err = l.meta.TimerWrite(timer.UUID, timer)

	return
}

func (l *logic) TimerDelete(timerID string) (err error) {
	l.Lock()
	defer l.Unlock()

	if err = l.meta.TimerDelete(timerID); err != nil {
		return
	}

	return
}

func (l *logic) TimerStart(id string, startTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//attempt to read the timer
	if timer, err = l.meta.TimerRead(id); err != nil {
		return
	}
	//check if timer is archived
	if timer.Archived {
		err = fmt.Errorf(data.ErrTimerIsArchivedf, timer.UUID)

		return
	}
	//check if the timer start is empty
	if timer.Start == 0 {
		//set the start time to now
		timer.Start = startTime.UnixNano()
	}
	//check if there's an active slice
	if timer.ActiveSliceUUID != "" {
		//read the active slice (e.g. resume the timer so you don't lose the slice if it was
		// stopped ungracefully)
		if timeSlice, err = l.meta.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
	} else {
		//generate the time slice id
		if timeSlice.UUID, err = data.GenerateID(); err != nil {
			return
		}
		//update the time slice's timer ID
		timeSlice.TimerUUID = timer.UUID
		//store the timeslice
		if err = l.meta.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
			//TODO: store the meta error
			return
		}
		//set the start time
		timeSlice.Start = startTime.UnixNano()
	}
	//store the timeSlice
	timer.ActiveSliceUUID = timeSlice.UUID
	//store the time slice
	if err = l.meta.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		//TODO: store meta error
		return
	}
	//store the timer
	if err = l.meta.TimerWrite(timer.UUID, timer); err != nil {
		//TODO: store meta error
		return
	}
	//REVIEW: if we create a GUID, we will need some way to overwrite the cached operation

	return
}

func (l *logic) TimerPause(timerID string, pauseTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//if the given timer exists, when we pause the timer, we will
	// grab the active time slice, set the finish time to now
	// set the active slice value to -1 update the elapsed time,
	// add all the timers that aren't archived together

	//attempt to read the timer
	if timer, err = l.meta.TimerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is, do nothing
	// if there isn't, generate an error since a paused timer
	// should be in-progress (noted by an active slice)
	if timer.ActiveSliceUUID == "" {
		//REVIEW: would it make more sense to output an error that says
		// unable to pause a timer that isn't in progress?
		err = fmt.Errorf(data.ErrNoActiveTimeSlicef, timer.UUID)

		return
	}
	//read the active time slice
	if timeSlice, err = l.meta.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
		return
	}
	//set the finish time
	timeSlice.Finish = pauseTime.UnixNano()
	//calculate the elapsed time
	timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
	//set archived to true
	timeSlice.Archived = true
	//calculate elapsed time
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	//update the timeslice
	if err = l.meta.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//update the timer, remove its time slice and update the elapsed time
	//set the activeSliceUUID to empty
	timer.ActiveSliceUUID = ""
	//update the timer
	if err = l.meta.TimerWrite(timer.UUID, timer); err != nil {
		return
	}

	return
}

func (l *logic) TimerSubmit(timerID string, submitTime time.Time) (timer data.Timer, err error) {
	l.Lock()
	defer l.Unlock()

	var timeSlice data.TimeSlice

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submitted" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)
	if timer, err = l.meta.TimerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is,
	if timer.ActiveSliceUUID != "" {
		//read the active time slice
		if timeSlice, err = l.meta.TimeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
		//set the finish time
		timeSlice.Finish = submitTime.UnixNano()
		//calculate the elapsed time
		timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
		//set archived to true
		timeSlice.Archived = true
		//update the timer, remove its time slice and update the elapsed time
		//set the activeSliceUUID to empty
		timer.ActiveSliceUUID = ""
	}
	//set finish time
	timer.Finish = submitTime.UnixNano()
	//calculate elapsed time
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	//set completed
	timer.Completed = true
	//update the timeslice
	if err = l.meta.TimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//update the timer
	if err = l.meta.TimerWrite(timer.UUID, timer); err != nil {
		return
	}

	return
}

func (l *logic) TimeSliceRead(timeSliceID string) (timeSlice data.TimeSlice, err error) {
	l.RLock()
	defer l.RUnlock()

	return l.meta.TimeSliceRead(timeSliceID)
}
