package bludgeonserverfunctional

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	common "github.com/antonio-alexander/go-bludgeon/bludgeon/server/common"
)

type server struct {
	sync.RWMutex                        //mutex for threadsafe operations
	sync.WaitGroup                      //waitgroup to track go routines
	started        bool                 //whether or not the business logic has starting
	config         common.Configuration //configuration
	stopper        chan struct{}        //stopper to stop goRoutines
	chExternal     chan struct{}
	tokens         map[string]bludgeon.Token //map of tokens
	logError       *log.Logger
	log            *log.Logger
	meta           interface { //storage interface
		bludgeon.MetaTimer
		bludgeon.MetaTimeSlice
	}
}

func NewServer(log, logError *log.Logger, meta interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}) interface {
	bludgeon.Logger
	Owner
	Manage
	Functional
} {
	//REVIEW: should add cases to confirm that meta/remote aren't nil since
	// basic functionality won't work?
	if meta == nil {
		panic("meta is nil")
	}
	//populate client pointer
	return &server{
		meta:     meta,
		log:      log,
		logError: logError,
		tokens:   make(map[string]bludgeon.Token),
	}
}

func (s *server) Println(v ...interface{}) {
	if s.log != nil {
		s.log.Println(v...)
	}
}

func (s *server) Printf(format string, v ...interface{}) {
	if s.log != nil {
		s.log.Printf(format, v...)
	}
}

func (s *server) Print(v ...interface{}) {
	if s.log != nil {
		s.log.Print(v...)
	}
}

func (s *server) Error(err error) {
	if s.logError != nil {
		s.logError.Println(err)
	}
}

func (s *server) Errorf(format string, v ...interface{}) {
	if s.logError != nil {
		s.logError.Printf(format, v...)
	}
}

func (s *server) launchPurge() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
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
	}()
	<-started
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
	s.config = common.Configuration{}
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
	Start(config common.Configuration) (chExternal <-chan struct{}, err error)

	//
	Stop() (err error)
}

func (s *server) Start(config common.Configuration) (chExternal <-chan struct{}, err error) {
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
	s.chExternal = make(chan struct{})
	chExternal = s.chExternal
	//launch go routines
	s.launchPurge()
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
	close(s.chExternal)
	//set started flag to false
	s.started = false

	return
}

type Functional interface {
	//CommandHandler
	CommandHandler(command bludgeon.CommandServer, dataIn interface{}, token bludgeon.Token) (dataOut interface{}, err error)
}

//ensure that cache implements Functional
var _ Functional = &server{}

func (s *server) CommandHandler(command bludgeon.CommandServer, dataIn interface{}, token bludgeon.Token) (dataOut interface{}, err error) {
	switch command {
	case bludgeon.CommandServerTimerCreate:
		dataOut, err = s.TimerCreate()
	case bludgeon.CommandServerTimerRead:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			dataOut, err = s.TimerRead(d.ID)
		}
	case bludgeon.CommandServerTimerUpdate:
		if timer, ok := dataIn.(bludgeon.Timer); !ok {
			err = errors.New("Unable to cast into timer")
		} else {
			dataOut, err = s.TimerUpdate(timer)
		}
	case bludgeon.CommandServerTimerDelete:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			err = s.TimerDelete(d.ID)
		}
	case bludgeon.CommandServerTimerStart:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			dataOut, err = s.TimerStart(d.ID, d.StartTime)
		}
	case bludgeon.CommandServerTimerPause:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			dataOut, err = s.TimerPause(d.ID, d.PauseTime)
		}
	case bludgeon.CommandServerTimerSubmit:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			dataOut, err = s.TimerSubmit(d.ID, d.FinishTime)
		}
	case bludgeon.CommandServerTimeSliceRead:
		if d, ok := dataIn.(common.CommandData); !ok {
			err = errors.New("Unable to cast into command data")
		} else {
			dataOut, err = s.TimeSliceRead(d.ID)
		}
	default:
		err = fmt.Errorf("command not supported: %s", command)
	}

	return
}

//
func (s *server) TimerCreate() (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//create the timer
	timer, err = bludgeon.TimerCreate(s.meta)

	return
}

//
func (s *server) TimerRead(id string) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//read the timer
	timer, err = bludgeon.TimerRead(id, s.meta)

	return
}

//
func (s *server) TimerUpdate(t bludgeon.Timer) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//update the timer
	timer, err = bludgeon.TimerUpdate(t, s.meta)

	return
}

//
func (s *server) TimerDelete(id string) (err error) {
	s.Lock()
	defer s.Unlock()

	//delete the timer
	err = bludgeon.TimerDelete(id, s.meta)

	return
}

func (s *server) TimerStart(id string, startTime time.Time) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//start the timer
	timer, err = bludgeon.TimerStart(id, startTime, s.meta)

	return
}

func (s *server) TimerPause(id string, pauseTime time.Time) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//pause the time
	timer, err = bludgeon.TimerPause(id, pauseTime, s.meta)

	return
}

func (s *server) TimerSubmit(id string, submitTime time.Time) (timer bludgeon.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//submit timer
	timer, err = bludgeon.TimerSubmit(id, submitTime, s.meta)

	return
}

func (s *server) TimeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	s.Lock()
	defer s.Unlock()

	//read the time slice
	timeSlice, err = bludgeon.TimeSliceRead(timeSliceID, s.meta)

	return
}

// type Token interface {
// 	//Generate creates a token, outputs a token
// 	Generate() string

// 	//Verify will check if a token exists
// 	Verify(tokenID string) bool

// 	//Purge deletes tokens that are as old as wait
// 	Purge()

// 	//Delete will remove a token
// 	Delete(token string) (deleted bool)
// }

// //Generate creates a token, outputs a token
// func (s *server) TokenGenerate() string {
// 	s.Lock()
// 	defer s.Unlock()

// 	b := make([]byte, 8)
// 	var t bludgeon.Token

// 	for {
// 		rand.Read(b) //generate token
// 		token := fmt.Sprintf("%x-%x-%x-%x", b[0:2], b[2:4], b[4:6], b[6:8])

// 		if _, ok := s.tokens[token]; !ok {
// 			t.Token = token
// 			t.Time = time.Now().UnixNano()
// 			s.tokens[token] = t

// 			return t.Token
// 		}
// 	}
// }

// //Verify will check if a token exists
// func (s *server) TokenVerify(tokenID string) (valid bool) {
// 	s.RLock()
// 	defer s.RUnlock()

// 	if token, ok := s.tokens[tokenID]; ok {
// 		if (time.Now().UnixNano() - token.Time) > s.config.TokenWait/int64(time.Nanosecond) {
// 			delete(s.tokens, tokenID)
// 		} else {
// 			valid = true
// 		}
// 	}

// 	return
// }

// //Purge deletes tokens that are as old as wait
// func (s *server) TokenPurge() {
// 	s.Lock()
// 	defer s.Unlock()

// 	for key, token := range s.tokens {
// 		if time.Now().UnixNano()-token.Time > s.config.TokenWait {
// 			delete(s.tokens, key)
// 		}
// 	}
// }

// //Delete will remove a token
// func (s *server) TokenDelete(token string) (deleted bool) {
// 	s.Lock()
// 	defer s.Unlock()

// 	if _, ok := s.tokens[token]; ok {
// 		delete(s.tokens, token)
// 		deleted = true
// 	}

// 	return
// }
