package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/pkg/config"
	"github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"google.golang.org/grpc"
)

const LogAlias string = "GrpcServer"

type grpcServer struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	*grpc.Server
	initialized   bool
	configured    bool
	config        *Configuration
	registrations []RegisterFx
}

func New() interface {
	common.Configurer
	common.Initializer
	common.Parameterizer
	grpc.ServiceRegistrar
} {
	return &grpcServer{Logger: logger.NewNullLogger()}
}

func (s *grpcServer) launchServe(listener net.Listener) {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		close(started)
		if err := s.Server.Serve(listener); err != nil {
			s.Error("%s %s", LogAlias, err)
		}
	}()
	<-started
	s.Info("%s listening on %v", LogAlias, listener.Addr())
}

func (s *grpcServer) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case Registerer:
			s.registrations = append(s.registrations, p.Register)
		}
	}
}

func (s *grpcServer) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *grpcServer) Configure(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	s.config = c
	s.configured = true
	return nil
}

func (s *grpcServer) Initialize() error {
	s.Lock()
	defer s.Unlock()

	if s.initialized {
		return errors.New("already initialized")
	}
	if !s.configured {
		return errors.New("not configured")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
	if err != nil {
		return err
	}
	s.Server = grpc.NewServer(s.config.Options...)
	for _, registration := range s.registrations {
		registration(s.Server)
	}
	s.launchServe(listener)
	s.initialized = true
	return nil
}

func (s *grpcServer) Shutdown() {
	s.Lock()
	defer s.Unlock()

	if !s.initialized {
		return
	}
	s.Server.GracefulStop()
	s.Wait()
	s.configured, s.initialized = false, false
	s.registrations = nil
}
