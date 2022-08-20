package logic

import (
	"context"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-queue/finite"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	goqueue "github.com/antonio-alexander/go-queue"
	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"
)

type handler struct {
	stopper chan struct{}
	queue   interface {
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
		goqueue.Owner
	}
}

type logic struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	meta.Change
	meta.Registration
	meta.RegistrationChange
	ctx         context.Context
	cancel      context.CancelFunc
	handlersMux sync.RWMutex
	handlers    map[string]handler
	initialized bool
}

func changeDequeue(queue goqueue.Dequeuer) (*data.Change, bool) {
	item, underflow := queue.Dequeue()
	if underflow {
		return nil, true
	}
	change, _ := item.(*data.Change)
	return change, false
}

//New will generate a new instance of logic that implements
// the interfaces Logic and Owner, from the provided parameters
// we can set the logger and the employee meta (required)
func New() interface {
	Logic
	internal.Initializer
	internal.Parameterizer
} {
	return &logic{
		handlers: make(map[string]handler),
		Logger:   logger.NewNullLogger(),
		ctx:      context.Background(),
	}
}

func (l *logic) handlersWrite(stopper chan struct{}, queue interface {
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.Event
	goqueue.Owner
}) (handlerId string) {
	l.handlersMux.Lock()
	defer l.handlersMux.Unlock()
	handlerId = uuid.Must(uuid.NewRandom()).String()
	l.handlers[handlerId] = handler{
		stopper: stopper,
		queue:   queue,
	}
	return
}

func (l *logic) handlersRead() (handlers map[string]goqueue.Enqueuer) {
	l.handlersMux.RLock()
	defer l.handlersMux.RUnlock()
	handlers = make(map[string]goqueue.Enqueuer)
	for handlerId, handler := range l.handlers {
		handlers[handlerId] = handler.queue
	}
	return
}

func (l *logic) handlersDelete(handlerId string) error {
	l.handlersMux.Lock()
	defer l.handlersMux.Unlock()
	handler, ok := l.handlers[handlerId]
	if !ok {
		return errors.New("handler not found")
	}
	delete(l.handlers, handlerId)
	select {
	default:
		close(handler.stopper)
	case <-handler.stopper:
	}
	handler.queue.Close()
	return nil
}

func (l *logic) launchHandler(handlerId string, handleFx HandlerFx, stopper chan struct{}, queue interface {
	goqueue.Dequeuer
	goqueue.Event
	goqueue.Owner
}) {
	started := make(chan struct{})
	l.Add(1)
	go func() {
		defer l.Done()
		defer func() { l.Debug(logAlias+"stopped handler: %s", handlerId) }()
		defer func() {
			select {
			default:
				close(stopper)
			case <-stopper:
			}
		}()

		signalIn := queue.GetSignalIn()
		close(started)
		l.Debug(logAlias+"started handler: %s", handlerId)
		for {
			select {
			case <-stopper:
				return
			case <-l.ctx.Done():
				return
			case <-signalIn:
				if change, underflow := changeDequeue(queue); !underflow {
					handleFx(l.ctx, handlerId, []*data.Change{change})
				}
			}
		}
	}()
	<-started
}

func (l *logic) changeBroadcast(change *data.Change) {
	handlers := l.handlersRead()
	if len(handlers) == 0 {
		l.Debug(logAlias + "received change to broadcast, but no handlers")
		return
	}
	for handlerId, queue := range l.handlersRead() {
		if overflow := queue.Enqueue(change); overflow {
			l.Trace(logAlias+"overflow while attempting to broadcast change %s to %s", change.Id, handlerId)
			continue
		}
		l.Trace(logAlias+"broadcasted change %s to %s", change.Id, handlerId)
	}
}

func (l *logic) SetUtilities(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case logger.Logger:
			l.Logger = p
		}
	}
}

func (l *logic) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			meta.Change
			meta.Registration
			meta.RegistrationChange
		}:
			l.Change = p
			l.Registration = p
			l.RegistrationChange = p
		case meta.Registration:
			l.Registration = p
		case meta.RegistrationChange:
			l.RegistrationChange = p
		case meta.Change:
			l.Change = p
		}
	}
	switch {
	case l.Registration == nil:
		panic(PanicRegistrationMetaNotSet)
	case l.RegistrationChange == nil:
		panic(PanicRegistrationChangeMetaNotSet)
	case l.Change == nil:
		panic(PanicChangeMetaNotSet)
	}
}

func (l *logic) Initialize() error {
	l.Lock()
	defer l.Unlock()
	if l.initialized {
		return errors.New("logic already initialized")
	}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.initialized = true
	l.Info(logAlias + "initialized")
	return nil
}

func (l *logic) Shutdown() {
	l.Lock()
	defer l.Unlock()
	if !l.initialized {
		return
	}
	l.cancel()
	for handlerId := range l.handlersRead() {
		l.handlersDelete(handlerId)
	}
	l.Wait()
	l.Info(logAlias + "shutdown")
	l.initialized = false
}

func (l *logic) ChangeUpsert(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error) {
	change, err := l.Change.ChangeCreate(ctx, changePartial)
	if err != nil {
		return nil, err
	}
	if err := l.RegistrationChangeUpsert(ctx, change.Id); err != nil {
		return nil, err
	}
	l.Trace(logAlias+"upserted change: %s", change.Id)
	l.changeBroadcast(change)
	return change, nil
}

func (l *logic) RegistrationChangesRead(ctx context.Context, registrationId string) ([]*data.Change, error) {
	changeIds, err := l.RegistrationChange.RegistrationChangesRead(ctx, registrationId)
	if err != nil {
		return nil, err
	}
	//KIM: because this change is inclusive, if we provide an empty list of change ids
	// it will return ALL changes
	if len(changeIds) == 0 {
		return nil, nil
	}
	return l.ChangesRead(ctx, data.ChangeSearch{ChangeIds: changeIds})
}

func (l *logic) RegistrationChangeAcknowledge(ctx context.Context, serviceName string, changeIds ...string) error {
	//REVIEW: there's an opportunity for data inconsistency if the
	// deletion of the changes fail, maybe store them in memory?
	changeIdsToDelete, err := l.RegistrationChange.RegistrationChangeAcknowledge(ctx, serviceName, changeIds...)
	if err != nil {
		return err
	}
	l.Trace(logAlias+"acknowledged change(s) %v for %s", changeIds, serviceName)
	if len(changeIdsToDelete) > 0 {
		if err := l.ChangesDelete(ctx, changeIdsToDelete...); err != nil {
			return err
		}
		l.Trace(logAlias+"deleted change(s): %v", changeIdsToDelete)
	}
	return nil
}

func (l *logic) HandlerCreate(ctx context.Context, handleFx HandlerFx) (string, error) {
	stopper, queue := make(chan struct{}), finite.New(QueueSize)
	handlerId := l.handlersWrite(stopper, queue)
	l.launchHandler(handlerId, handleFx, stopper, queue)
	l.Trace(logAlias+"created handler: %s", handlerId)
	return handlerId, nil
}

func (l *logic) HandlerDelete(ctx context.Context, handlerId string) error {
	if err := l.handlersDelete(handlerId); err != nil {
		return err
	}
	l.Trace(logAlias+"deleted handler %s", handlerId)
	return nil
}
