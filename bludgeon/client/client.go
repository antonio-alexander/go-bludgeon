package bludgeonclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	api "github.com/antonio-alexander/go-bludgeon/bludgeon/api"
	meta "github.com/antonio-alexander/go-bludgeon/bludgeon/meta"
	remote "github.com/antonio-alexander/go-bludgeon/bludgeon/remote"
)

//REVIEW: the future version of the client, will have a remote and meta pointer provided
// at new, at least one must not be nil, but if only remote, it'll use the server for all
// operations, if remote is nil but meta is not, it'll use meta instead for local persistence

type client struct {
	sync.RWMutex                 //mutex for threadsafe operations
	sync.WaitGroup               //waitgroup to track go routines
	started        bool          //whether or not the business logic has starting
	config         Configuration //configuration
	meta           interface {   //storage interface
		meta.Meta
		meta.MetaTimer
		meta.MetaTimeSlice
	}
	remote interface { //remote interface
		remote.Remote
	}
	lookupTimers     map[string]string //lookup for timers
	lookupTimeSlices map[string]string //lookup for time slices
	stopper          chan struct{}     //stopper to stop goRoutines
	//some set of flags or method to know what to update
}

func NewClient(meta interface {
	meta.Meta
	meta.MetaTimer
	meta.MetaTimeSlice
}, remote interface {
	remote.Remote
}) interface {
	Owner
	Manage
	Functional
	API
} {
	//create internal maps
	lookupTimers := make(map[string]string)
	lookupTimeSlices := make(map[string]string)
	//TODO: need a way to prepopulate or generate the
	// lookups from scratch? will need to have the full data some how
	//populate client pointer
	return &client{
		meta:             meta,
		remote:           remote,
		lookupTimers:     lookupTimers,
		lookupTimeSlices: lookupTimeSlices,
	}
}

//
func (c *client) timeSliceCreate(timerID string) (timeSlice bludgeon.TimeSlice, err error) {
	//use the api to create a time slice
	if timeSlice.UUID, err = api.TimeSliceCreate(timerID); err != nil {
		//TODO: cache the operation
		//generate the time slice id
		if timeSlice.UUID, err = bludgeon.GenerateID(); err != nil {
			return
		}
	}
	//cache the timeslice lookup
	c.lookupTimeSlices[timeSlice.UUID] = timeSlice.UUID
	//store the timeslice
	err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)

	return
}

//
func (c *client) timeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	//attempt to use the api to query the timeslice
	if timeSlice, err = api.TimeSliceRead(timeSliceID); err != nil {
		//TODO: attempt to use the api to cache the active slice
		//attempt to read the timeSlice locally
		if timeSlice, err = c.meta.MetaTimeSliceRead(timeSliceID); err != nil {
			return
		}
	}
	//add the time slice to the lookup
	c.lookupTimeSlices[timeSliceID] = timeSliceID

	return
}

//
func (c *client) timerRead(timerID string) (timer bludgeon.Timer, err error) {
	//attempt to use the api to get the provided timer
	if timer, err = api.TimerRead(timerID); err != nil {
		//attempt to read the timer locally
		if timer, err = c.meta.MetaTimerRead(timerID); err != nil {
			return
		}
	} else {
		//cache the lookup
		c.lookupTimers[timer.UUID] = timer.UUID
		//store the timer in the local cache
		err = c.meta.MetaTimerWrite(timerID, timer)
	}

	return
}

type Owner interface {
	//Close
	Close()

	//Serialize
	Serialize() (bytes []byte, err error)

	//Deserialize
	Deserialize(bytes []byte) (err error)
}

func (c *client) Close() {
	c.Lock()
	defer c.Unlock()

	//use the remote and/or storage functions to:
	// attempt to synchronize the current values
	// attempt to serialize what's remaining

	//set internal configuration to default
	c.config = Configuration{}
	//set internal pointers to nil
	c.lookupTimeSlices, c.lookupTimers = nil, nil
	c.meta, c.remote = nil, nil

	return
}

