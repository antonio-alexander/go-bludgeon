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
//REVIEW: generally cached operations are done after the remote fails, maybe we should do it
// only if there's some kind of meta? especially in cases where a guid is created locally rather
// than remotely, if we cache right after the failure, we'd need to ensure its guid matched what
// was present locally

type lazyMeta interface {
	MetaTimer
	MetaTimeSlice
}

type lazyRemote interface {
	//RemoteTimer
	RemoteTimer

	//RemoteTimeSlice
	RemoteTimeSlice
}

//sortMetaRemote the goal of this function is to sort varadics and output
// meta and remote (to simplify the api)
func sortMetaRemote(i []interface{}) (meta lazyMeta, remote lazyRemote, err error) {
	//check if an appropriate number of varadics
	if len(i) <= 0 || len(i) > 2 {
		err = errors.New("Too few or too many varadics")

		return
	}
	//switch over variables and
	for _, i := range i {
		if i != nil {
			//switch on the interface type
			switch v := i.(type) {
			case lazyRemote:
				remote = v
			case lazyMeta:
				meta = v
			default:
				err = fmt.Errorf("unsupported varatic: %t", v)

				return
			}
		}
	}

	return
}

//timerRead will use meta and remote to read a single timer from remote/meta while
// prioritizing meta and falling back to remote if meta fails or in addition to
func timerRead(timerID string, meta lazyMeta, remote lazyRemote) (timer Timer, err error) {
	//check if remote is nil
	if remote != nil {
		//attempt to use the api to get the provided timer
		if timer, err = remote.TimerRead(timerID); err != nil {
			//TODO: store the remote error
		} else {
			//if meta is not nil, update meta with the newly read timer
			if meta != nil {
				//store the timer in the local cache
				if err = meta.MetaTimerWrite(timerID, timer); err != nil {
					return
				}
			}
			//we don't want to continue with the remainder of the logic since it pre-supposes
			// that the remote is nil and meta is non-nil
			return
		}
	}
	//read the meta timer if remote is nil or there is an error
	if (err != nil || remote == nil) && meta != nil {
		//attempt to read the timer locally
		if timer, err = meta.MetaTimerRead(timerID); err != nil {
			//TODO: store meta error
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
func timeSliceCreate(timerUUID string, i ...interface{}) (timeSlice TimeSlice, err error) {
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, _, err = sortMetaRemote(i); err != nil {
		return
	}
	//
	if meta != nil {
		//generate the time slice id
		if timeSlice.UUID, err = GenerateID(); err != nil {
			return
		}
		//update the time slice's timer ID
		timeSlice.TimerUUID = timerUUID
		//store the timeslice
		if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
			//TODO: store the meta error
			return
		}
	}

	return
}

//timeSliceUpdate
func timeSliceUpdate(timeSlice TimeSlice, i ...interface{}) (err error) {
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, _, err = sortMetaRemote(i); err != nil {
		return
	}
	//store the timeslice if meta is not nil
	if meta != nil {
		if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
			//TODO: store the meta error
			return
		}
	}

	return
}

//TimeSliceDelete
func timeSliceDelete(timeSliceID string, i ...interface{}) (err error) {
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, _, err = sortMetaRemote(i); err != nil {
		return
	}
	if meta != nil {
		//delete the timeslice
		err = meta.MetaTimeSliceDelete(timeSliceID)
	}

	return
}

//TimeSliceRead
func TimeSliceRead(timeSliceID string, i ...interface{}) (timeSlice TimeSlice, err error) {
	var meta lazyMeta
	var remote lazyRemote

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//
	if remote != nil {
		//attempt to read time slice
		if timeSlice, err = remote.TimeSliceRead(timeSliceID); err != nil {
			//TODO: store remote error
		} else {
			//store timeSlice in meta
			if meta != nil {
				if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
					//TODO: store meta error

					return
				}
			}
			//return, there's no reason to attempt to read the timeslice
			return
		}
	}
	//only query meta if remote fails, is nil and meta is not nil
	if (err != nil || remote == nil) && meta != nil {
		if timeSlice, err = meta.MetaTimeSliceRead(timeSliceID); err != nil {
			//TODO: store meta eror
			return
		}
	}

	return
}

//TimerCreate
func TimerCreate(i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//check if remote is nil
	if remote != nil {
		//attempt execute the api
		if timer, err = remote.TimerCreate(); err != nil {
			//TODO: cache the operation
		} else {
			if meta != nil {
				//update the timer in meta
				if err = meta.MetaTimerWrite(timer.UUID, timer); err != nil {
					//TODO: cache meta error
					return
				}
			}
			//return since remote succeeded
			return
		}
	}
	//generate the uuid if there's a remote error or if there is no remote
	if (err != nil || remote == nil) && meta != nil {
		//generate a UUID for the timer locally
		if timer.UUID, err = GenerateID(); err != nil {
			return
		}
		//update the timer in meta
		if err = meta.MetaTimerWrite(timer.UUID, timer); err != nil {
			//TODO: store meta error
			return
		}
	}

	return
}

//get a timer
func TimerRead(id string, i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta
	var timeSlice TimeSlice

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
		//read the time slice to store in meta if remote
		if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta, remote); err != nil {
			return
		}
		//update elapsed time in realtime if there is an active timeslice
		timer.ElapsedTime += (time.Now().UnixNano() - timeSlice.Start)
	}

	return
}

