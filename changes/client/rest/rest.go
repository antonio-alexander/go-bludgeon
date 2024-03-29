package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"syscall"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal_cache "github.com/antonio-alexander/go-bludgeon/changes/internal/cache"
	internal_queue "github.com/antonio-alexander/go-bludgeon/changes/internal/queue"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	restclient "github.com/antonio-alexander/go-bludgeon/internal/rest/client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type restClient struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	handlers map[string]*handler
	queue    internal_queue.Queue
	cache    interface {
		internal_cache.Cache
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
	}
	initialized bool
	configured  bool
	config      *Configuration
	client      interface {
		restclient.Client
		internal.Configurer
		internal.Parameterizer
	}
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
	cache := internal_cache.New()

	queue := internal_queue.New(QueueSize)
	return &restClient{
		client:   restclient.New(),
		handlers: make(map[string]*handler),
		cache:    cache,
		queue:    queue,
		Logger:   logger.NewNullLogger(),
	}
}

func (r *restClient) cacheRead(changeId string) *data.Change {
	if r.config.DisableCache {
		return nil
	}
	return r.cache.Read(changeId)
}

func (r *restClient) cacheWrite(change *data.Change) {
	if r.config.DisableCache {
		return
	}
	r.cache.Write(change)
}

func (r *restClient) cacheDelete(changeId string) {
	if r.config.DisableCache {
		return
	}
	r.cache.Delete(changeId)
}

func (r *restClient) queueEnqueue(item interface{}) bool {
	return r.queue.Enqueue(item)
}

func (r *restClient) launchUpsertQueue() {
	if r.config.DisableQueue {
		return
	}
	started := make(chan struct{})
	r.Add(1)
	go func() {
		defer r.Done()

		tQueueReadeRate := time.NewTicker(r.config.UpsertQueueRate)
		defer tQueueReadeRate.Stop()
		queueReadFx := func() {
			var changePartialsToEnqueue []data.ChangePartial

			for _, changePartial := range internal_queue.ChangePartialFlush(r.queue) {
				if _, err := r.changeUpsert(r.ctx, changePartial); err != nil {
					r.Error("error while attempting to upsert: %s", err)
					changePartialsToEnqueue = append(changePartialsToEnqueue,
						changePartial)
				}
			}
			if _, overflow := internal_queue.ChangePartialEnqueueMultiple(r.queue, changePartialsToEnqueue); overflow {
				r.Info("overflow encountered while enqueueing multiple changes")
			}
		}
		close(started)
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-tQueueReadeRate.C:
				queueReadFx()
			}
		}
	}()
	<-started
}

func (r *restClient) doRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error) {
	bytes, statusCode, err := r.client.DoRequest(ctx, uri, method, data)
	if err != nil {
		return nil, err
	}
	switch statusCode {
	default:
		err := new(internal_errors.Error)
		if err := json.Unmarshal(bytes, err); err != nil {
			r.Error("error while unmarshalling error: %s", err)
			return nil, errors.Errorf("failed status code: %d", statusCode)
		}
		return nil, err
	case http.StatusOK, http.StatusNoContent, http.StatusNotModified:
		return bytes, nil
	}
}

