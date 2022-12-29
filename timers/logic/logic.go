package logic

import (
	"context"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"
	employeesdata "github.com/antonio-alexander/go-bludgeon/employees/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

type logic struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	meta.Timer
	meta.TimeSlice
	stopper        chan struct{}
	changesClient  changesclient.Client
	changesHandler changesclient.Handler
	initialized    bool
	configured     bool
	handlerId      string
	config         *Configuration
}

// New will instantiate a logic pointer that
// implements the Logic interface
func New() interface {
	Logic
	internal.Initializer
	internal.Parameterizer
	internal.Configurer
} {
	return &logic{Logger: logger.NewNullLogger()}
}

func (l *logic) changeUpsert(changePartial changesdata.ChangePartial) {
	l.Add(1)
	go func() {
		defer l.Done()

		ctx, cancel := context.WithTimeout(context.Background(), l.config.ChangesTimeout)
		defer cancel()
		change, err := l.changesClient.ChangeUpsert(ctx, changePartial)
		if err != nil {
			l.Error("error while upserting change (%s:%s->%s): %s", *changePartial.DataType, *changePartial.DataId, *changePartial.DataAction, err)
			return
		}
		l.Debug("Upserted change: %s (%s:%s->%s)", change.Id, change.DataType, change.DataId, change.DataAction)
	}()
}
func (l *logic) registrationChangeAcknowledge(serviceName string, changeIds ...string) {
	if len(changeIds) <= 0 {
		return
	}
	l.Add(1)
	go func() {
		defer l.Done()

		ctx, cancel := context.WithTimeout(context.Background(), l.config.ChangesTimeout)
		defer cancel()
		if err := l.changesClient.RegistrationChangeAcknowledge(ctx, l.config.ChangesRegistrationId, changeIds...); err != nil {
			l.Error("error while acknowledging changes: %s", err)
		}
	}()
}

func (l *logic) handleChanges(changes ...*changesdata.Change) error {
	var changesToAcknowledge []string

	for _, change := range changes {
		switch {
		case change.DataType == employeesdata.ChangeTypeEmployee &&
			change.DataAction == employeesdata.ChangeActionDelete:
			timers, err := l.TimersRead(context.Background(), data.TimerSearch{EmployeeID: &change.DataId})
			if err != nil {
				l.Error("error while reading timers: %s", err)
				break
			}
			failure := false
			for _, timer := range timers {
				err := l.TimerDelete(context.Background(), timer.ID)
				if err != nil {
					l.Error("error while deleting timers: %s", err)
					failure = true
					break
				}
				l.Debug("Deleted timer: %s, due to deletion of employee: %s ", timer.ID, change.DataId)
			}
			if failure {
				break
			}
			changesToAcknowledge = append(changesToAcknowledge, change.Id)
		}
	}
	l.registrationChangeAcknowledge(l.config.ChangesRegistrationId, changesToAcknowledge...)
	return nil
}

func (l *logic) launchChangeHandler() {
	started := make(chan struct{})
	l.Add(1)
	go func() {
		defer l.Done()

		checkChangesFx := func() {
			ctx, cancel := context.WithTimeout(context.Background(), l.config.ChangesTimeout)
			defer cancel()
			changesRead, err := l.changesClient.RegistrationChangesRead(ctx, l.config.ChangesRegistrationId)
			if err != nil {
				l.Error("error while reading registration changes: %s", err)
				return
			}
			if len(changesRead) == 0 {
				return
			}
			if err := l.handleChanges(changesRead...); err != nil {
				l.Error("error while reading registration changes: %s", err)
			}
		}
		tCheck := time.NewTicker(l.config.ChangeRateRead)
		defer tCheck.Stop()
		close(started)
		for {
			select {
			case <-l.stopper:
				return
			case <-tCheck.C:
				checkChangesFx()
			}
		}
	}()
	<-started
}

func (l *logic) launchChangeRegistration() {
	started := make(chan struct{})
	l.Add(1)
	go func() {
		defer l.Done()

		var registered, handlerSet bool
		var err error

		tRegister := time.NewTicker(l.config.ChangeRateRegistration)
		defer tRegister.Stop()
		close(started)
		for {
			select {
			case <-l.stopper:
				return
			case <-tRegister.C:
				if !handlerSet {
					if l.handlerId, err = l.changesHandler.HandlerCreate(l.handleChanges); err != nil {
						l.Error("error while creating change handler: %s", err)
						break
					}
					l.Debug("Change handler created: %s (%s)", l.handlerId, l.config.ChangesRegistrationId)
					handlerSet = true
				}
				if !registered {
					ctx, cancel := context.WithTimeout(context.Background(), l.config.ChangesTimeout)
					defer cancel()
					if err := l.changesClient.RegistrationUpsert(ctx, l.config.ChangesRegistrationId); err != nil {
						l.Error("error while upserting change registration: %s", err)
						break
					}
					l.Debug("Change registration upserted for: %s", l.config.ChangesRegistrationId)
					registered = true
				}
				if handlerSet && registered {
					return
				}
			}
		}
	}()
	<-started
}

