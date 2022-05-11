package meta

import (
	"strings"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/internal/errors"
)

//these constants are used to generate employee specific errors
const (
	EmployeeNotFound       string = "employee not found"
	EmployeeNotUpdated     string = "employee not updated"
	EmployeeNotCreated     string = "employee not created, email address not provided"
	EmployeeConflictCreate string = "cannot create employee; email address in use"
	EmployeeConflictUpdate string = "cannot update employee; email address in use"
)

//these are error variables used within the employee meta
var (
	ErrEmployeeNotFound       = errors.NewNotFound(EmployeeNotFound)
	ErrEmployeeNotUpdated     = errors.NewNotupdated(EmployeeNotUpdated)
	ErrEmployeeNotCreated     = errors.NewNotCreated(EmployeeNotCreated)
	ErrEmployeeConflictCreate = errors.NewConflict(EmployeeConflictCreate)
	ErrEmployeeConflictUpdate = errors.NewConflict(EmployeeConflictUpdate)
)

//SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Employees map[string]data.Employee `json:"employees"`
}

//Type is a special string that can be used to identify the type
// of meta
type Type string

const (
	TypeInvalid Type = "invalid"
	TypeFile    Type = "file"
	TypeMySQL   Type = "mysql"
)

//String can be used to convert a meta type into a string
// type
func (m Type) String() string {
	switch m {
	case TypeFile:
		return "file"
	case TypeMySQL:
		return "mysql"
	default:
		return "invalid"
	}
}

//AtoType can be used to attempt to convert a string to a
// meta type
func AtoType(s string) Type {
	switch strings.ToLower(s) {
	case "file":
		return TypeFile
	case "mysql":
		return TypeMySQL
	default:
		return TypeInvalid
	}
}

//Serializer is an interface that can be used to convert the contents of
// meta into a scalar type
type Serializer interface {
	//Serialize can be used to convert all available metadata
	// to a single pointer to be used to serialize to bytes
	Serialize() (*SerializedData, error)

	//Deserialize can be used to provide metadata as a single pointer
	// once it's been deserialized from bytes
	Deserialize(data *SerializedData) error
}

//Owner contains methods that shuold only be used by the constructor of the pointer
// and allows use of functions that can affect the underlying pointer
type Owner interface {
	//Shutdown can be used to "stop" the meta and any asynchronous
	// processes
	Shutdown()
}

//Employee is an interface that groups functions to interact with one or more
// employees
type Employee interface {
	//EmployeeCreate can be used to create a single Employee
	// the employee email address is required and must be unique
	// at the time of creation
	EmployeeCreate(e data.EmployeePartial) (*data.Employee, error)

	//EmployeeRead can be used to read a single employee given a
	// valid id
	EmployeeRead(id string) (*data.Employee, error)

	//EmployeeUpdate can be used to update the properties of a given employee
	EmployeeUpdate(id string, e data.EmployeePartial) (*data.Employee, error)

	//EmployeeDelete can be used to delete a single employee given a
	// valid id
	EmployeeDelete(id string) error

	//EmployeesRead can be used to read one or more employees, given a set of
	// search parameters
	EmployeesRead(search data.EmployeeSearch) ([]*data.Employee, error)
}
