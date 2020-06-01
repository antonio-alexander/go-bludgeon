package bludgeon

import (
	"errors"
	"fmt"
	"time"
)

//The goal of this is to create a varratic layer for any of the functions with common logic being shared
// between the server and client, the goal is to prevent having duplicate logic in both the server and the
// client to facilitate having a client with only a remote and not a meta (for example)

//REVIEW: How should we handle/output remote and meta errors
//REVIEW: there's a huge problem with time synchronization between remote and meta pointers

//sortMetaRemote the goal of this function is to sort varadics and output
// meta and remote (to simplify the api)
func sortMetaRemote(i ...interface{}) (meta Meta, remote Remote, err error) {
	//check if an appropriate number of varadics
	if len(i) <= 0 || len(i) > 2 {
		err = errors.New("Too few or too many varadics")

		return
	}
	//switch over variables and
	for _, i := range i {
		//switch on the interface type
		switch v := i.(type) {
		case Remote:
			remote = v
		case Meta:
			meta = v
		default:
			err = fmt.Errorf("unsupported varatic: %t", v)

			return
		}
	}

	return
}

//timerRead will use meta and remote to read a single timer from remote/meta while
// prioritizing meta and falling back to remote if meta fails or in addition to
func timerRead(timerID string, meta Meta, remote Remote) (timer Timer, err error) {
	//check if remote is nil
	if remote != nil {
		//attempt to use the api to get the provided timer
		if timer, err = remote.TimerRead(timerID); err != nil {
			//TODO: store the remote error
			//TODO: cache the read?
		} else {
			//if meta is not nil, update meta with the newly read timer
			if meta != nil {
				//store the timer in the local cache
				if err = meta.MetaTimerWrite(timerID, timer); err != nil {
					return
				}
			}
		}
	}
	//read the meta timer if remote is nil or there is no err
	if err != nil || remote == nil {
		//attempt to read the timer locally
		if timer, err = meta.MetaTimerRead(timerID); err != nil {
			return
		}
	}

	return
}

//LaunchCache
func LaunchCache(stopper <-chan struct{}, wg interface {
	Add(int)
	Done()
}, chCache <-chan CacheData) {
	//create channel to block until go routine enters business logic
	started := make(chan struct{})
	//launch go Routine
	wg.Add(1)
	go func() {
		defer wg.Done()

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
			case <-chCache:
			case <-stopper:
				return
			}
		}
	}()
	//wait for goRoutine to enter business logic
	<-started
}

//TimeSliceCreate
func TimeSliceCreate(timerUUID string, i ...interface{}) (timeSlice TimeSlice, err error) {
	var meta Meta
	var remote Remote

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//use the api to create a time slice if remote is not nil
	if remote != nil {
		if timeSlice, err = remote.TimeSliceCreate(timerUUID); err != nil {
			//TODO: cache the operation
		}
	}
	//REVIEW: what if the first remote call succeeds, but the second one fails? should probably update the api to include
	//only generate the idea if the remote operations failed
	if err != nil || remote == nil {
		//generate the time slice id
		if timeSlice.UUID, err = GenerateID(); err != nil {
			return
		}
		//update the time slice's timer ID
		timeSlice.TimerUUID = timerUUID
	}
	//store the time slice of meta is not nil
	if meta != nil {
		//store the timeslice
		err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)
	}

	return
}

//TimeSliceRead
func TimeSliceRead(timeSliceID string, i ...interface{}) (timeSlice TimeSlice, err error) {
	var meta Meta
	var remote Remote

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	if remote != nil {
		//attempt to use the api to query the timeslice
		if timeSlice, err = remote.TimeSliceRead(timeSliceID); err != nil {
			//TODO: cache attempt to read time slice?
		} else {
			if meta != nil {
				//write time slice to meta
				if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
					return
				}
			}
		}
	}
	//only read from meta if remote is nil or there's an error
	if remote == nil || err != nil {
		//attempt to read the timeSlice locally
		if timeSlice, err = meta.MetaTimeSliceRead(timeSliceID); err != nil {
			return
		}
	}

	return
}

//timeSliceUpdate
func TimeSliceUpdate(timeSlice TimeSlice, i ...interface{}) (err error) {
	var meta Meta
	var remote Remote

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//use the api to create a time slice if remote is not nil
	if remote != nil {
		if err = remote.TimeSliceUpdate(timeSlice); err != nil {
			//TODO: cache the operation
		}
	}
	if meta != nil {
		//store the timeslice
		err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)
	}

	return
}

//TimeSliceDelete
func TimeSliceDelete(timeSliceID string, i ...interface{}) (err error) {
	var meta Meta
	var remote Remote

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//use the api to create a time slice if remote is not nil
	if remote != nil {
		if err = remote.TimeSliceDelete(timeSliceID); err != nil {
			//TODO: cache the operation
		}
	}
	if meta != nil {
		//store the timeslice
		err = meta.MetaTimeSliceDelete(timeSliceID)
	}

	return
}