func (l *logic) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			meta.Timer
			meta.TimeSlice
		}:
			l.Timer = p
			l.TimeSlice = p
		case meta.Timer:
			l.Timer = p
		case meta.TimeSlice:
			l.TimeSlice = p
		case interface {
			changesclient.Handler
			changesclient.Client
		}:
			l.changesHandler = p
			l.changesClient = p
		case changesclient.Handler:
			l.changesHandler = p
		case changesclient.Client:
			l.changesClient = p
		}
	}
	switch {
	case l.changesHandler == nil:
		panic("changes handler not set")
	case l.changesClient == nil:
		panic("changes client not set")
	case l.TimeSlice == nil:
		panic("no meta found for time slice")
	case l.Timer == nil:
		panic("no meta found for timer")
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

func (l *logic) Configure(items ...interface{}) error {
	l.Lock()
	defer l.Unlock()

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
	l.config = c
	l.configured = true
	return nil
}

func (l *logic) Initialize() error {
	l.Lock()
	defer l.Unlock()
	if l.initialized {
		return nil
	}
	l.stopper = make(chan struct{})
	l.launchChangeHandler()
	l.launchChangeRegistration()
	l.initialized = true
	return nil
}

func (l *logic) Shutdown() {
	l.Lock()
	defer l.Unlock()

	if !l.initialized {
		return
	}
	close(l.stopper)
	l.Wait()
	if l.handlerId != "" {
		if err := l.changesHandler.HandlerDelete(l.handlerId); err != nil {
			l.Error("error while deleting change handler: %s", err)
		}
		l.Debug("Change handler deleted: %s (%s)", l.handlerId, l.config.ChangesRegistrationId)
	}
	l.initialized = false
}

// IsConnected can be used to determine whether or not
// the underlying change handler is connected
func (l *logic) IsConnected() bool {
	connected, err := l.changesHandler.HandlerConnected(l.handlerId)
	if err != nil && err.Error() != "handler not found" { // changesclient.ErrHandlerNotFound {
		l.Error("error while checking if handler connected: %s", err)
		return false
	}
	return connected
}

// TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (l *logic) TimerCreate(ctx context.Context, timerPartial data.TimerPartial) (*data.Timer, error) {
	timer, err := l.Timer.TimerCreate(ctx, timerPartial)
	if err != nil {
		return nil, err
	}
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &timer.LastUpdated,
		ChangedBy:       &timer.LastUpdatedBy,
		DataId:          &timer.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionCreate,
		DataVersion:     &timer.Version,
	})
	return timer, nil
}

// TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (l *logic) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	timer, err := l.Timer.TimerStart(ctx, id)
	if err != nil {
		return nil, err
	}
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &timer.LastUpdated,
		ChangedBy:       &timer.LastUpdatedBy,
		DataId:          &timer.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionStart,
		DataVersion:     &timer.Version,
	})
	return timer, nil
}

// TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (l *logic) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	timer, err := l.Timer.TimerStop(ctx, id)
	if err != nil {
		return nil, err
	}
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &timer.LastUpdated,
		ChangedBy:       &timer.LastUpdatedBy,
		DataId:          &timer.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionStop,
		DataVersion:     &timer.Version,
	})
	return timer, nil
}

// TimerDelete can be used to delete a timer if it exists
func (l *logic) TimerDelete(ctx context.Context, id string) error {
	if err := l.Timer.TimerDelete(ctx, id); err != nil {
		return err
	}
	tNow := time.Now().UnixNano()
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &tNow,
		DataId:          &id,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionDelete,
	})
	return nil
}

// TimerSubmit can be used to stop a timer and set completed to true
func (l *logic) TimerSubmit(ctx context.Context, id string, submitTime int64) (*data.Timer, error) {
	if submitTime <= 0 {
		submitTime = time.Now().UnixNano()
	}
	timer, err := l.Timer.TimerSubmit(ctx, id, submitTime)
	if err != nil {
		return nil, err
	}
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &timer.LastUpdated,
		ChangedBy:       &timer.LastUpdatedBy,
		DataId:          &timer.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionSubmit,
		DataVersion:     &timer.Version,
	})
	return timer, nil
}

// TimerUpdate can be used to update values a given timer
// not associated with timer operations, values such as:
// comment, archived and completed
func (l *logic) TimerUpdate(ctx context.Context, id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	timer, err := l.Timer.TimerUpdate(ctx, id, data.TimerPartial{
		Completed: timerPartial.Completed,
		Archived:  timerPartial.Archived,
		Comment:   timerPartial.Comment,
	})
	if err != nil {
		return nil, err
	}
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &timer.LastUpdated,
		ChangedBy:       &timer.LastUpdatedBy,
		DataId:          &timer.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeTimer,
		DataAction:      &data.ChangeActionUpdate,
		DataVersion:     &timer.Version,
	})
	return timer, nil
}
