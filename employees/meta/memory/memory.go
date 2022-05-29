package memory

import (
	"strings"
	"sync"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"
)

const lastUpdatedBy string = "bludgeon_meta_memory"

type memory struct {
	sync.RWMutex                            //mutex for threadsafe functionality
	logger.Logger                           //logger
	employees     map[string]*data.Employee //map to store employees
}

func New(parameters ...interface{}) interface {
	meta.Owner
	meta.Employee
	meta.Serializer
} {
	m := &memory{
		employees: make(map[string]*data.Employee),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *memory) validateEmployee(e data.EmployeePartial, create bool, ids ...string) error {
	var id string

	if len(ids) > 0 {
		id = strings.ToLower(ids[0])
	}
	if create {
		if e.EmailAddress == nil || *e.EmailAddress == "" {
			return meta.ErrEmployeeNotCreated
		}
	}
	for _, employee := range m.employees {
		if employee.ID == id {
			continue
		}
		if e.EmailAddress != nil && strings.EqualFold(employee.EmailAddress, *e.EmailAddress) {
			if create {
				return meta.ErrEmployeeConflictCreate
			}
			return meta.ErrEmployeeConflictUpdate
		}
	}
	return nil
}

func (m *memory) Shutdown() {
	m.Lock()
	defer m.Unlock()
	m.employees = nil
}

func (m *memory) EmployeeCreate(e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	if err := m.validateEmployee(e, true); err != nil {
		return nil, err
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	employee := &data.Employee{
		ID:           id,
		EmailAddress: *e.EmailAddress, //KIM: validate will ensure that email address isn't empty or nil
		Audit: data.Audit{
			LastUpdated:   time.Now().UnixNano(),
			LastUpdatedBy: lastUpdatedBy,
			Version:       1,
		},
	}
	employee.EmailAddress = *e.EmailAddress
	if e.FirstName != nil {
		employee.FirstName = *e.FirstName
	}
	if e.LastName != nil {
		employee.LastName = *e.LastName
	}
	m.employees[id] = employee
	return copyEmployee(employee), nil
}

func (m *memory) EmployeeRead(id string) (*data.Employee, error) {
	m.RLock()
	defer m.RUnlock()
	employee, ok := m.employees[id]
	if !ok {
		return nil, meta.ErrEmployeeNotFound
	}
	return copyEmployee(employee), nil
}

func (m *memory) EmployeeUpdate(id string, e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	updated := false
	if err := m.validateEmployee(e, false, id); err != nil {
		return nil, err
	}
	employee, ok := m.employees[id]
	if !ok {
		return nil, meta.ErrEmployeeNotFound
	}
	if e.EmailAddress != nil {
		employee.EmailAddress = *e.EmailAddress
		updated = true
	}
	if e.FirstName != nil {
		employee.FirstName = *e.FirstName
		updated = true
	}
	if e.LastName != nil {
		employee.LastName = *e.LastName
		updated = true
	}
	if !updated {
		return nil, meta.ErrEmployeeNotUpdated
	}
	employee.LastUpdated = time.Now().UnixNano()
	employee.Version++
	return copyEmployee(employee), nil
}

func (m *memory) EmployeeDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.employees[id]; !ok {
		return meta.ErrEmployeeNotFound
	}
	delete(m.employees, id)
	return nil
}

func (m *memory) EmployeesRead(search data.EmployeeSearch) ([]*data.Employee, error) {
	m.RLock()
	defer m.RUnlock()
	searchFx := func(e *data.Employee) bool {
		//KIM: this is an inclusive search and is computationally expensive
		if len(search.IDs) > 0 {
			found := false
			for _, id := range search.IDs {
				if e.ID == id {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.FirstName != nil:
			if e.FirstName != *search.FirstName {
				return false
			}
		case len(search.FirstNames) > 0:
			found := false
			for _, firstName := range search.FirstNames {
				if e.FirstName == firstName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.LastName != nil:
			if e.LastName != *search.LastName {
				return false
			}
		case len(search.LastNames) > 0:
			found := false
			for _, lastName := range search.LastNames {
				if e.LastName == lastName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.EmailAddress != nil:
			if e.EmailAddress != *search.EmailAddress {
				return false
			}
		case len(search.EmailAddresses) > 0:
			found := false
			for _, emailAddress := range search.EmailAddresses {
				if e.EmailAddress == emailAddress {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	var employees []*data.Employee
	for _, employee := range m.employees {
		if searchFx(employee) {
			employees = append(employees, copyEmployee(employee))
		}
	}
	return employees, nil
}

func (m *memory) Serialize() (*meta.SerializedData, error) {
	m.Lock()
	defer m.Unlock()
	serializedData := &meta.SerializedData{
		Employees: make(map[string]data.Employee),
	}
	for id, employee := range m.employees {
		serializedData.Employees[id] = *employee
	}
	return serializedData, nil
}

func (m *memory) Deserialize(serializedData *meta.SerializedData) error {
	m.Lock()
	defer m.Unlock()
	if serializedData == nil {
		return errors.New("serialized data is nil")
	}
	m.employees = make(map[string]*data.Employee)
	for id, employee := range serializedData.Employees {
		m.employees[id] = &employee
	}
	return nil
}
