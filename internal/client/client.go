package bludgeonclient

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	config "github.com/antonio-alexander/go-bludgeon/bludgeon/client/config"
	common "github.com/antonio-alexander/go-bludgeon/internal/common"
)

type client struct {
	sync.RWMutex                         //mutex for threadsafe operations
	sync.WaitGroup                       //waitgroup to track go routines
	started        bool                  //whether or not the business logic has starting
	config         config.Configuration  //configuration
	stopper        chan struct{}         //stopper to stop goRoutines
	chExternal     chan struct{}         //channel for owners to block on
	chCache        chan common.CacheData //channel for cache data
	log            *log.Logger           //logger for non error messages
	logError       *log.Logger           //logger for errors
	meta           interface {           //storage interface
		common.MetaTimer
	}
	remote interface { //remote interface
		common.FunctionalTimer
		common.FunctionalTimeSlice
	}
}

func NewClient(log, logError *log.Logger, meta interface {
	common.MetaTimer
}, remote interface {
	common.FunctionalTimer
	common.FunctionalTimeSlice
}) interface {
	Owner
	Manage
	common.FunctionalTimer
} {
	//ensure that either meta or remote is not nil
	if meta == nil && remote == nil {
		panic("meta/remote pointers nil")
	}
	//populate client pointer
	return &client{
		meta:     meta,
		remote:   remote,
		log:      log,
		logError: logError,
	}
}

func (c *client) Println(v ...interface{}) {
	if c.log != nil {
		c.log.Println(v...)
	}
}

func (c *client) Printf(format string, v ...interface{}) {
	if c.log != nil {
		c.log.Printf(format, v...)
	}
}

func (c *client) Print(v ...interface{}) {
	if c.log != nil {
		c.log.Print(v...)
	}
}

func (c *client) Error(err error) {
	if c.logError != nil {
		c.logError.Println(err)
	}
}

func (c *client) Errorf(format string, v ...interface{}) {
	if c.logError != nil {
		c.logError.Printf(format, v...)
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
	c.config = config.Configuration{}
	//set internal pointers to nil
	c.meta, c.remote = nil, nil
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
	Start(config config.Configuration) (chExternal <-chan struct{}, err error)

	//
	Stop() (err error)
}

func (c *client) Start(conf config.Configuration) (chExternal <-chan struct{}, err error) {
	c.Lock()
	defer c.Unlock()

	//check if already started
	if c.started {
		err = errors.New(ErrStarted)

		return
	}
	//store configuration
	c.config = conf
	//create stopper
	c.stopper = make(chan struct{})
	//create channel external
	c.chExternal = make(chan struct{})
	chExternal = c.chExternal
	//launch cache goRoutine
	common.LaunchCache(c.stopper, c, c.chCache)
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
	close(c.chExternal)
	//set started flag to false
	c.started = false

	return
}

//ensure that cache implements Functional
var _ common.FunctionalTimer = &client{}

//add a timer
func (c *client) TimerCreate() (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//create a timer
	timer, err = common.TimerCreate(c.meta, c.remote)

	return
}

//get a timer
func (c *client) TimerRead(id string) (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//read the timer
	timer, err = common.TimerRead(id, c.meta, c.remote)

	return
}

//add a timer
func (c *client) TimerUpdate(t common.Timer) (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//create a timer
	timer, err = common.TimerUpdate(t, c.meta, c.remote)

	return
}

//start a timer
func (c *client) TimerStart(id string, startTime time.Time) (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//start the timer
	timer, err = common.TimerStart(id, startTime, c.meta, c.remote)

	return
}

//pause a timer
func (c *client) TimerPause(timerID string, pauseTime time.Time) (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//pause the timer
	timer, err = common.TimerPause(timerID, pauseTime, c.meta, c.remote)

	return
}

//submit a timer
func (c *client) TimerSubmit(timerID string, submitTime time.Time) (timer common.Timer, err error) {
	c.Lock()
	defer c.Unlock()

	//submit the timer
	timer, err = common.TimerSubmit(timerID, submitTime, c.meta, c.remote)

	return
}

//delete a timer
func (c *client) TimerDelete(timerID string) (err error) {
	c.Lock()
	defer c.Unlock()

	//delete the timer
	err = common.TimerDelete(timerID, c.meta, c.remote)

	return
}