//Serialize
func (c *client) Serialize() (bytes []byte, err error) {
	c.RLock()
	defer c.RUnlock()

	//serialize data into bytes
	bytes, err = json.Marshal(&SerializedData{
		LookupTimers:     c.lookupTimers,
		LookupTimeSlices: c.lookupTimeSlices,
	})

	return
}

//Deserialize
func (c *client) Deserialize(bytes []byte) (err error) {
	c.Lock()
	defer c.Unlock()

	var serializedData SerializedData

	//check to see if started, otherwise, don't deserialzie the bytes
	if c.started {
		err = errors.New(ErrStarted)

		return
	}
	//attempt to de-serialize the data
	if err = json.Unmarshal(bytes, &serializedData); err != nil {
		return
	}
	//store de-serialzied time slices
	//range over timers and populate lookup
	for key, value := range serializedData.LookupTimeSlices {
		c.lookupTimers[key] = value
	}
	//range over timers and populate lookup
	for key, value := range c.lookupTimeSlices {
		c.lookupTimeSlices[key] = value
	}

	return
}

type Manage interface {
	//
	Start(config Configuration) (err error)

	//
	Stop() (err error)
}

func (c *client) Start(config Configuration) (err error) {
	c.Lock()
	defer c.Unlock()

	//check if already started
	if c.started {
		err = errors.New(ErrStarted)

		return
	}
	//store configuration
	c.config = config
	//create stopper
	c.stopper = make(chan struct{})
	//launch go routines
	c.LaunchCache()
	//set started to true
	c.started = true

	return
}

func (c *client) Stop() (err error) {
	c.Lock()
	defer c.Unlock()

	//check if not started
	if !c.started {
		err = errors.New(ErrNotStarted)

		return
	}
	//close the stopper
	close(c.stopper)
	c.Wait()
	//set started flag to false
	c.started = false

	return
}

type Functional interface {
	//LaunchCache
	LaunchCache()

	//
	CommandHandler(command bludgeon.CommandClient, dataIn interface{}) (dataOut interface{}, err error)
}

//ensure that cache implements Functional
var _ Functional = &client{}

//Launch for cache logic
func (c *client) LaunchCache() {
	started := make(chan struct{})
	c.Add(1)
	go c.goCache(started)
	<-started
}

//go routine for cache logic
func (c *client) goCache(started chan struct{}) {
	defer c.Done()

	//idea
	// could just store the function call in a slice and have the business
	// logic just attempt to run and if successful, delete from the slice

	//create ticker to periodically do business logic
	//create signal or way to receive cached operations that
	// need to be synchronized
	//close started to indicate that business logic has started
	close(started)
	//start business logic
	for {
		select {
		//periodically attempt to execute business logic
		//signal to store cached operations
		case <-c.stopper:
			return
		}
	}
}

func (c *client) CommandHandler(command bludgeon.CommandClient, dataIn interface{}) (dataOut interface{}, err error) {
	//execute the command
	switch command {
	case bludgeon.CommandClientTimerCreate:
		dataOut, err = c.TimerCreate()
	case bludgeon.CommandClientTimerRead:
		var id string

		//switch on the data
		switch v := dataIn.(type) {
		case bludgeon.Timer:
			id = v.UUID
		case string:
			id = v
		default:
			//generate error
		}
		//check
		if err == nil {
			dataOut, err = c.TimerRead(id)
		}
	default:
		//TODO: generate error
	}

	return
}

//API
type API interface {
	//add a timer
	TimerCreate() (timer bludgeon.Timer, err error)

	//get a timer
	TimerRead(timerID string) (timer bludgeon.Timer, err error)

	//delete a timer
	TimerDelete(timerID string) (err error)

	//start a timer
	TimerStart(timerID string, startTime time.Time) (err error)

	//pause a timer
	TimerPause(timerID string, pauseTime time.Time) (err error)

	//submit a timer
	TimerSubmit(timerID string, submitTime time.Time) (err error)
}

