package rest

import (
	"errors"
	"sync"
	"time"

	common "github.com/antonio-alexander/go-bludgeon/common"
	rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	config "github.com/antonio-alexander/go-bludgeon/server/config"
)

type server struct {
	sync.RWMutex                 //mutex for threadsafe operations
	sync.WaitGroup               //waitgroup to track go routines
	common.Logger                //
	started        bool          //whether or not the business logic has starting
	config         config.Rest   //configuration
	stopper        chan struct{} //stopper to stop goRoutines
	chExternal     chan struct{}
	tokens         map[string]common.Token //map of tokens
	server         rest.Server             //
	meta           interface {             //storage interface
		common.MetaTimer
		common.MetaTimeSlice
	}
}

func New(logger common.Logger, meta interface {
	common.MetaTimer
	common.MetaTimeSlice
}) interface {
	Owner
	Manage
	common.FunctionalTimer
	common.FunctionalTimeSlice
} {
	//REVIEW: should add cases to confirm that meta/remote aren't nil since
	// basic functionality won't work?
	if meta == nil {
		panic("meta is nil")
	}
	//populate client pointer
	return &server{
		meta:   meta,
		Logger: logger,
		tokens: make(map[string]common.Token),
		server: rest.New(logger),
	}
}

func (s *server) launchPurge() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		tPurge := time.NewTicker(30 * time.Second)
		defer tPurge.Stop()
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

//buildRoutes will create all the routes and their functions to execute when received
func (s *server) buildRoutes() []rest.HandleFuncConfig {
	//REVIEW: in the future when we add tokens, we'll need to create some way to check tokens for
	// certain functions, we may need to implement varratics to add support for tokens for the
	// server actions, may not be able to re-use the existing endpoints
	return []rest.HandleFuncConfig{
		//timer
		{Route: common.RouteTimerCreate, Method: POST, HandleFx: TimerCreate(s, s)},
		{Route: common.RouteTimerRead, Method: POST, HandleFx: TimerRead(s, s)},
		{Route: common.RouteTimerUpdate, Method: POST, HandleFx: TimerUpdate(s, s)},
		{Route: common.RouteTimerDelete, Method: POST, HandleFx: TimerDelete(s, s)},
		{Route: common.RouteTimerStart, Method: POST, HandleFx: TimerStart(s, s)},
		{Route: common.RouteTimerPause, Method: POST, HandleFx: TimerPause(s, s)},
		{Route: common.RouteTimerSubmit, Method: POST, HandleFx: TimerSubmit(s, s)},
		//time slice
		{Route: common.RouteTimeSliceRead, Method: POST, HandleFx: TimeSliceRead(s, s)},
		//route stop
		{Route: common.RouteStop, Method: POST, HandleFx: Stop(s, s)},
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
	s.config = config.Rest{}
	//set internal pointers to nil
	s.meta = nil
	s.server = nil
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
	Start(config config.Rest) (chExternal <-chan struct{}, err error)

	//
	Stop() (err error)
}

func (s *server) Start(config config.Rest) (chExternal <-chan struct{}, err error) {
	s.Lock()
	defer s.Unlock()

	//check if already started, store configuration, create stopper,
	// slaunch go routines, set started to true
	if s.started {
		err = errors.New(ErrStarted)

		return
	}
	s.config = config
	s.stopper = make(chan struct{})
	s.chExternal = make(chan struct{})
	chExternal = s.chExternal
	routes := s.buildRoutes()
	if err = s.server.BuildRoutes(routes); err != nil {
		return
	}
	s.launchPurge()
	s.started = true

	return
}

func (s *server) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	//check if not started
	//close the stopper
	//set started flag to false
	if !s.started {
		err = errors.New(ErrNotStarted)

		return
	}
	s.server.Stop()
	close(s.stopper)
	s.Wait()
	close(s.chExternal)
	s.started = false

	return
}

var _ common.FunctionalTimer = &server{}

//
func (s *server) TimerCreate() (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//create the timer
	timer, err = common.TimerCreate(s.meta)

	return
}

//
func (s *server) TimerRead(id string) (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//read the timer
	timer, err = common.TimerRead(id, s.meta)

	return
}

//
func (s *server) TimerUpdate(t common.Timer) (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//update the timer
	timer, err = common.TimerUpdate(t, s.meta)

	return
}

//
func (s *server) TimerDelete(id string) (err error) {
	s.Lock()
	defer s.Unlock()

	//delete the timer
	err = common.TimerDelete(id, s.meta)

	return
}

func (s *server) TimerStart(id string, startTime time.Time) (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//start the timer
	timer, err = common.TimerStart(id, startTime, s.meta)

	return
}

func (s *server) TimerPause(id string, pauseTime time.Time) (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//pause the time
	timer, err = common.TimerPause(id, pauseTime, s.meta)

	return
}

func (s *server) TimerSubmit(id string, submitTime time.Time) (timer common.Timer, err error) {
	s.Lock()
	defer s.Unlock()

	//submit timer
	timer, err = common.TimerSubmit(id, submitTime, s.meta)

	return
}

var _ common.FunctionalTimeSlice = &server{}

func (s *server) TimeSliceRead(timeSliceID string) (timeSlice common.TimeSlice, err error) {
	s.Lock()
	defer s.Unlock()

	//read the time slice
	timeSlice, err = common.TimeSliceRead(timeSliceID, s.meta)

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
// 	var t common.Token

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
