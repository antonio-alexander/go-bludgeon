package restserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

type server struct {
	sync.RWMutex
	sync.WaitGroup
	context.Context
	logger.Logger
	*mux.Router
	*http.Server
	cancel      context.CancelFunc
	initialized bool
	configured  bool
	config      *Configuration
}

func New() interface {
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
	Router
	context.Context
} {
	router := mux.NewRouter()
	return &server{
		Logger:  logger.NewNullLogger(),
		Context: context.Background(),
		Router:  router,
		Server: &http.Server{
			Handler: router,
		},
	}
}

func (s *server) launchServer() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.WaitGroup.Done()

		if !s.config.CorsDisabled {
			s.Server.Handler = cors.New(cors.Options{
				AllowedOrigins:   s.config.AllowedOrigins,
				AllowCredentials: s.config.AllowCredentials,
				AllowedMethods:   s.config.AllowedMethods,
				Debug:            s.config.CorsDebug,
			}).Handler(s.Router)
			if s.config.AllowCredentials {
				s.Debug(logAlias + "CORS configured with Allow Credentials")
			}
			if s.config.CorsDebug {
				s.Debug(logAlias + "CORS configured with Debug")
			}
			if len(s.config.AllowedOrigins) > 0 {
				s.Debug(logAlias+"CORS configured with Allowed Origins \"%s\"", strings.Join(s.config.AllowedOrigins, ","))
			}
			if len(s.config.AllowedMethods) > 0 {
				s.Debug(logAlias+"CORS configured with Allowed Methods \"%s\"", strings.Join(s.config.AllowedMethods, ","))
			}
		} else {
			s.Debug(logAlias + "CORS disabled")
		}
		close(started)
		if err := s.Server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.Error(logAlias+"%s %s", err)
			}
		}
		//REVIEW: Do we need to account for a situation where the rest server kills itself
		// unexepctedly?
	}()
	<-started
}

func (s *server) Done() <-chan struct{} {
	return s.Context.Done()
}

func (s *server) HandleFunc(config HandleFuncConfig) {
	switch config.Method {
	default:
		s.Router.HandleFunc(config.Route, config.HandleFx).Methods(config.Method)
	case "":
		s.Router.HandleFunc(config.Route, config.HandleFx)
	}
}

func (s *server) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (s *server) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *server) Configure(items ...interface{}) error {
	var envs map[string]string
	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		c = new(Configuration)
		c.Default()
		c.FromEnv(envs)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	s.config = c
	s.configured = true
	return nil
}

func (s *server) Initialize() (err error) {
	s.Lock()
	defer s.Unlock()
	if s.initialized {
		return errors.New(ErrStarted)
	}
	if !s.configured {
		return errors.New("not configured)")
	}
	s.Context, s.cancel = context.WithCancel(context.Background())
	s.Server.Addr = fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
	s.launchServer()
	s.initialized = true
	s.Info(logAlias+"listening on %s:%s", s.config.Address, s.config.Port)
	return
}

func (s *server) Shutdown() {
	s.Lock()
	defer s.Unlock()
	if !s.initialized {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Error(logAlias+"error while shutting down the server: %s", err)
	}
	s.cancel()
	s.Wait()
	s.initialized, s.configured = false, false
	s.Info(logAlias + "stopped")
}
