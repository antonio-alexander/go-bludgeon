package restclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"

	"github.com/google/uuid"
)

type restClient struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	client interface {
		internal_restclient.Client
		internal.Configurer
		internal.Parameterizer
	}
	ctx         context.Context
	cancel      context.CancelFunc
	handlers    map[string]*handler
	initialized bool
	configured  bool
	config      *Configuration
}

// New can be used to create a concrete instance of the rest client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	client.Client
	client.Handler
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	r := &restClient{
		client:   internal_restclient.New(),
		handlers: make(map[string]*handler),
		Logger:   logger.NewNullLogger(),
	}
	return r
}

func (r *restClient) SetUtilities(parameters ...interface{}) {
	r.client.SetUtilities(parameters)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			r.Logger = p
		}
	}
}

func (r *restClient) SetParameters(parameters ...interface{}) {
	r.client.SetParameters(parameters)
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			internal_restclient.Client
			internal.Configurer
			internal.Parameterizer
		}:
			r.client = p
		}
	}
}

func (r *restClient) Configure(items ...interface{}) error {
	r.Lock()
	defer r.Unlock()

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
	if err := r.client.Configure(&c.Rest); err != nil {
		return err
	}
	r.config = c
	r.configured = true
	return nil
}

// Initialize can be used to ready the underlying pointer for use
func (r *restClient) Initialize() error {
	r.Lock()
	defer r.Unlock()

	if r.initialized {
		return nil
	}
	if !r.configured {
		return errors.New("not configured")
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.initialized = true
	return nil
}

func (r *restClient) Shutdown() {
	r.Lock()
	defer r.Unlock()

	if !r.initialized {
		return
	}
	r.cancel()
	r.Wait()
	r.initialized = false
}

func (r *restClient) ChangeUpsert(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error) {
	bytes, err := json.Marshal(&data.RequestChange{ChangePartial: changePartial})
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChanges, r.config.Rest.Address, r.config.Rest.Port)
	bytes, err = r.client.DoRequest(ctx, uri, data.MethodChangeUpsert, bytes)
	if err != nil {
		return nil, err
	}
	response := &data.ResponseChange{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	return &response.Change, nil
}

func (r *restClient) ChangeRead(ctx context.Context, changeId string) (*data.Change, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, r.config.Rest.Address, r.config.Rest.Port, changeId)
	bytes, err := r.client.DoRequest(ctx, uri, data.MethodChangeRead, nil)
	if err != nil {
		return nil, err
	}
	response := &data.ResponseChange{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	return &response.Change, nil
}

func (r *restClient) ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesSearch+search.ToParams(), r.config.Rest.Address, r.config.Rest.Port)
	bytes, err := r.client.DoRequest(ctx, uri, data.MethodChangeRead, nil)
	if err != nil {
		return nil, err
	}
	response := &data.ChangeDigest{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	return response.Changes, nil
}

func (r *restClient) ChangeDelete(ctx context.Context, changeId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, r.config.Rest.Address, r.config.Rest.Port, changeId)
	if _, err := r.client.DoRequest(ctx, uri, data.MethodChangeDelete, nil); err != nil {
		return err
	}
	return nil
}

func (r *restClient) RegistrationUpsert(ctx context.Context, registrationId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamf,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	if _, err := r.client.DoRequest(ctx, uri, data.MethodRegistrationUpsert, nil); err != nil {
		return err
	}
	return nil
}

func (r *restClient) RegistrationChangesRead(ctx context.Context, registrationId string) ([]*data.Change, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	bytes, err := r.client.DoRequest(ctx, uri, data.MethodChangeRead, nil)
	if err != nil {
		return nil, err
	}
	response := &data.ChangeDigest{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	return response.Changes, nil
}

func (r *restClient) RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) error {
	bytes, err := json.Marshal(&data.RequestAcknowledge{ChangeIds: changeIds})
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationServiceIdAcknowledgef,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	if _, err = r.client.DoRequest(ctx, uri, data.MethodRegistrationChangeAcknowledge, bytes); err != nil {
		return err
	}
	return nil
}

func (r *restClient) RegistrationDelete(ctx context.Context, registrationId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamf,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	if _, err := r.client.DoRequest(ctx, uri, data.MethodRegistrationDelete, nil); err != nil {
		return err
	}
	return nil
}

func (r *restClient) HandlerCreate(handlerFx client.HandlerFx) (string, error) {
	r.Lock()
	defer r.Unlock()

	handlerId := uuid.Must(uuid.NewRandom()).String()
	r.handlers[handlerId] = newHandler(r.ctx, r, handlerId, r.config, handlerFx)
	return handlerId, nil
}

func (r *restClient) HandlerConnected(handlerId string) (bool, error) {
	r.Lock()
	defer r.Unlock()

	handler, ok := r.handlers[handlerId]
	if !ok {
		return false, errors.New("handler not found")
	}
	return handler.isConnected(), nil
}

func (r *restClient) HandlerDelete(handlerId string) error {
	r.Lock()
	defer r.Unlock()

	handler, ok := r.handlers[handlerId]
	if !ok {
		return errors.New("handler not found")
	}
	delete(r.handlers, handlerId)
	handler.Close()
	return nil
}
