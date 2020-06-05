package bludgeonclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

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
	stopper chan struct{} //stopper to stop goRoutines
	chCache chan bludgeon.CacheData
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
} {
	//ensure that either meta or remote is not nil
	if meta == nil && remote == nil {
		panic("meta/remote pointers nil")
	}
	//populate client pointer
	return &client{
		meta:   meta,
		remote: remote,
	}
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
	c.meta, c.remote = nil, nil

	return
}

//Serialize
func (c *client) Serialize() (bytes []byte, err error) {
	c.RLock()
	defer c.RUnlock()

	//serialize data into bytes
	bytes, err = json.Marshal(&SerializedData{
		//
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
	//launch cache goRoutine
	bludgeon.LaunchCache(c.stopper, c, c.chCache)
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
	//CommandHandler
	CommandHandler(command bludgeon.CommandClient, dataIn interface{}) (dataOut interface{}, err error)
}

//ensure that cache implements Functional
var _ Functional = &client{}

func (c *client) CommandHandler(command bludgeon.CommandClient, dataIn interface{}) (dataOut interface{}, err error) {
	var id string
	var ok bool

	//execute the command
	switch command {
	case bludgeon.CommandClientTimerCreate:
		dataOut, err = c.TimerCreate()
	case bludgeon.CommandClientTimerStart:
		if id, ok = dataIn.(string); !ok {
			err = errors.New("Unable to cast into string")
		} else {
			//start the timer
			dataOut, err = c.TimerStart(id, time.Now())
		}
	case bludgeon.CommandClientTimerPause, bludgeon.CommandClientTimerStop:
		if id, ok = dataIn.(string); !ok {
			err = errors.New("Unable to cast into string")
		} else {
			//start the timer
			dataOut, err = c.TimerPause(id, time.Now())
		}
	case bludgeon.CommandClientTimerRead:
		//cast into string and use
		if id, ok = dataIn.(string); !ok {
			err = errors.New("Unable to cast into string")
		} else {
			//read the timer
			dataOut, err = c.TimerRead(id)
		}
	case bludgeon.CommandClientTimerSubmit:
		if id, ok = dataIn.(string); !ok {
			err = errors.New("Unable to cast into string")
		} else {
			//start the timer
			dataOut, err = c.TimerSubmit(id, time.Now())
		}
	default:
		err = fmt.Errorf("Unsupported command: %s", command)
	}

	return
}

//add a timer
func (c *client) TimerCreate() (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//create a timer
	timer, err = bludgeon.TimerCreate(c.meta, c.remote)

	return
}

//get a timer
func (c *client) TimerRead(id string) (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//read the timer
	timer, err = bludgeon.TimerRead(id, c.meta, c.remote)

	return
}

//start a timer
func (c *client) TimerStart(id string, startTime time.Time) (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//start the timer
	timer, err = bludgeon.TimerStart(id, startTime, c.meta, c.remote)

	return
}

//pause a timer
func (c *client) TimerPause(timerID string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//pause the timer
	timer, err = bludgeon.TimerPause(timerID, pauseTime, c.meta, c.remote)

	return
}

//submit a timer
func (c *client) TimerSubmit(timerID string, submitTime time.Time) (timer bludgeon.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//submit the timer
	timer, err = bludgeon.TimerSubmit(timerID, submitTime, c.meta, c.remote)

	return
}

//delete a timer
func (c *client) TimerDelete(timerID string) (err error) {
	c.Lock()
	defer c.Unlock()

	//delete the timer
	err = bludgeon.TimerDelete(timerID, c.meta, c.remote)

	return
}
