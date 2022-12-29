package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	logic "github.com/antonio-alexander/go-bludgeon/changes/logic"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	internal "github.com/antonio-alexander/go-bludgeon/internal"

	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_rest "github.com/antonio-alexander/go-bludgeon/internal/rest/server"
	internal_websocket "github.com/antonio-alexander/go-bludgeon/internal/websocket/server"

	"github.com/pkg/errors"
)

type restServer struct {
	sync.WaitGroup
	logger.Logger
	ctx    context.Context
	logic  logic.Logic
	router internal_rest.Router
}

func New() interface {
	internal.Parameterizer
} {
	return &restServer{
		Logger: logger.NewNullLogger(),
		ctx:    context.Background(),
	}
}

func (s *restServer) handleResponse(writer http.ResponseWriter, err error, item interface{}) error {
	if err != nil {
		switch {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, meta.ErrChangeNotFound):
			writer.WriteHeader(http.StatusNotFound)
		case errors.Is(err, meta.ErrChangeNotWritten):
			writer.WriteHeader(http.StatusNotModified)
		case errors.Is(err, meta.ErrChangeConflictWrite):
			writer.WriteHeader(http.StatusConflict)
		}
		//REVIEW: I don't like to throw away/shadow this error, but it's not
		bytes, jsonErr := json.Marshal(&internal_errors.Error{Error: err.Error()})
		if jsonErr != nil {
			s.Error(logAlias+"json error on handle response: %s", jsonErr)
		}
		_, err = writer.Write(bytes)
		return err
	}
	if item != nil {
		bytes, err := json.Marshal(item)
		if err != nil {
			s.Error(logAlias+"json error on handle response: %s", err)
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = writer.Write(bytes)
		return err
	}
	writer.WriteHeader(http.StatusNoContent)
	_, err = writer.Write([]byte{})
	return err
}

func (s *restServer) endpointChangeUpsert() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var change *data.Change
		var bytes []byte
		var err error

		changeRequest, changeResponse := &data.RequestChange{}, &data.ResponseChange{}
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, &changeRequest); err == nil {
				if change, err = s.logic.ChangeUpsert(request.Context(), changeRequest.ChangePartial); err == nil {
					changeResponse.Change = *change
					s.Debug(logAlias+"upserted change: %s", change.Id)
				}
			}
		}
		if err = s.handleResponse(writer, err, changeResponse); err != nil {
			s.Error(logAlias+"upserted changes: %s", err)
			return
		}
	}
}

func (s *restServer) endpointChangeRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var change *data.Change
		var err error

		changeResponse := &data.ResponseChange{}
		if changeId, _ := valueFromPath(data.PathChangeId, internal_rest.Vars(request)); err == nil {
			if change, err = s.logic.ChangeRead(request.Context(), changeId); err == nil {
				changeResponse.Change = *change
			}
		}
		if err = s.handleResponse(writer, err, changeResponse); err != nil {
			s.Error(logAlias+"change read:  %s", err)
		}
	}
}

func (s *restServer) endpointChangesRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var search data.ChangeSearch

		search.FromParams(request.URL.Query())
		changes, err := s.logic.ChangesRead(request.Context(), search)
		if err = s.handleResponse(writer, err, &data.ChangeDigest{
			Changes: changes,
		}); err != nil {
			s.Error(logAlias+"changes read: %s", err)
		}
	}
}

func (s *restServer) endpointChangeDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		changeId, _ := valueFromPath(data.PathChangeId, internal_rest.Vars(request))
		err := s.logic.ChangesDelete(request.Context(), changeId)
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error(logAlias+"change delete: %s", err)
			return
		}
		s.Debug(logAlias+"deleted change: %s", changeId)
	}
}

func (s *restServer) endpointRegistrationUpsert() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		registrationId, _ := valueFromPath(data.PathRegistrationId, internal_rest.Vars(request))
		err := s.logic.RegistrationUpsert(request.Context(), registrationId)
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error(logAlias+"registration upsert -  %s", err)
			return
		}
		s.Debug(logAlias+"upserted registration: %s", registrationId)
	}
}

func (s *restServer) endpointRegistrationChangesRead() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		registrationId, _ := valueFromPath(data.PathRegistrationId, internal_rest.Vars(request))
		changes, err := s.logic.RegistrationChangesRead(request.Context(), registrationId)
		if err = s.handleResponse(writer, err, &data.ChangeDigest{
			Changes: changes,
		}); err != nil {
			s.Error(logAlias+"registration changes read: %s", err)
		}
	}
}