//ensure that cache implements API
var _ API = &client{}

//add a timer
func (c *client) TimerCreate() (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//attempt execute the api
	if timer, err = api.TimerCreate(); err != nil {
		//cache the operation
		// attempt to create a timer

		//generate a UUID for the timer locally
		if timer.UUID, err = bludgeon.GenerateID(); err != nil {
			//unable to generate the timer, quit and give an error
			return
		}
	}
	// update the lookup for the timer
	c.lookupTimers[timer.UUID] = timer.UUID
	//check if meta is not nil
	if c.meta != nil {
		//update the timer in meta
		if err = c.meta.MetaTimerWrite(timer.UUID, timer); err != nil {
			return
		}
	}

	return
}

//get a timer
func (c *client) TimerRead(id string) (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//attempt to read the timer
	if timer, err = c.timerRead(id); err != nil {
		//use meta to write to the timer
		//check if meta is not nil
		if c.meta != nil {
			//read the timer in meta
			if timer, err = c.meta.MetaTimerRead(id); err != nil {
				return
			}
		}

		return
	}
	//REVIEW: how doe sthis work with meta?
	//attempt to get the activeSlice if it exists
	if timer.ActiveSliceUUID != "" {
		//attempt to read the slice
		if _, err = c.timeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
	}

	return
}

//start a timer
func (c *client) TimerStart(id string, startTime time.Time) (err error) {
	c.Lock()
	defer c.Unlock()

	var timer bludgeon.Timer
	var timeSlice bludgeon.TimeSlice

	//attempt to read the timer
	if timer, err = c.timerRead(id); err != nil {
		return
	}
	//check if timer is archived
	if timer.Archived {
		err = fmt.Errorf(ErrTimerIsArchivedf, timer.UUID)

		return
	}
	//check if the timer start is empty
	if timer.Start == 0 {
		//set the start time to now
		timer.Start = startTime.UnixNano()
	}
	//check if there's an active slice, if there is, do nothing
	// if there isn't, create an active slice
	if timer.ActiveSliceUUID != "" {
		if _, err = c.timeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
	} else {
		//create the timeSlice
		if timeSlice, err = c.timeSliceCreate(timer.UUID); err != nil {
			return
		}
		//set the start time
		timeSlice.Start = startTime.UnixNano()
		//update the slice using the API
		if err = api.TimeSliceUpdate(timeSlice.UUID, timeSlice); err != nil {
			//cache operation

			//cache the timeSlice
			c.timeSlices[timeSlice.UUID] = timeSlice
			//cache the lookup
			c.lookupTimeSlices[timeSlice.UUID] = timeSlice.UUID
			//clear the error
			err = nil
		}
		//store the timeSlice
		timer.ActiveSliceUUID = timeSlice.UUID
		//update the timer
		if err = api.TimerUpdate(timer.UUID, timer); err != nil {
			//cache the operation

			//cache the lookup
			c.lookupTimers[timer.UUID] = timer.UUID
			//cache the timer
			c.timers[timer.UUID] = timer
			err = nil
		}
	}

	return
}

