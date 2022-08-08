package file

import (
	"context"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	File interface {
		internal_file.File
		internal_file.Owner
	}
	logger.Logger
	meta.Serializer
	meta.Employee
	internal_meta.Owner
}

func New(parameters ...interface{}) File {
	var c *internal_file.Configuration

	memory := memory.New(parameters)
	file := &file{
		Serializer: memory,
		Employee:   memory,
		File:       internal_file.New(parameters...),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case *internal_file.Configuration:
			c = p
		case logger.Logger:
			file.Logger = p
		}
	}
	if c != nil {
		if err := file.Initialize(c); err != nil {
			panic(err)
		}
	}
	return file
}

func (m *file) write() error {
	serializedData, err := m.Serialize()
	if err != nil {
		return err
	}
	return m.File.Write(serializedData)
}

//Initialize
func (m *file) Initialize(config *internal_file.Configuration) error {
	m.Lock()
	defer m.Unlock()

	if err := m.File.Initialize(config); err != nil {
		return err
	}
	serializedData := &meta.SerializedData{}
	if err := m.File.Read(serializedData); err != nil {
		return err
	}
	return m.Deserialize(serializedData)
}

//Shutdown
func (m *file) Shutdown() {
	m.Lock()
	defer m.Unlock()
	serializedData, err := m.Serialize()
	if err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
	if err := m.File.Write(serializedData); err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
	m.Owner.Shutdown()
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
