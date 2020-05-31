package bludgeonclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
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
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
	}
	remote interface { //remote interface
		bludgeon.Remote
	}
	lookupTimers     map[string]string //lookup for timers
	lookupTimeSlices map[string]string //lookup for time slices
	stopper          chan struct{}     //stopper to stop goRoutines
	//some set of flags or method to know what to update
}

func NewClient(meta interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, remote interface {
	bludgeon.Remote
}) interface {
	Owner
	Manage
	Functional
	// API
} {
	//REVIEW: should add cases to confirm that meta/remote aren't nil since
	// basic functionality won't work?
	if meta == nil {
		panic("meta is nil")
	}
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
func (c *client) timeSliceCreate(timerUUID string) (timeSlice bludgeon.TimeSlice, err error) {
	//use the api to create a time slice if remote is not nil
	if c.remote != nil {
		if timeSlice, err = c.remote.TimeSliceCreate(timerUUID); err != nil {
			//TODO: cache the operation
		}
	}
	//REVIEW: what if the first remote call succeeds, but the second one fails? should probably update the api to include
	//only generate the idea if the remote operations failed
	if err != nil || c.remote == nil {
		//generate the time slice id
		if timeSlice.UUID, err = bludgeon.GenerateID(); err != nil {
			return
		}
		//update the time slice's timer ID
		timeSlice.TimerUUID = timerUUID
	}
	//cache the timeslice lookup
	c.lookupTimeSlices[timeSlice.UUID] = timeSlice.UUID
	//store the timeslice
	err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)

	return
}

//
func (c *client) timeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	if c.remote != nil {
		//attempt to use the api to query the timeslice
		if timeSlice, err = c.remote.TimeSliceRead(timeSliceID); err != nil {
			//TODO: cache attempt to read time slice?
		} else {
			//write time slice to meta
			if err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
				return
			}
		}
	}
	//only read from meta if remote is nil or there's an error
	if c.remote == nil || err != nil {
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
func (c *client) timeSliceUpdate(timeSlice bludgeon.TimeSlice) (err error) {
	//use the api to create a time slice if remote is not nil
	if c.remote != nil {
		if err = c.remote.TimeSliceUpdate(timeSlice); err != nil {
			//TODO: cache the operation
		}
	}
	//store the timeslice
	err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)

	return
}

//
func (c *client) timerRead(timerID string) (timer bludgeon.Timer, err error) {
	//check if remote is nil
	if c.remote != nil {
		//attempt to use the api to get the provided timer
		if timer, err = c.remote.TimerRead(timerID); err != nil {
			//TODO: cache the read?
		} else {
			//store the timer in the local cache
			if err = c.meta.MetaTimerWrite(timerID, timer); err != nil {
				return
			}
		}
	}
	//read the meta timer if remote is nil or there is no err
	if err != nil || c.remote == nil {
		//attempt to read the timer locally
		if timer, err = c.meta.MetaTimerRead(timerID); err != nil {
			return
		}
	}
	//cache the lookup
	c.lookupTimers[timer.UUID] = timer.UUID

	return
}