func (r *restClient) changeUpsert(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error) {
	bytes, err := json.Marshal(&changePartial)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChanges, r.config.Rest.Address, r.config.Rest.Port)
	bytes, err = r.doRequest(ctx, uri, data.MethodChangeUpsert, bytes)
	if err != nil {
		return nil, err
	}
	change := &data.Change{}
	if err = json.Unmarshal(bytes, change); err != nil {
		return nil, err
	}
	return change, nil
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
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case internal_queue.Queue:
			if r.queue != nil {
				r.queue.Close()
			}
			r.queue = p
		case interface {
			internal_cache.Cache
			internal.Configurer
			internal.Initializer
			internal.Parameterizer
		}:
			r.cache = p
		case interface {
			restclient.Client
			internal.Configurer
			internal.Parameterizer
		}:
			r.client = p
		}
	}
	switch {
	case r.queue == nil:
		panic("queue is nil")
	case r.cache == nil:
		panic("cache is nil")
	case r.client == nil:
		panic("client is nil")
	}
	r.client.SetParameters(parameters...)
	r.cache.SetParameters(parameters...)
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
		c = NewConfiguration()
		c.Default()
		c.FromEnv(envs)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	if err := r.client.Configure(c.Rest); err != nil {
		return err
	}
	if !c.DisableCache {
		if err := r.cache.Configure(c.Cache); err != nil {
			return err
		}
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
	r.launchUpsertQueue()
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
	if !r.config.DisableQueue {
		for _, changePartial := range internal_queue.ChangePartialFlush(r.queue) {
			if _, err := r.changeUpsert(context.Background(), changePartial); err != nil {
				r.queueEnqueue(changePartial)
			}
		}
	}
	r.initialized = false
}

func (r *restClient) ChangeUpsert(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error) {
	r.Lock()
	defer r.Unlock()

	if !r.initialized {
		return nil, errors.New("not initialized")
	}
	change, err := r.changeUpsert(ctx, changePartial)
	if err != nil {
		r.Error("error while upserting change: %s", err)
		switch {
		default:
			return nil, err
		case errors.Is(err, syscall.ECONNREFUSED),
			errors.Is(err, syscall.ECONNRESET),
			errors.Is(err, syscall.ECONNABORTED),
			errors.Is(err, syscall.ETIMEDOUT):
			if !r.config.DisableQueue {
				r.Trace("attempting to enqueue change")
				if overflow := r.queueEnqueue(changePartial); overflow {
					r.Trace("failed to enqueue change")
					return nil, err
				}
				r.Trace("successfully enqueued change")
				return changePartialToChange(changePartial), nil
			}
			return nil, err
		}
	}
	r.cacheWrite(change)
	return change, nil
}

func (r *restClient) ChangeRead(ctx context.Context, changeId string) (*data.Change, error) {
	if change := r.cacheRead(changeId); change != nil {
		return change, nil
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesParamf, r.config.Rest.Address, r.config.Rest.Port, changeId)
	bytes, err := r.doRequest(ctx, uri, data.MethodChangeRead, nil)
	if err != nil {
		return nil, err
	}
	change := &data.Change{}
	if err = json.Unmarshal(bytes, change); err != nil {
		return nil, err
	}
	r.cacheWrite(change)
	return change, nil
}

func (r *restClient) ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesSearch+search.ToParams(), r.config.Rest.Address, r.config.Rest.Port)
	bytes, err := r.doRequest(ctx, uri, data.MethodChangeRead, nil)
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
	if _, err := r.doRequest(ctx, uri, data.MethodChangeDelete, nil); err != nil {
		return err
	}
	r.cacheDelete(changeId)
	return nil
}

func (r *restClient) RegistrationUpsert(ctx context.Context, registrationId string) error {
	bytes, err := json.Marshal(&data.RequestRegister{RegistrationId: registrationId})
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistration,
		r.config.Rest.Address, r.config.Rest.Port)
	if _, err := r.doRequest(ctx, uri, data.MethodRegistrationUpsert, bytes); err != nil {
		return err
	}
	return nil
}

func (r *restClient) RegistrationChangesRead(ctx context.Context, registrationId string) ([]*data.Change, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamChangesf,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	bytes, err := r.doRequest(ctx, uri, data.MethodChangeRead, nil)
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
	if _, err = r.doRequest(ctx, uri, data.MethodRegistrationChangeAcknowledge, bytes); err != nil {
		return err
	}
	return nil
}

func (r *restClient) RegistrationDelete(ctx context.Context, registrationId string) error {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteChangesRegistrationParamf,
		r.config.Rest.Address, r.config.Rest.Port, registrationId)
	if _, err := r.doRequest(ctx, uri, data.MethodRegistrationDelete, nil); err != nil {
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
	return handler.client.IsConnected(), nil
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
