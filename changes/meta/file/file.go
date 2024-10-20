package file

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"

	common "github.com/antonio-alexander/go-bludgeon/common"
	internal_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	internal_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/pkg/meta/file"
)

type file struct {
	sync.RWMutex
	internal_logger.Logger
	internal_file.File
	file interface {
		common.Configurer
		common.Initializer
		common.Parameterizer
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
	common.Configurer
	common.Initializer
	common.Parameterizer
} {
	memory := memory.New()
	internalFile := internal_file.New()
	return &file{
		file:               internalFile,
		File:               internalFile,
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
	return m.Write(serializedData)
}

func (m *file) SetUtilities(parameters ...interface{}) {
	m.Lock()
	defer m.Unlock()

	for _, p := range parameters {
		switch p := p.(type) {
		case internal_logger.Logger:
			m.Logger = p
		}
	}
	m.file.SetUtilities(parameters...)
}

func (m *file) SetParameters(parameters ...interface{}) {
	m.Lock()
	defer m.Unlock()

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
	switch {
	case m.Serializer == nil:
		panic("serializer not set")
	case m.Change == nil:
		panic("change not set")
	case m.Registration == nil:
		panic("registration not set")
	case m.RegistrationChange == nil:
		panic("registration change not set")
	}
	m.file.SetParameters(parameters...)
}

func (m *file) Configure(items ...interface{}) error {
	m.Lock()
	defer m.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		default:
			c = new(Configuration)
			if err := internal_config.Get(item, configKey, c); err != nil {
				return err
			}
		case internal_config.Envs:
			c = new(Configuration)
			c.Default()
			c.FromEnv(v)
		case *Configuration:
			c = v
		case Configuration:
			c = &v
		}
	}
	if c == nil {
		return errors.New(internal_config.ErrConfigurationNotFound)
	}
	if err := m.file.Configure(c.Configuration); err != nil {
		return err
	}
	return nil
}

func (m *file) Initialize() error {
	m.Lock()
	defer m.Unlock()

	if err := m.file.Initialize(); err != nil {
		return err
	}
	serializedData := &meta.SerializedData{}
	if err := m.Read(serializedData); err != nil && !os.IsNotExist(err) {
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
	if err := m.Write(serializedData); err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
}

func (m *file) ChangeCreate(ctx context.Context, c data.ChangePartial) (*data.Change, error) {
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
