package logic

import (
	"github.com/antonio-alexander/go-bludgeon/employees/meta"

	"github.com/pkg/errors"
)

const (
	LogAlias                 string = "Logic"
	PanicEmployeeMetaNotSet  string = "employee meta not set"
	PanicChangesclientNotSet string = "changes client not set"
)

var (
	ErrEmployeeIDNotProvided = errors.New("employee id not provided")
)

// Logic is an interface that provides functionality to interact with
// Employee objects
type Logic interface {
	meta.Employee
}