func (c *client) timerUpdate(timer bludgeon.Timer) (err error) {
	//use the api to create a time slice if remote is not nil
	if c.remote != nil {
		if err = c.remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//store the timeslice
	err = c.meta.MetaTimerWrite(timer.UUID, timer)

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

	//check to see if started, otherwise, don't de-serialize the bytes
	if c.started {
		err = errors.New(ErrStarted)

		return
	}
	//attempt to de-serialize the data
	if err = json.Unmarshal(bytes, &serializedData); err != nil {
		return
	}
	//store de-serialized time slices
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

	//CommandHandler
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
	var id string
	var ok bool

	//execute the command
	switch command {
	case bludgeon.CommandClientTimerCreate:
		dataOut, err = c.TimerCreate()
	case bludgeon.CommandClientTimerStart:
		if id, ok = dataIn.(string); !ok {
			//TODO: generate error
			return
		}
		//start the timer
		err = c.TimerStart(id, time.Now())
	case bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerStop:
		if id, ok = dataIn.(string); !ok {
			//TODO: generate error
			return
		}
		//start the timer
		err = c.TimerPause(id, time.Now())
	case bludgeon.CommandClientTimerRead:
		//cast into string and use
		if id, ok = dataIn.(string); !ok {
			//TODO: generate error
			return
		}
		//read the timer
		dataOut, err = c.TimerRead(id)
	case bludgeon.CommandClientTimerSubmit:
		if id, ok = dataIn.(string); !ok {
			//TODO: generate error
			return
		}
		//start the timer
		err = c.TimerSubmit(id, time.Now())
	default:
		//TODO: generate error
	}

	return
}

// //API
// type API interface {
// 	//add a timer
// 	TimerCreate() (timer bludgeon.Timer, err error)

// 	//get a timer
// 	TimerRead(timerID string) (timer bludgeon.Timer, err error)

// 	//delete a timer
// 	TimerDelete(timerID string) (err error)

// 	//start a timer
// 	TimerStart(timerID string, startTime time.Time) (err error)

// 	//pause a timer
// 	TimerPause(timerID string, pauseTime time.Time) (err error)

// 	//submit a timer
// 	TimerSubmit(timerID string, submitTime time.Time) (err error)
// }

// //ensure that cache implements API
// var _ API = &client{}

//add a timer
func (c *client) TimerCreate() (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//check if remote is nil
	if c.remote != nil {
		//attempt execute the api
		if timer, err = c.remote.TimerCreate(); err != nil {
			//TODO: cache the operation
		}
	}
	//generate the uuid if there's a remote error or if there is no remote
	if err != nil || c.remote == nil {
		//generate a UUID for the timer locally
		if timer.UUID, err = bludgeon.GenerateID(); err != nil {
			return
		}
	}
	// update the lookup for the timer
	c.lookupTimers[timer.UUID] = timer.UUID
	//update the timer in meta
	if err = c.meta.MetaTimerWrite(timer.UUID, timer); err != nil {
		return
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
	//REVIEW: how does this work with meta?
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
	//check if there's an active slice
	if timer.ActiveSliceUUID != "" {
		//read the active slice (e.g. resume the timer so you don't lose the slice if it was
		// stopped ungracefully)
		if timeSlice, err = c.timeSliceRead(timer.ActiveSliceUUID); err != nil {
			return
		}
	} else {
		//since there is no active slice, create the timeSlice
		if timeSlice, err = c.timeSliceCreate(timer.UUID); err != nil {
			return
		}
		//set the start time
		timeSlice.Start = startTime.UnixNano()
		//update the time slice with its new start time
		if c.remote != nil {
			//update the slice using the API
			if err = c.remote.TimeSliceUpdate(timeSlice); err != nil {
				//TODO: cache operation
			}
		}
	}
	//store the slice
	if err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//cache the lookup
	c.lookupTimeSlices[timeSlice.UUID] = timeSlice.UUID
	//store the timeSlice
	timer.ActiveSliceUUID = timeSlice.UUID
	if c.remote != nil {
		//update the timer
		if err = c.remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//cache the lookup
	c.lookupTimers[timer.UUID] = timer.UUID
	//cache the timer
	err = c.meta.MetaTimerWrite(timer.UUID, timer)

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

	//attempt to read the timer
	if timer, err = c.timerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is, do nothing
	// if there isn't, generate an error since a paused timer
	// should be in-progress (noted by an active slice)
	if timer.ActiveSliceUUID == "" {
		//REVIEW: would it make more sense to output an error that says
		// unable to pause a timer that isn't in progress?
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
	if c.remote != nil {
		//update the time slice
		if err = c.remote.TimeSliceUpdate(timeSlice); err != nil {
			//TODO: cache the operation
		}
	}
	// because the service side logic should complete, we don't really have to do anything
	// since the logic is going to delete the slice anyway
	if err = c.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//remove the time slice since we know the server has the most up to date
	// copy
	delete(c.lookupTimeSlices, timeSlice.UUID)
	//update the timer, remove its time slice and update the elapsed time
	//set the activeSliceUUID to empty
	timer.ActiveSliceUUID = ""
	//calculate elapsed time
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	//update the timer
	if c.remote != nil {
		if err = c.remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//update the timer
	err = c.meta.MetaTimerWrite(timer.UUID, timer)

	return
}

//submit a timer
func (c *client) TimerSubmit(timerID string, submitTime time.Time) (err error) {
	c.Lock()
	defer c.Unlock()

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submitted" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)

	var timer bludgeon.Timer
	var timeSlice bludgeon.TimeSlice

	//attempt to use the api to get the provided timer
	if timer, err = c.timerRead(timerID); err != nil {
		return
	}
	//check if there's an active slice, if there is,
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
		if c.remote != nil {
			if err = c.remote.TimeSliceUpdate(timeSlice); err != nil {
				//TODO: cache the operation
			}
			//remove the time slice since we know the server has the most up to date
			// copy
			delete(c.lookupTimeSlices, timeSlice.UUID)
		}
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
	if c.remote != nil {
		//update the timer
		if err = c.remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//update the timer in meta
	err = c.meta.MetaTimerWrite(timer.UUID, timer)

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
	if c.remote != nil {
		if err = c.remote.TimerUpdate(timer); err != nil {
			//TODOcache the operation
		}
	}
	err = c.meta.MetaTimerWrite(timerID, timer)

	return
}
