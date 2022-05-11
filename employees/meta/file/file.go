package file

import (
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

type file struct {
	sync.RWMutex
	File interface {
		internal_file.File
		internal_file.Owner
	}
	logger.Logger
	meta.Owner
	meta.Serializer
	meta.Employee
}

func New(parameters ...interface{}) interface {
	internal_file.Owner
	meta.Owner
	meta.Employee
} {
	memory := memory.New(parameters)
	m := &file{
		Owner:      memory,
		Serializer: memory,
		Employee:   memory,
		File:       internal_file.New(parameters...),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
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

func (m *file) EmployeeCreate(e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeCreate(e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *file) EmployeeUpdate(id string, e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeUpdate(id, e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *file) EmployeeDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Employee.EmployeeDelete(id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}
