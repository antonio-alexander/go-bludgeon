package service

import (
	"encoding/json"
	"net/http"

	data "github.com/antonio-alexander/go-bludgeon/healthcheck/data"
	logic "github.com/antonio-alexander/go-bludgeon/healthcheck/logic"

	common "github.com/antonio-alexander/go-bludgeon/common"
	logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	server "github.com/antonio-alexander/go-bludgeon/pkg/rest/server"
)

type restServer struct {
	logger.Logger
	logic logic.Logic
}

func New() interface {
	common.Parameterizer
	server.RouteBuilder
} {
	return &restServer{
		Logger: logger.NewNullLogger(),
	}
}

func (s *restServer) endpointHealthCheck() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var healthCheck *data.HealthCheck
		var bytes []byte
		var err error

		ctx := request.Context()
		if healthCheck, err = s.logic.HealthCheck(ctx); err == nil {
			bytes, err = json.Marshal(healthCheck)
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employee create -  %s", err)
		}
	}
}

func (s *restServer) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logic.Logic:
			s.logic = p
		}
	}
	switch {
	case s.logic == nil:
		panic("logic not set")
	}
}

func (s *restServer) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *restServer) BuildRoutes() []server.HandleFuncConfig {
	return []server.HandleFuncConfig{
		{Route: data.RouteHealthCheck, Method: http.MethodGet, HandleFx: s.endpointHealthCheck()},
	}
}
