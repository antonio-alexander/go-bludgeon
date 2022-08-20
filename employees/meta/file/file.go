package file

import (
	"context"
	"os"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	internal "github.com/antonio-alexander/go-bludgeon/internal"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	logger.Logger
	internal.Configurer
	internal_file.File
	internal.Initializer
	meta.Serializer
	meta.Employee
}

func New() interface {
	meta.Employee
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	memory := memory.New()
	internalFile := internal_file.New()
	return &file{
		Logger:      logger.NewNullLogger(),
		Configurer:  internalFile,
		File:        internalFile,
		Initializer: internalFile,
		Serializer:  memory,
		Employee:    memory,
	}
}

func (m *file) SetParameters(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case interface {
			meta.Serializer
			meta.Employee
		}:
			m.Serializer = p
			m.Employee = p
		case meta.Employee:
			m.Employee = p
		case meta.Serializer:
			m.Serializer = p
		}
	}
}

func (m *file) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *file) write() error {
	serializedData, err := m.Serialize()
	if err != nil {
		return err
	}
	return m.Write(serializedData)
}

// Initialize
func (m *file) Initialize() error {
	m.Lock()
	defer m.Unlock()

	if err := m.Initializer.Initialize(); err != nil {
		return err
	}
	serializedData := &meta.SerializedData{}
	if err := m.Read(serializedData); err != nil && !os.IsNotExist(err) {
		return err
	}
	return m.Deserialize(serializedData)
}

// Shutdown
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
	m.Initializer.Shutdown()
}

func (m *file) EmployeeCreate(ctx context.Context, e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeCreate(ctx, e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *file) EmployeeUpdate(ctx context.Context, id string, e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeUpdate(ctx, id, e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *file) EmployeeDelete(ctx context.Context, id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Employee.EmployeeDelete(ctx, id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}
