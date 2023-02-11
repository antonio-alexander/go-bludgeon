package logic

import (
	"context"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	config "github.com/antonio-alexander/go-bludgeon/internal/config"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"
)

type logic struct {
	sync.WaitGroup
	sync.RWMutex
	logger.Logger
	meta          meta.Employee
	changesClient changesclient.Client
	configured    bool
	config        *Configuration
}

// New will generate a new instance of logic that implements
// the interfaces Logic and Owner, from the provided parameters
// we can set the logger and the employee meta (required)
func New() interface {
	Logic
	internal.Parameterizer
	internal.Configurer
	internal.Shutdowner
} {
	return &logic{
		Logger: logger.NewNullLogger(),
	}
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

func (l *logic) SetParameters(parameters ...interface{}) {
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case changesclient.Client:
			l.changesClient = p
		case meta.Employee:
			l.meta = p
		}
	}
	switch {
	case l.meta == nil:
		panic(PanicEmployeeMetaNotSet)
	case l.changesClient == nil:
		panic(PanicChangesclientNotSet)
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

func (l *logic) Shutdown() {
	l.Lock()
	defer l.Unlock()

	l.Wait()
}

// EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (l *logic) EmployeeCreate(ctx context.Context, e data.EmployeePartial) (*data.Employee, error) {
	employee, err := l.meta.EmployeeCreate(ctx, e)
	if err != nil {
		return nil, err
	}
	l.Debug("%s created employee %s", LogAlias, employee.ID)
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &employee.LastUpdated,
		ChangedBy:       &employee.LastUpdatedBy,
		DataId:          &employee.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeEmployee,
		DataAction:      &data.ChangeActionCreate,
		DataVersion:     &employee.Version,
	})
	return employee, nil
}

// EmployeeRead can be used to read a single employee given a
// valid id, logic will ensure that the id is not empty
func (l *logic) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	if id == "" {
		return nil, ErrEmployeeIDNotProvided
	}
	employee, err := l.meta.EmployeeRead(ctx, id)
	if err != nil {
		return nil, err
	}
	l.Debug("%s read employee %s", LogAlias, employee.ID)
	return employee, nil
}

// EmployeeUpdate can be used to update the properties of a given employee
func (l *logic) EmployeeUpdate(ctx context.Context, id string, e data.EmployeePartial) (*data.Employee, error) {
	employee, err := l.meta.EmployeeUpdate(ctx, id, e)
	if err != nil {
		return nil, err
	}
	l.Debug("%s updated employee %s", LogAlias, employee.ID)
	l.changeUpsert(changesdata.ChangePartial{
		WhenChanged:     &employee.LastUpdated,
		ChangedBy:       &employee.LastUpdatedBy,
		DataId:          &employee.ID,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeEmployee,
		DataAction:      &data.ChangeActionUpdate,
		DataVersion:     &employee.Version,
	})
	return employee, nil
}

// EmployeeDelete can be used to delete a single employee given a
// valid id, logic will ensure that the id is not empty
func (l *logic) EmployeeDelete(ctx context.Context, employeeId string) error {
	if employeeId == "" {
		return ErrEmployeeIDNotProvided
	}
	if err := l.meta.EmployeeDelete(ctx, employeeId); err != nil {
		return err
	}
	l.Debug("%s deleted employee %s", LogAlias, employeeId)
	l.changeUpsert(changesdata.ChangePartial{
		DataId:          &employeeId,
		DataServiceName: &data.ServiceName,
		DataType:        &data.ChangeTypeEmployee,
		DataAction:      &data.ChangeActionDelete,
	})
	return nil
}

// EmployeesRead can be used to read one or more employees, given a set of
// search parameters
func (l *logic) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	employees, err := l.meta.EmployeesRead(ctx, search)
	if err != nil {
		return nil, err
	}
	l.Debug("%s read %d employees", LogAlias, len(employees))
	return employees, nil
}
