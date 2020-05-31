package bludgeonserver

import (
	"errors"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

//REVIEW: the future version of the client, will have a remote and meta pointer provided
// at new, at least one must not be nil, but if only remote, it'll use the server for all
// operations, if remote is nil but meta is not, it'll use meta instead for local persistence

type server struct {
	sync.RWMutex                 //mutex for threadsafe operations
	sync.WaitGroup               //waitgroup to track go routines
	started        bool          //whether or not the business logic has starting
	config         Configuration //configuration
	meta           interface {   //storage interface
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
	}
	stopper chan struct{} //stopper to stop goRoutines
	//some set of flags or method to know what to update
	tokens map[string]bludgeon.Token //map of tokens
}

func NewServer(meta interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}) interface {
	Owner
	Manage
	Functional
} {
	//REVIEW: should add cases to confirm that meta/remote aren't nil since
	// basic functionality won't work?
	if meta == nil {
		panic("meta is nil")
	}
	//create internal maps
	//TODO: need a way to prepopulate or generate the
	// lookups from scratch? will need to have the full data some how
	//populate client pointer
	return &server{
		meta:   meta,
		tokens: make(map[string]bludgeon.Token),
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

func (s *server) Close() {
	s.Lock()
	defer s.Unlock()

	//use the remote and/or storage functions to:
	// attempt to synchronize the current values
	// attempt to serialize what's remaining

	//set internal configuration to default
	s.config = Configuration{}
	//set internal pointers to nil
	s.meta = nil

	return
}

//Serialize
func (s *server) Serialize() (bytes []byte, err error) {
	s.RLock()
	defer s.RUnlock()

	return
}

//Deserialize
func (s *server) Deserialize(bytes []byte) (err error) {
	s.Lock()
	defer s.Unlock()

	return
}

type Manage interface {
	//
	Start(config Configuration) (err error)

	//
	Stop() (err error)
}

func (s *server) Start(config Configuration) (err error) {
	s.Lock()
	defer s.Unlock()

	//check if already started
	if s.started {
		err = errors.New(ErrStarted)

		return
	}
	//store configuration
	s.config = config
	//create stopper
	s.stopper = make(chan struct{})
	//launch go routines
	s.LaunchPurge()
	//set started to true
	s.started = true

	return
}

func (s *server) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	//check if not started
	if !s.started {
		err = errors.New(ErrNotStarted)

		return
	}
	//close the stopper
	close(s.stopper)
	s.Wait()
	//set started flag to false
	s.started = false

	return
}

type Functional interface {
	//
	LaunchPurge()

	//CommandHandler
	CommandHandler(command bludgeon.CommandServer, dataIn interface{}) (dataOut interface{}, err error)

	//CommandHandler
	CommandHandlerToken(command bludgeon.CommandServer, dataIn interface{}, token bludgeon.Token) (dataOut interface{}, err error)
}

//ensure that cache implements Functional
var _ Functional = &server{}

func (s *server) LaunchPurge() {
	started := make(chan struct{})
	s.Add(1)
	go s.goPurge(started)
	<-started
}

func (s *server) goPurge(started chan struct{}) {
	defer s.Done()

	tPurge := time.NewTicker(30 * time.Second)
	close(started)
	for {
		select {
		case <-tPurge.C:
			// if err := s.TokenPurge(); err != nil {
			// 	fmt.Printf(err)
			// }
		case <-s.stopper:
			return
		}
	}
}

func (s *server) CommandHandler(command bludgeon.CommandServer, dataIn interface{}) (dataOut interface{}, err error) {
	//execute the command
	switch command {
	case bludgeon.CommandServerTimerCreate:
		dataOut, err = s.TimerCreate()
	case bludgeon.CommandServerTimerRead:
		if id, ok := dataIn.(string); !ok {
			//TODO: generate error
		} else {
			//read time slice
			dataOut, err = s.TimerRead(id)
		}
	case bludgeon.CommandServerTimerUpdate:
		if timer, ok := dataIn.(bludgeon.Timer); !ok {
			//TODO: generate error
		} else {
			err = s.TimerUpdate(timer)
		}
	case bludgeon.CommandServerTimerDelete:
		if id, ok := dataIn.(string); !ok {
			//TODO: generate error
		} else {
			//read time slice
			err = s.TimerDelete(id)
		}
	case bludgeon.CommandServerTimeSliceCreate:
		if id, ok := dataIn.(string); !ok {
			//TODO: generate error
		} else {
			dataOut, err = s.TimeSliceCreate(id)
		}
	case bludgeon.CommandServerTimeSliceRead:
		if id, ok := dataIn.(string); !ok {
			//TODO: generate error
		} else {
			//read time slice
			dataOut, err = s.TimeSliceRead(id)
		}
	case bludgeon.CommandServerTimeSliceUpdate:
		if timeSlice, ok := dataIn.(bludgeon.TimeSlice); !ok {
			//TODO: generate error
		} else {
			err = s.TimeSliceUpdate(timeSlice)
		}
	case bludgeon.CommandServerTimeSliceDelete:
		if id, ok := dataIn.(string); !ok {
			//TODO: generate error
		} else {
			err = s.TimeSliceDelete(id)
		}
	default:
		//TODO: generate error
	}

	return
}

func (s *server) CommandHandlerToken(command bludgeon.CommandServer, dataIn interface{}, token bludgeon.Token) (dataOut interface{}, err error) {
	//TODO: check token
	//execute command
	dataOut, err = s.CommandHandler(command, dataIn)

	return
}

//
func (s *server) TimerCreate() (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//generate timer uuid
	if timer.UUID, err = bludgeon.GenerateID(); err != nil {
		return
	}
	//write timer
	err = s.meta.MetaTimerWrite(timer.UUID, timer)

	return
}

//
func (s *server) TimerRead(id string) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//read the timer
	timer, err = s.meta.MetaTimerRead(id)

	return
}

