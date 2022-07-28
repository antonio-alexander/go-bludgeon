package restserver

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

const LogAlias string = "HttpServer"

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
	if s.Logger == nil {
		s.Logger = logger.New()
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
		if s.config.AllowCredentials {
			s.Debug("%s CORS configured with Allow Credentials", LogAlias)
		}
		if s.config.CorsDebug {
			s.Debug("%s CORS configured with Debug", LogAlias)
		}
		s.Debug("%s CORS configured with Allowed Origins \"%s\"", LogAlias, strings.Join(s.config.AllowedOrigins, ","))
		s.Debug("%s CORS configured with Allowed Methods \"%s\"", LogAlias, strings.Join(s.config.AllowedMethods, ","))
		close(started)
		if err := s.Server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.Error("%s %s", LogAlias, err)
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
	s.Info("%s listening on %s:%s", LogAlias, s.config.Address, s.config.Port)
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
		s.Error("%s %s", LogAlias, err)
	}
	close(s.stopper)
	s.Wait()
	s.started = false
	s.Info("%s stopped", LogAlias)
}
