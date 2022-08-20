package logic

import (
	"context"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	"github.com/antonio-alexander/go-bludgeon/internal"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

type logic struct {
	logger.Logger
	meta          meta.Employee
	changesClient changesclient.Client
}

// New will generate a new instance of logic that implements
// the interfaces Logic and Owner, from the provided parameters
// we can set the logger and the employee meta (required)
func New() interface {
	Logic
	internal.Parameterizer
} {
	return &logic{
		Logger: logger.NewNullLogger(),
	}
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

// EmployeeCreate can be used to create a single Employee
// the employee email address is required and must be unique
// at the time of creation
func (l *logic) EmployeeCreate(ctx context.Context, e data.EmployeePartial) (*data.Employee, error) {
	employee, err := l.meta.EmployeeCreate(ctx, e)
	if err != nil {
		return nil, err
	}
	l.Debug("%s created employee %s", LogAlias, employee.ID)
	if _, err := l.changesClient.ChangeUpsert(ctx, changesdata.ChangePartial{
		WhenChanged:     &employee.LastUpdated,
		ChangedBy:       &employee.LastUpdatedBy,
		DataId:          &employee.ID,
		DataServiceName: &ServiceName,
		DataType:        &ChangeTypeEmployee,
		DataAction:      &ChangeActionCreate,
		DataVersion:     &employee.Version,
	}); err != nil {
		l.Error("Error while upserting change: %s", err)
	}
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
	if _, err := l.changesClient.ChangeUpsert(ctx, changesdata.ChangePartial{
		WhenChanged:     &employee.LastUpdated,
		ChangedBy:       &employee.LastUpdatedBy,
		DataId:          &employee.ID,
		DataServiceName: &ServiceName,
		DataType:        &ChangeTypeEmployee,
		DataAction:      &ChangeActionUpdate,
		DataVersion:     &employee.Version,
	}); err != nil {
		l.Error("Error while upserting change: %s", err)
	}
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
	if _, err := l.changesClient.ChangeUpsert(ctx, changesdata.ChangePartial{
		DataId:          &employeeId,
		DataServiceName: &ServiceName,
		DataType:        &ChangeTypeEmployee,
		DataAction:      &ChangeActionDelete,
	}); err != nil {
		l.Error("Error while upserting change: %s", err)
	}
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
