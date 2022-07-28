package grpcclient

import (
	"errors"
	"fmt"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	sync.RWMutex
	logger.Logger
	*grpc.ClientConn
	config      *Configuration
	initialized bool
}

func New(parameters ...interface{}) interface {
	Client
	grpc.ClientConnInterface
} {
	var config *Configuration

	s := &grpcClient{}
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
	return s
}

func (s *grpcClient) Initialize(config *Configuration) error {
	s.Lock()
	defer s.Unlock()

	if s.initialized {
		return errors.New("already initialized")
	}
	s.config = config
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port), s.config.Options...)
	if err != nil {
		return err
	}
	s.ClientConn = clientConn
	s.initialized = true
	return nil
}

func (s *grpcClient) Shutdown() {
	s.Lock()
	defer s.Unlock()

	if !s.initialized {
		return
	}
	if err := s.Close(); err != nil {
		s.Error("error while closing: %s", err)
	}
	s.initialized = false
}
