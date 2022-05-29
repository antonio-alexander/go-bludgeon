package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

type server struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	*mux.Router
	*http.Server
	config  Configuration
	stopper chan struct{}
	started bool
}

func New(parameters ...interface{}) interface {
	Owner
	Router
} {
	var config *Configuration

	router := mux.NewRouter()
	s := &server{
		Router: router,
		Server: &http.Server{
			Handler: router,
		},
	}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			config = p
		case logger.Logger:
			s.Logger = p
		}
	}
	if config != nil {
		if err := s.Start(config); err != nil {
			panic(err)
		}
	}
	return s
}

func (s *server) launchServer() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		s.Server.Handler = cors.New(cors.Options{
			AllowedOrigins:   s.config.AllowedOrigins,
			AllowCredentials: s.config.AllowCredentials,
			AllowedMethods:   s.config.AllowedMethods,
			Debug:            s.config.CorsDebug,
		}).Handler(s.Router)
		s.Debug("HttpServer: CORS configured with Allow Credentials: %t", s.config.AllowCredentials)
		s.Debug("HttpServer: CORS configured with Allowed Origins: %s", strings.Join(s.config.AllowedOrigins, ","))
		s.Debug("HttpServer: CORS configured with Allowed Methods: %s", strings.Join(s.config.AllowedMethods, ","))
		s.Debug("HttpServer: CORS configured with Debug: %t", s.config.CorsDebug)
		close(started)
		if err := s.Server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.Debug("HttpServer: ListenAndServe() error: %s", err)
			}
		}
		//REVIEW: Do we need to account for a situation where the rest server kills itself
		// unexepctedly?
	}()
	<-started
}

func (s *server) HandleFunc(config HandleFuncConfig) {
	s.Router.HandleFunc(config.Route, config.HandleFx).Methods(config.Method)
}

func (s *server) Start(config *Configuration) (err error) {
	s.Lock()
	defer s.Unlock()
	if s.started {
		return errors.New(ErrStarted)
	}
	s.config = *config
	s.stopper = make(chan struct{})
	s.Server.Addr = fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
	s.launchServer()
	s.started = true
	s.Info("HttpServer: started, listening on %s:%s", s.config.Address, s.config.Port)
	return
}

func (s *server) Stop() {
	s.Lock()
	defer s.Unlock()
	if !s.started {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Error("HttpServer: stopping - %s", err)
	}
	close(s.stopper)
	s.Wait()
	s.started = false
	s.Info("HttpServer: stopped")
}
