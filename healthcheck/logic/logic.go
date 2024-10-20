package logic

import (
	"context"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/healthcheck/data"

	common "github.com/antonio-alexander/go-bludgeon/common"
	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
)

type logic struct {
	sync.WaitGroup
	sync.RWMutex
	logger.Logger
}

func New() interface {
	Logic
	common.Parameterizer
	common.Shutdowner
} {
	return &logic{
		Logger: logger.NewNullLogger(),
	}
}

func (l *logic) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch parameter.(type) {
		//
		}
	}
}

func (l *logic) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			l.Logger = p
		}
	}
}

func (l *logic) Shutdown() {
	l.Lock()
	defer l.Unlock()

	l.Wait()
}

func (l *logic) HealthCheck(ctx context.Context) (*data.HealthCheck, error) {
	return &data.HealthCheck{
		Time: time.Now().UnixNano(),
	}, nil
}
