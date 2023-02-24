package grpcclient

import (
	"errors"
	"fmt"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"google.golang.org/grpc"
)

type grpcClient struct {
	sync.RWMutex
	logger.Logger
	*grpc.ClientConn
	initialized bool
	configured  bool
	config      *Configuration
}

func New() interface {
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
	grpc.ClientConnInterface
} {
	return &grpcClient{Logger: logger.NewNullLogger()}
}

func (s *grpcClient) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (s *grpcClient) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *grpcClient) Configure(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	var c *Configuration
	var envs map[string]string

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

func (s *grpcClient) Initialize() error {
	s.Lock()
	defer s.Unlock()

	if s.initialized {
		return errors.New("already initialized")
	}
	if !s.configured {
		return errors.New("not configured")
	}
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
	if err := s.ClientConn.Close(); err != nil {
		s.Error("error while closing: %s", err)
	}
	s.initialized, s.configured = false, false
}