func (s *restServer) endpointRegistrationChangeAcknowledge() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bytes []byte
		var err error

		registrationId, _ := valueFromPath(data.PathRegistrationId, internal_rest.Vars(request))
		acknowledgeRequest := &data.RequestAcknowledge{}
		if bytes, err = io.ReadAll(request.Body); err == nil {
			if err = json.Unmarshal(bytes, acknowledgeRequest); err == nil {
				if err = s.logic.RegistrationChangeAcknowledge(request.Context(),
					registrationId, acknowledgeRequest.ChangeIds...); err == nil {
					s.Debug(logAlias+"%s acknowledged change(s) %v",
						registrationId, acknowledgeRequest.ChangeIds)
				}
			}
		}
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error(logAlias+"change acknowledge: %s", err)
		}
	}
}

func (s *restServer) endpointRegistrationDelete() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		registrationId, _ := valueFromPath(data.PathRegistrationId, internal_rest.Vars(request))
		err := s.logic.RegistrationDelete(request.Context(), registrationId)
		if err = s.handleResponse(writer, err, nil); err != nil {
			s.Error(logAlias+"registration delete -  %s", err)
			return
		}
		s.Debug(logAlias+"deleted registration: %s", registrationId)
	}
}

func (s *restServer) endpointWebsocket() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ws := internal_websocket.New(writer, request, s.Logger)
		if ws == nil {
			err := errors.New("unable to create websocket")
			if err := s.handleResponse(writer, err, nil); err != nil {
				s.Error(logAlias+"websocket -  %s", err)
			}
		}
		_, err := s.logic.HandlerCreate(s.ctx, func(ctx context.Context, handlerId string, changes []*data.Change) error {
			for _, change := range changes {
				wrapper := data.ToWrapper(change)
				if err := ws.Write(wrapper); err != nil {
					s.Error(logAlias+"error while handling change: %s", err)
					if err := s.logic.HandlerDelete(ctx, handlerId); err != nil {
						s.Error(logAlias+"error while deleting handler: %s", err)
					}
					ws.Close()
					return nil
				}
			}
			return nil
		})
		if err != nil {
			err := errors.New("unable to create websocket")
			if err := s.handleResponse(writer, err, nil); err != nil {
				s.Error(logAlias+"websocket -  %s", err)
			}
		}
	}
}

func (s *restServer) buildRoutes() {
	for _, route := range []internal_rest.HandleFuncConfig{
		{Route: data.RouteChangesWebsocket, HandleFx: s.endpointWebsocket()},
		{Route: data.RouteChanges, Method: data.MethodChangeUpsert, HandleFx: s.endpointChangeUpsert()},
		{Route: data.RouteChangesSearch, Method: data.MethodChangeRead, HandleFx: s.endpointChangesRead()},
		{Route: data.RouteChangesParam, Method: data.MethodChangeRead, HandleFx: s.endpointChangeRead()},
		{Route: data.RouteChangesParam, Method: data.MethodChangeDelete, HandleFx: s.endpointChangeDelete()},
		{Route: data.RouteChangesRegistrationServiceIdAcknowledge, Method: data.MethodRegistrationChangeAcknowledge, HandleFx: s.endpointRegistrationChangeAcknowledge()},
		{Route: data.RouteChangesRegistration, Method: data.MethodRegistrationUpsert, HandleFx: s.endpointRegistrationUpsert()},
		{Route: data.RouteChangesRegistrationParamChanges, Method: data.MethodChangeRead, HandleFx: s.endpointRegistrationChangesRead()},
		{Route: data.RouteChangesRegistrationParam, Method: data.MethodRegistrationDelete, HandleFx: s.endpointRegistrationDelete()},
	} {
		s.router.HandleFunc(route)
	}
}

func (s *restServer) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			s.Logger = p
		}
	}
}

func (s *restServer) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			internal_rest.Router
			context.Context
		}:
			s.ctx = p
			s.router = p
			s.buildRoutes()
		case logic.Logic:
			s.logic = p
		case context.Context:
			s.ctx = p
		case internal_rest.Router:
			s.router = p
		}
	}
	switch {
	case s.logic == nil:
		panic("logic not set")
	case s.router == nil:
		panic("router not set")
	case s.ctx == nil:
		panic("context not set")
	}
}