//pause a timer
func (c *client) TimerPause(timerID string, pauseTime time.Time) (err error) {
	c.Lock()
	defer c.Unlock()

	//if the given timer exists, when we pause the timer, we will
	// grab the active time slice, set the finish time to now
	// set the active slice value to -1 update the elapsed time,
	// add all the timers that aren't archived together

	var timer bludgeon.Timer
	var timeSlice bludgeon.TimeSlice

	//attempt to use the api to get the provided timer
	if timer, err = c.timerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is, do nothing
	// if there isn't, create an active slice
	if timer.ActiveSliceUUID != "" {
		err = fmt.Errorf(ErrNoActiveTimeSlicef, timer.UUID)

		return
	}
	//read the active time slice
	if timeSlice, err = c.timeSliceRead(timer.ActiveSliceUUID); err != nil {
		return
	}
	//set the finish time
	timeSlice.Finish = pauseTime.UnixNano()
	//calculate the elapsed time
	timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
	//set archived to true
	timeSlice.Archived = true
	//update the slice using the API, there's an understanding that this should
	// trigger some server side logic to update the timer
	if err = api.TimeSliceUpdate(timeSlice.UUID, timeSlice); err != nil {
		//cache the operation
		// because the service side logic should complete, we don't really have to do anything
		// since the logic is going to delete the slice anyway

		//update timeslice for completeness
		c.timeSlices[timeSlice.UUID] = timeSlice
		//clear error
		err = nil
	} else {
		//remove the time slice since we know the server has the most up to date
		// copy
		delete(c.lookupTimeSlices, timeSlice.UUID)
		delete(c.timeSlices, timeSlice.UUID)
	}
	//update the timer, remove its time slice and update the elapsed time
	//set the activeSliceUUID to empty
	timer.ActiveSliceUUID = ""
	//update the timer
	if err = api.TimerUpdate(timer.UUID, timer); err != nil {
		//cache the operation

		//server logic should update the elapsed time
		//update the timer locally
		//calculate elapsed time
		timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
		//update the timer
		c.timers[timer.UUID] = timer
	} else {
		//re-read the timer
		_, err = c.timerRead(timer.UUID)
	}

	return
}

//submit a timer
func (c *client) TimerSubmit(timerID string, submitTime time.Time) (err error) {
	c.Lock()
	defer c.Unlock()

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submittd" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)

	var timer bludgeon.Timer
	var timeSlice bludgeon.TimeSlice

	//attempt to use the api to get the provided timer
	if timer, err = c.timerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is, do nothing
	// if there isn't, create an active slice
	if timer.ActiveSliceUUID != "" {
		//read the active time slice
		if timeSlice, err = c.timeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
		//set the finish time
		timeSlice.Finish = submitTime.UnixNano()
		//calculate the elapsed time
		timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
		//set archived to true
		timeSlice.Archived = true
		//update the slice using the API, there's an understanding that this should
		// trigger some server side logic to update the timer
		if err = api.TimeSliceUpdate(timeSlice.UUID, timeSlice); err != nil {
			//cache the operation
			// because the service side logic should complete, we don't really have to do anything
			// since the logic is going to delete the slice anyway

			//update timeslice for completeness
			c.timeSlices[timeSlice.UUID] = timeSlice
			//clear error
			err = nil
		} else {
			//remove the time slice since we know the server has the most up to date
			// copy
			delete(c.lookupTimeSlices, timeSlice.UUID)
			delete(c.timeSlices, timeSlice.UUID)
		}
		//update the timer, remove its time slice and update the elapsed time
		//set the activeSliceUUID to empty
		timer.ActiveSliceUUID = ""
	}
	//set finish time
	timer.Finish = submitTime.UnixNano()
	//update the timer
	if err = api.TimerUpdate(timer.UUID, timer); err != nil {
		//cache the operation

		//server logic should update the elapsed time
		//update the timer locally
		//calculate elapsed time
		timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
		//update the timer
		c.timers[timer.UUID] = timer
	} else {
		//re-read the timer
		_, err = c.timerRead(timer.UUID)
	}

	return
}

//delete a timer
func (c *client) TimerDelete(timerID string) (err error) {
	c.Lock()
	defer c.Unlock()

	var timer bludgeon.Timer

	//attempt to use the api to get the provided timer
	if timer, err = c.timerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice
	if timer.ActiveSliceUUID != "" {
		//delete the active slice
	}
	timer.Archived = true
	//update the timer
	if err = api.TimerUpdate(timer.UUID, timer); err != nil {
		//cache the operation

		//update the timer
		c.timers[timer.UUID] = timer
	}

	return
}