//TimerCreate
func TimerCreate(i ...interface{}) (timer Timer, err error) {
	var remote Remote
	var meta Meta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//check if remote is nil
	if remote != nil {
		//attempt execute the api
		if timer, err = remote.TimerCreate(); err != nil {
			//TODO: cache the operation
		}
	}
	//generate the uuid if there's a remote error or if there is no remote
	if meta != nil {
		//generate a UUID for the timer locally
		if timer.UUID, err = GenerateID(); err != nil {
			return
		}
		//update the timer in meta
		if err = meta.MetaTimerWrite(timer.UUID, timer); err != nil {
			return
		}
	}

	return
}

//get a timer
func TimerRead(id string, i ...interface{}) (timer Timer, err error) {
	var remote Remote
	var meta Meta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//attempt to read the timer
	if timer, err = timerRead(id, meta, remote); err != nil {
		return
	}
	//attempt to get the activeSlice if it exists
	if timer.ActiveSliceUUID != "" {
		//functionally, this is only useful if meta is not nil (as a kind of cache)
		if meta != nil {
			//attempt to read the slice
			if _, err = TimeSliceRead(timer.ActiveSliceUUID, meta, remote); err != nil {
				return
			}
		}
	}

	return
}

//TimerUpdate
func TimerUpdate(timer Timer, i ...interface{}) (err error) {
	var remote Remote
	var meta Meta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//use the api to create a time slice if remote is not nil
	if remote != nil {
		if err = remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	if meta != nil {
		//store the timeslice
		err = meta.MetaTimerWrite(timer.UUID, timer)
	}

	return
}

//delete a timer
func TimerDelete(timerID string, i ...interface{}) (err error) {
	var timer Timer
	var remote Remote
	var meta Meta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//attempt to use the api to get the provided timer
	if timer, err = timerRead(timerID, meta, remote); err != nil {
		return
	}
	//check if there's an active slice
	if timer.ActiveSliceUUID != "" {
		//TODO: delete all associated slices??
	}
	timer.Archived = true
	//update the timer
	if remote != nil {
		if err = remote.TimerUpdate(timer); err != nil {
			//TODOcache the operation
		}
	}
	if meta != nil {
		err = meta.MetaTimerWrite(timerID, timer)
	}
	return
}

//start a timer
func TimerStart(id string, startTime time.Time, i ...interface{}) (err error) {
	var timer Timer
	var timeSlice TimeSlice
	var remote Remote
	var meta Meta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//attempt to read the timer
	if timer, err = timerRead(id, meta, remote); err != nil {
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
		if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta, remote); err != nil {
			return
		}
	} else {
		//since there is no active slice, create the timeSlice
		if timeSlice, err = TimeSliceCreate(timer.UUID, meta, remote); err != nil {
			return
		}
		//set the start time
		timeSlice.Start = startTime.UnixNano()
		//update the time slice with its new start time
		if remote != nil {
			//update the slice using the API
			if err = remote.TimeSliceUpdate(timeSlice); err != nil {
				//TODO: cache operation
			}
		}
	}
	//store the slice
	if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//store the timeSlice
	timer.ActiveSliceUUID = timeSlice.UUID
	if remote != nil {
		//update the timer
		if err = remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//cache the timer
	err = meta.MetaTimerWrite(timer.UUID, timer)

	return
}

//pause a timer
func TimerPause(timerID string, pauseTime time.Time, i ...interface{}) (err error) {
	var remote Remote
	var meta Meta
	var timer Timer
	var timeSlice TimeSlice

	//if the given timer exists, when we pause the timer, we will
	// grab the active time slice, set the finish time to now
	// set the active slice value to -1 update the elapsed time,
	// add all the timers that aren't archived together

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//attempt to read the timer
	if timer, err = timerRead(timerID, meta, remote); err != nil {
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
	if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta, remote); err != nil {
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
	if remote != nil {
		//update the time slice
		if err = remote.TimeSliceUpdate(timeSlice); err != nil {
			//TODO: cache the operation
		}
	}
	// because the service side logic should complete, we don't really have to do anything
	// since the logic is going to delete the slice anyway
	if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
		return
	}
	//update the timer, remove its time slice and update the elapsed time
	//set the activeSliceUUID to empty
	timer.ActiveSliceUUID = ""
	//calculate elapsed time
	timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
	//update the timer
	if remote != nil {
		if err = remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//update the timer
	err = meta.MetaTimerWrite(timer.UUID, timer)

	return
}

//submit a timer
func TimerSubmit(timerID string, submitTime time.Time, i ...interface{}) (err error) {
	var timer Timer
	var timeSlice TimeSlice
	var remote Remote
	var meta Meta

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submitted" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//attempt to use the api to get the provided timer
	if timer, err = timerRead(timerID, meta, remote); err != nil {
		return
	}
	//check if there's an active slice, if there is,
	if timer.ActiveSliceUUID != "" {
		//read the active time slice
		if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta, remote); err != nil {
			return
		}
		//set the finish time
		timeSlice.Finish = submitTime.UnixNano()
		//calculate the elapsed time
		timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
		//set archived to true
		timeSlice.Archived = true
		//update the slice using the API, there's an understanding that this should
		if remote != nil {
			if err = remote.TimeSliceUpdate(timeSlice); err != nil {
				//TODO: cache the operation
			}
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
	if remote != nil {
		//update the timer
		if err = remote.TimerUpdate(timer); err != nil {
			//TODO: cache the operation
		}
	}
	//update the timer in meta
	err = meta.MetaTimerWrite(timer.UUID, timer)

	return
}
