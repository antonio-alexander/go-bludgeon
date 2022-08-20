package file

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	memory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	logger.Logger
	internal_file.Owner
	internal_file.File
	meta.Serializer
	meta.Timer
	meta.TimeSlice
}

func New(parameters ...interface{}) File {
	memory := memory.New(parameters...)
	internalFile := internal_file.New(parameters...)
	m := &file{
		Serializer: memory,
		Timer:      memory,
		TimeSlice:  memory,
		File:       internalFile,
		Owner:      internalFile,
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *file) Write() error {
	serializedData, err := m.Serialize()
	if err != nil {
		return err
	}
	return m.File.Write(serializedData)
}

func (m *file) Read() error {
	serializedData := &meta.SerializedData{}
	if err := m.File.Read(serializedData); err != nil && !os.IsNotExist(err) {
		return err
	}
	return m.Deserialize(serializedData)
}

func (m *file) Initialize(config *internal_file.Configuration) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Owner.Initialize(config); err != nil {
		return err
	}
	if err := m.Read(); err != nil {
		folder := filepath.Dir(config.File)
		return os.MkdirAll(folder, os.ModePerm)
	}
	return nil
}

func (m *file) Shutdown() {
	m.Lock()
	defer m.Unlock()
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
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
	if err := m.Write(); err != nil {
		return err
	}
	return nil
}
