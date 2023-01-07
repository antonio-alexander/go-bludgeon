package meta

import (
	"context"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// these constants are used to generate employee specific errors
const (
	EmployeeNotFound       string = "employee not found"
	EmployeeNotUpdated     string = "employee not updated"
	EmployeeNotCreated     string = "employee not created, email address not provided"
	EmployeeConflictCreate string = "cannot create employee; email address in use"
	EmployeeConflictUpdate string = "cannot update employee; email address in use"
)

// these are error variables used within the employee meta
var (
	ErrEmployeeNotFound       = errors.NewNotFound(errors.New(EmployeeNotFound))
	ErrEmployeeNotUpdated     = errors.NewNotUpdated(errors.New(EmployeeNotUpdated))
	ErrEmployeeNotCreated     = errors.NewNotCreated(errors.New(EmployeeNotCreated))
	ErrEmployeeConflictCreate = errors.NewConflict(errors.New(EmployeeConflictCreate))
	ErrEmployeeConflictUpdate = errors.NewConflict(errors.New(EmployeeConflictUpdate))
)

// SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Employees map[string]data.Employee `json:"employees"`
}

// Serializer is an interface that can be used to convert the contents of
// meta into a scalar type
type Serializer interface {
	//Serialize can be used to convert all available metadata
	// to a single pointer to be used to serialize to bytes
	Serialize() (*SerializedData, error)

	//Deserialize can be used to provide metadata as a single pointer
	// once it's been deserialized from bytes
	Deserialize(data *SerializedData) error
}

// Employee is an interface that groups functions to interact with one or more
// employees
type Employee interface {
	//EmployeeCreate can be used to create a single Employee
	// the employee email address is required and must be unique
	// at the time of creation
	EmployeeCreate(ctx context.Context, e data.EmployeePartial) (*data.Employee, error)

	//EmployeeRead can be used to read a single employee given a
	// valid id
	EmployeeRead(ctx context.Context, id string) (*data.Employee, error)

	//EmployeeUpdate can be used to update the properties of a given employee
	EmployeeUpdate(ctx context.Context, id string, e data.EmployeePartial) (*data.Employee, error)

	//EmployeeDelete can be used to delete a single employee given a
	// valid id
	EmployeeDelete(ctx context.Context, id string) error

	//EmployeesRead can be used to read one or more employees, given a set of
	// search parameters
	EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error)
}