//TimerUpdate
func TimerUpdate(t Timer, i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//use the api to create a time slice if remote is not nil
	if remote != nil {
		if timer, err = remote.TimerUpdate(t); err != nil {
			//TODO: store remote error
			//TODO: cache the operation
		} else {
			if meta != nil {
				if err = meta.MetaTimerWrite(timer.UUID, timer); err != nil {
					//TODO: store meta error

					return
				}
			}

			return
		}
	}
	if meta != nil {
		//get the current timer
		if timer, err = timerRead(t.UUID, meta, remote); err != nil {
			return
		}
		//update only the items that have changed
		// update the comment
		if t.Comment != "" {
			timer.Comment = t.Comment
		}
		//store the timeslice
		err = meta.MetaTimerWrite(timer.UUID, timer)
	}

	return
}

//delete a timer
func TimerDelete(timerID string, i ...interface{}) (err error) {
	var remote lazyRemote
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//if remote is non nil, delete the timer
	if remote != nil {
		if err = remote.TimerDelete(timerID); err != nil {
			//TODO: store remote error
			//TODOcache the operation
		}
	}
	//if meta is non nil, delete the timer
	if meta != nil {
		if err = meta.MetaTimerDelete(timerID); err != nil {
			//TODO: store meta error
			return
		}
	}

	return
}

//TimerStart will
func TimerStart(id string, startTime time.Time, i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//check if remote is not nil
	if remote != nil {
		//TODO: should this output the updated timer
		if timer, err = remote.TimerStart(id, startTime); err != nil {
			//TODO: cache operation
			//TODO: store remote error
		} else {
			if meta != nil {
				//store updated timer in meta
				if err = meta.MetaTimerWrite(id, timer); err != nil {
					//TODO: store meta timer
					return
				}
			}
			//return since the remote was successful
			return
		}
	}
	//only perform the below if the remote fails or is nil
	if (err != nil || remote == nil) && meta != nil {
		var timeSlice TimeSlice

		//attempt to read the timer
		if timer, err = timerRead(id, meta, nil); err != nil {
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
			if timeSlice, err = timeSliceCreate(timer.UUID, meta, remote); err != nil {
				return
			}
			//set the start time
			timeSlice.Start = startTime.UnixNano()
		}
		//store the timeSlice
		timer.ActiveSliceUUID = timeSlice.UUID
		//store the time slice
		if err = meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice); err != nil {
			//TODO: store meta error
			return
		}
		//store the timer
		if err = meta.MetaTimerWrite(timer.UUID, timer); err != nil {
			//TODO: store meta error
			return
		}
		//REVIEW: if we create a GUID, we will need some way to overwrite the cached operation
	}

	return
}

//pause a timer
func TimerPause(timerID string, pauseTime time.Time, i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta

	//if the given timer exists, when we pause the timer, we will
	// grab the active time slice, set the finish time to now
	// set the active slice value to -1 update the elapsed time,
	// add all the timers that aren't archived together

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//
	if remote != nil {
		if timer, err = remote.TimerPause(timerID, pauseTime); err != nil {
			//TODO: cache remote error
			//TODO: cache operation
		} else {
			//REVIEW: should we also store the slice??
			if meta != nil {
				//store updated timer
				if err = meta.MetaTimerWrite(timerID, timer); err != nil {
					//TODO: store meta error
					return
				}
			}
			//
			return
		}
	}
	//
	if (err != nil || remote == nil) && meta != nil {
		var timeSlice TimeSlice

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
		if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta); err != nil {
			return
		}
		//set the finish time
		timeSlice.Finish = pauseTime.UnixNano()
		//calculate the elapsed time
		timeSlice.ElapsedTime = timeSlice.Finish - timeSlice.Start
		//set archived to true
		timeSlice.Archived = true
		//update the timer, remove its time slice and update the elapsed time
		//set the activeSliceUUID to empty
		timer.ActiveSliceUUID = ""
		//calculate elapsed time
		timer.ElapsedTime = timer.ElapsedTime + timeSlice.ElapsedTime
		//update the timeslice
		if err = timeSliceUpdate(timeSlice, meta); err != nil {
			return
		}
		//update the timer
		if _, err = TimerUpdate(timer, meta); err != nil {
			return
		}
	}

	return
}

//submit a timer
func TimerSubmit(timerID string, submitTime time.Time, i ...interface{}) (timer Timer, err error) {
	var remote lazyRemote
	var meta lazyMeta

	//when the timer is submitted, the stop time is updated, the active
	// time slice is completed with the current time and the timer is
	// set to "submitted" so any changes to it after the fact are known as
	// changes that shouldn't involve the time slices (and that they're now
	// invalid)

	//sort the varadics into meta/remote
	if meta, remote, err = sortMetaRemote(i); err != nil {
		return
	}
	//
	if remote != nil {
		if timer, err = remote.TimerSubmit(timerID, submitTime); err != nil {
			//TODO: cache operation
			//TODO: store remote error
		} else {
			//REVIEW: should we also store the active time slice?
			if meta != nil {
				//TODO: store finished time
				if err = meta.MetaTimerWrite(timerID, timer); err != nil {
					//TODO: store meta error
					return
				}
			}
			//
			return
		}
	}
	//
	if (err != nil || remote == nil) && meta != nil {
		var timeSlice TimeSlice

		//attempt to use the api to get the provided timer
		if timer, err = timerRead(timerID, meta, remote); err != nil {
			return
		}
		//check if there's an active slice, if there is,
		if timer.ActiveSliceUUID != "" {
			//read the active time slice
			if timeSlice, err = TimeSliceRead(timer.ActiveSliceUUID, meta); err != nil {
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
		if err = timeSliceUpdate(timeSlice, meta); err != nil {
			return
		}
		//update the timer
		if _, err = TimerUpdate(timer, meta); err != nil {
			return
		}
	}

	return
}
