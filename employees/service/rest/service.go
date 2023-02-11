package service

import (
	"encoding/json"
	"io"
	"net/http"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	server "github.com/antonio-alexander/go-bludgeon/internal/rest/server"

	"github.com/gorilla/mux"
)

type restServer struct {
	logger.Logger
	logic logic.Logic
}

func New() interface {
	internal.Parameterizer
	server.RouteBuilder
} {
	return &restServer{
		Logger: logger.NewNullLogger(),
	}
}

func (s *restServer) endpointEmployeeCreate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var employeePartial data.EmployeePartial
		var employee *data.Employee
		var bytes []byte
		var err error

		ctx := request.Context()
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &employeePartial); err == nil {
				if employee, err = s.logic.EmployeeCreate(ctx, employeePartial); err == nil {
					bytes, err = json.Marshal(employee)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employee create -  %s", err)
		}
	}
}

func (s *restServer) endpointEmployeeRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var employee *data.Employee
		var bytes []byte
		var err error

		ctx := request.Context()
		id := idFromPath(mux.Vars(request))
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if employee, err = s.logic.EmployeeRead(ctx, id); err == nil {
				bytes, err = json.Marshal(employee)
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employee create -  %s", err)
		}
	}
}

func (s *restServer) endpointEmployeesRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var employees []*data.Employee
		var bytes []byte
		var err error
		var search data.EmployeeSearch

		ctx := request.Context()
		search.FromParams(request.URL.Query())
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if employees, err = s.logic.EmployeesRead(ctx, search); err == nil {
				bytes, err = json.Marshal(employees)
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employees read -  %s", err)
		}
	}
}

func (s *restServer) endpointEmployeeUpdate() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var employeePartial data.EmployeePartial
		var employee *data.Employee
		var bytes []byte
		var err error

		ctx := request.Context()
		id := idFromPath(mux.Vars(request))
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &employeePartial); err == nil {
				if employee, err = s.logic.EmployeeUpdate(ctx, id, employeePartial); err == nil {
					bytes, err = json.Marshal(employee)
				}
			}
		}
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employee update -  %s", err)
		}
	}
}

func (s *restServer) endpointEmployeeDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error

		ctx := request.Context()
		id := idFromPath(mux.Vars(request))
		err = s.logic.EmployeeDelete(ctx, id)
		if err = handleResponse(writer, err, bytes); err != nil {
			s.Error("employee delete -  %s", err)
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
		{Route: data.RouteEmployees, Method: http.MethodPost, HandleFx: s.endpointEmployeeCreate()},
		{Route: data.RouteEmployeesSearch, Method: http.MethodGet, HandleFx: s.endpointEmployeesRead()},
		{Route: data.RouteEmployeesID, Method: http.MethodGet, HandleFx: s.endpointEmployeeRead()},
		{Route: data.RouteEmployeesID, Method: http.MethodPut, HandleFx: s.endpointEmployeeUpdate()},
		{Route: data.RouteEmployeesID, Method: http.MethodDelete, HandleFx: s.endpointEmployeeDelete()},
	}
}
