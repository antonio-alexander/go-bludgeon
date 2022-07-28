package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

const LogAlias string = "GrpcServer"

type grpcServer struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	*grpc.Server
	config      *Configuration
	initialized bool
}

func New(parameters ...interface{}) interface {
	Owner
	grpc.ServiceRegistrar
} {
	var config *Configuration

	s := &grpcServer{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			config = p
		case logger.Logger:
			s.Logger = p
		}
	}
	if config != nil {
		if err := s.Initialize(config); err != nil {
			panic(err)
		}
	}
	if s.Logger == nil {
		s.Logger = logger.New()
	}
	return s
}

func (s *grpcServer) launchServe(listener net.Listener) {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		close(started)
		if err := s.Serve(listener); err != nil {
			s.Error("%s %s", LogAlias, err)
		}
	}()
	<-started
	s.Info("%s listening on %v", LogAlias, listener.Addr())
}

func (s *grpcServer) Initialize(config *Configuration, registrations ...RegisterFx) error {
	s.Lock()
	defer s.Unlock()

	if s.initialized {
		return errors.New("already initialized")
	}
	s.config = config
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
	if err != nil {
		return err
	}
	s.Server = grpc.NewServer(s.config.Options...)
	for _, registration := range registrations {
		registration()
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
	s.GracefulStop()
	s.Wait()
	s.initialized = false
}
