package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/common"

	"github.com/gorilla/mux"
)

type server struct {
	sync.RWMutex
	sync.WaitGroup
	common.Logger
	router  *mux.Router
	server  *http.Server
	started bool
}

type Server interface {
	Start(address, port string) (err error)
	Started() bool
	BuildRoutes(routes []HandleFuncConfig) (err error)
	Stop()
}

func New(logger common.Logger) Server {
	return &server{
		router: mux.NewRouter(),
		Logger: logger,
	}
}

//Start uses the configured mux/router to start listening to responses via REST
func (s *server) Start(address, port string) (err error) {
	s.Lock()
	defer s.Unlock()

	if s.started {
		return
	}
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: s.router,
	}
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		close(started)
		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.Debug("Httpserver: ListenAndServe() error: %s", err)
			}
		}
		//REVIEW: Do we need to account for a situation where the rest server kills itself
		// unexepctedly?
	}()
	<-started
	s.started = true

	return
}

func (s *server) Started() bool {
	s.RLock()
	defer s.RUnlock()
	return s.started
}

func (s *server) BuildRoutes(routes []HandleFuncConfig) (err error) {
	s.Lock()
	defer s.Unlock()

	for _, route := range routes {
		s.router.HandleFunc(route.Route, route.HandleFx).Methods(route.Method)
	}

	return
}

func (s *server) Stop() {
	s.Lock()
	defer s.Unlock()

	if !s.started {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), ConfigShutdownTimeout)
	defer cancel()
	s.server.Shutdown(ctx)
	s.Wait()
}
