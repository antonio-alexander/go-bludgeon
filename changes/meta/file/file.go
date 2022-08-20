package file

import (
	"context"
	"os"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	meta "github.com/antonio-alexander/go-bludgeon/changes/meta"
	memory "github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	internal "github.com/antonio-alexander/go-bludgeon/internal"

	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	logger.Logger
	file interface {
		internal.Configurer
		internal_file.File
		internal.Initializer
		internal.Parameterizer
	}
	meta.Serializer
	meta.Change
	meta.Registration
	meta.RegistrationChange
}

func New() interface {
	meta.Change
	meta.Registration
	meta.RegistrationChange
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	memory := memory.New()
	return &file{
		file:               internal_file.New(),
		Serializer:         memory,
		Change:             memory,
		Registration:       memory,
		RegistrationChange: memory,
	}
}

func (m *file) write() error {
	serializedData, err := m.Serialize()
	if err != nil {
		return err
	}
	return m.file.Write(serializedData)
}

func (m *file) SetUtilities(parameters ...interface{}) {
	m.file.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *file) SetParameters(parameters ...interface{}) {
	m.file.SetParameters(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case interface {
			meta.Serializer
			meta.Change
			meta.Registration
			meta.RegistrationChange
		}:
			m.Serializer = p
			m.Change = p
			m.Registration = p
			m.RegistrationChange = p
		case meta.Serializer:
			m.Serializer = p
		case meta.Change:
			m.Change = p
		case meta.Registration:
			m.Registration = p
		case meta.RegistrationChange:
			m.RegistrationChange = p
		}
	}
}

func (m *file) Configure(items ...interface{}) error {
	m.Lock()
	defer m.Unlock()

	var c *internal_file.Configuration
	var envs map[string]string

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *internal_file.Configuration:
			c = v
		}
	}
	if c == nil {
		c = new(internal_file.Configuration)
		c.Default()
		c.FromEnv(envs)
	}
	return m.file.Configure(c)
}

func (m *file) Initialize() error {
	m.Lock()
	defer m.Unlock()

	if err := m.file.Initialize(); err != nil {
		return err
	}
	serializedData := &meta.SerializedData{}
	if err := m.file.Read(serializedData); err != nil && !os.IsNotExist(err) {
		return err
	}
	return m.Deserialize(serializedData)
}

func (m *file) Shutdown() {
	m.Lock()
	defer m.Unlock()
	serializedData, err := m.Serialize()
	if err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
	if err := m.file.Write(serializedData); err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
}

func (m *file) ChangeCreate(ctx context.Context, c data.ChangePartial) (*data.Change, error) {
	m.Lock()
	defer m.Unlock()
	change, err := m.Change.ChangeCreate(ctx, c)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return change, nil
}

func (m *file) ChangesDelete(ctx context.Context, changeIds ...string) error {
	if err := m.Change.ChangesDelete(ctx, changeIds...); err != nil {
		return err
	}
	return m.write()
}

func (m *file) RegistrationUpsert(ctx context.Context, registrationId string) error {
	if err := m.Registration.RegistrationUpsert(ctx, registrationId); err != nil {
		return err
	}
	return m.write()
}

func (m *file) RegistrationDelete(ctx context.Context, registrationId string) error {
	if err := m.Registration.RegistrationDelete(ctx, registrationId); err != nil {
		return err
	}
	return m.write()
}

func (m *file) RegistrationChangeUpsert(ctx context.Context, changeId string) error {
	if err := m.RegistrationChange.RegistrationChangeUpsert(ctx, changeId); err != nil {
		return err
	}
	return m.write()
}

func (m *file) RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) ([]string, error) {
	changeIdsToDelete, err := m.RegistrationChange.RegistrationChangeAcknowledge(ctx, registrationId, changeIds...)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return changeIdsToDelete, nil
}
