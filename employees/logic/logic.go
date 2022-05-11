package logic

import (
	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
)

type logic struct {
	logger.Logger
	meta.Employee
}

//New will generate a new instance of logic that implements
// the interfaces Logic and Owner, from the provided parameters
// we can set the logger and the employee meta (required)
func New(parameters ...interface{}) interface {
	Logic
} {
	l := &logic{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case meta.Employee:
			l.Employee = p
		case logger.Logger:
			l.Logger = p
		}
	}
	switch {
	case l.Employee == nil:
		panic(PanicEmployeeMetaNotSet)
	}
	return l
}

//EmployeeRead can be used to read a single employee given a
// valid id, logic will ensure that the id is not empty
func (l *logic) EmployeeRead(id string) (*data.Employee, error) {
	if id == "" {
		return nil, ErrEmployeeIDNotProvided
	}
	return l.Employee.EmployeeRead(id)
}

//EmployeeDelete can be used to delete a single employee given a
// valid id, logic will ensure that the id is not empty
func (l *logic) EmployeeDelete(id string) error {
	if id == "" {
		return ErrEmployeeIDNotProvided
	}
	return l.Employee.EmployeeDelete(id)
}
