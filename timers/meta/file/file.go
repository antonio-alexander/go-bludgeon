package file

import (
	"context"
	"os"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	memory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	logger.Logger
	memory interface {
		meta.Serializer
		internal.Parameterizer
		internal.Initializer
	}
	file interface {
		internal_file.File
		internal.Configurer
		internal.Initializer
		internal.Parameterizer
	}
	meta.Timer
	meta.TimeSlice
}

func New() interface {
	meta.Timer
	meta.TimeSlice
	internal.Initializer
	internal.Parameterizer
	internal.Configurer
} {
	memory := memory.New()
	return &file{
		Logger:    logger.NewNullLogger(),
		file:      internal_file.New(),
		memory:    memory,
		Timer:     memory,
		TimeSlice: memory,
	}
}

func (m *file) write() error {
	serializedData, err := m.memory.Serialize()
	if err != nil {
		return err
	}
	return m.file.Write(serializedData)
}

func (m *file) read() error {
	serializedData := &meta.SerializedData{}
	if err := m.file.Read(serializedData); err != nil {
		return err
	}
	return m.memory.Deserialize(serializedData)
}

func (m *file) SetParameters(parameters ...interface{}) {
	m.file.SetParameters(parameters...)
	m.memory.SetParameters(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case interface {
			meta.Timer
			meta.TimeSlice
			meta.Serializer
			internal.Parameterizer
			internal.Initializer
		}:
			m.memory = p
			m.Timer = p
			m.TimeSlice = p
		case meta.Timer:
			m.Timer = p
		case meta.TimeSlice:
			m.TimeSlice = p
		}
	}
}

func (m *file) SetUtilities(parameters ...interface{}) {
	m.file.SetUtilities(parameters...)
	m.memory.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *file) Configure(items ...interface{}) error {
	return m.file.Configure(items...)
}

func (m *file) Initialize() error {
	m.Lock()
	defer m.Unlock()
	if err := m.memory.Initialize(); err != nil {
		return err
	}
	if err := m.file.Initialize(); err != nil {
		return err
	}
	if err := m.read(); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (m *file) Shutdown() {
	m.Lock()
	defer m.Unlock()

	if err := m.write(); err != nil {
		m.Error("error while shutting down: %s", err.Error())
	}
}

func (m *file) TimerCreate(ctx context.Context, t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerCreate(ctx, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerUpdate(ctx context.Context, id string, t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerUpdate(ctx, id, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerDelete(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Timer.TimerDelete(ctx, id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}

func (m *file) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerStart(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerStop(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimeSliceCreate(ctx context.Context, t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	timeSlice, err := m.TimeSlice.TimeSliceCreate(ctx, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func (m *file) TimeSliceUpdate(ctx context.Context, id string, t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	timeSlice, err := m.TimeSlice.TimeSliceUpdate(ctx, id, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func (m *file) TimeSliceDelete(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.TimeSlice.TimeSliceDelete(ctx, id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}