//
func (s *server) TimerUpdate(timer bludgeon.Timer) (err error) {
	s.Lock()
	defer s.Unlock()

	//update the timer
	err = s.meta.MetaTimerWrite(timer.UUID, timer)

	return
}

//
func (s *server) TimerDelete(id string) (err error) {
	s.Lock()
	defer s.Unlock()

	//delete the timer
	err = s.meta.MetaTimerDelete(id)

	return
}

//
func (s *server) TimeSliceCreate(timerID string) (timeSlice bludgeon.TimeSlice, err error) {
	s.Lock()
	defer s.Unlock()

	//generate time slice uuid
	if timeSlice.UUID, err = bludgeon.GenerateID(); err != nil {
		return
	}
	timeSlice.TimerUUID = timerID
	//write timer
	err = s.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)

	return
}

//
func (s *server) TimeSliceRead(id string) (timeSlice bludgeon.TimeSlice, err error) {
	s.Lock()
	defer s.Unlock()

	//read the time slice
	timeSlice, err = s.meta.MetaTimeSliceRead(id)

	return
}

//
func (s *server) TimeSliceUpdate(timeSlice bludgeon.TimeSlice) (err error) {
	s.Lock()
	defer s.Unlock()

	//update the timeslice
	err = s.meta.MetaTimeSliceWrite(timeSlice.UUID, timeSlice)

	return
}

//
func (s *server) TimeSliceDelete(id string) (err error) {
	s.Lock()
	defer s.Unlock()

	//delete the time slice
	err = s.meta.MetaTimeSliceDelete((id))

	return
}

type Token interface {
	//Generate creates a token, outputs a token
	Generate() string

	//Verify will check if a token exists
	Verify(tokenID string) bool

	//Purge deletes tokens that are as old as wait
	Purge()

	//Delete will remove a token
	Delete(token string) (deleted bool)
}

// //Generate creates a token, outputs a token
// func (t *token) Generate() string {
// 	d.Lock()
// 	defer d.Unlock()

// 	b := make([]byte, 8)
// 	var t Token

// 	for {
// 		rand.Read(b) //generate token
// 		token := fmt.Sprintf("%x-%x-%x-%x", b[0:2], b[2:4], b[4:6], b[6:8])

// 		if _, ok := d.Tokens[token]; !ok {
// 			t.Token = token
// 			t.Time = time.Now().UnixNano()
// 			d.Tokens[token] = t

// 			return t.Token
// 		}
// 	}
// }

// //Verify will check if a token exists
// func (d *DataToken) Verify(tokenID string) bool {
// 	d.RLock()
// 	defer d.RUnlock()

// 	var valid bool

// 	if token, ok := d.Tokens[tokenID]; ok {
// 		if (time.Now().UnixNano() - token.Time) > int64(d.Wait/time.Nanosecond) {
// 			delete(d.Tokens, tokenID)
// 		} else {
// 			valid = true
// 		}
// 	}
// 	return valid
// }

// //Purge deletes tokens that are as old as wait
// func (d *DataToken) Purge() {
// 	d.Lock()
// 	defer d.Unlock()

// 	time :=
// 	for key, value := range d.Tokens {
// 		if time.Now().UnixNano()  - token.Time > d.Wait {
// 			delete(d.Tokens, key)
// 		}
// 	}
// }

// //Delete will remove a token
// func (d *DataToken) Delete(token string) (deleted bool) {
// 	d.Lock()
// 	defer d.Unlock()

// 	if _, ok := d.Tokens[token]; ok {
// 		delete(d.Tokens, token)
// 		deleted = true
// 	}

// 	return
// }
