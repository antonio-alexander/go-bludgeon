package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"
	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/pkg/errors"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	driver_mysql "github.com/go-sql-driver/mysql" //import for driver support
)

const (
	tableEmployees   string = "employees"
	tableEmployeesV1 string = "employees_v1"
	lastUpdatedBy    string = "bludgeon_meta_mysql"
)

type mysql struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	*internal_mysql.DB
}

func New() interface {
	meta.Employee
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	return &mysql{
		DB:     internal_mysql.New(),
		Logger: logger.NewNullLogger(),
	}
}

func (m *mysql) SetParameters(parameters ...interface{}) {
	m.DB.SetParameters(parameters...)
}

func (m *mysql) SetUtilities(parameters ...interface{}) {
	m.DB.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *mysql) EmployeeCreate(ctx context.Context, employeePartial data.EmployeePartial) (*data.Employee, error) {
	var args []interface{}
	var columns []string
	var values []string

	if firstName := employeePartial.FirstName; firstName != nil {
		args = append(args, firstName)
		values = append(values, "?")
		columns = append(columns, "first_name")
	}
	if lastName := employeePartial.LastName; lastName != nil {
		args = append(args, lastName)
		values = append(values, "?")
		columns = append(columns, "last_name")
	}
	if emailAddress := employeePartial.EmailAddress; emailAddress != nil {
		args = append(args, emailAddress)
		values = append(values, "?")
		columns = append(columns, "email_address")
	}
	args = append(args, lastUpdatedBy)
	values = append(values, "?")
	columns = append(columns, "last_updated_by")
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", tableEmployees, strings.Join(columns, ","), strings.Join(values, ","))
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		switch err := err.(type) {
		default:
			return nil, err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return nil, err
			case 1364:
				return nil, meta.ErrEmployeeConflictCreate
			}
		}
	}
	if err := rowsAffected(result, meta.ErrEmployeeNotCreated); err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	employee, err := employeeRead(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *mysql) EmployeeRead(ctx context.Context, id string) (*data.Employee, error) {
	return employeeRead(ctx, m, id)
}

func (m *mysql) EmployeeUpdate(ctx context.Context, id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	var args []interface{}
	var updates []string

	if emailAddress := employeePartial.EmailAddress; emailAddress != nil {
		args = append(args, emailAddress)
		updates = append(updates, "email_address = ?")
	}
	if firstName := employeePartial.FirstName; firstName != nil {
		args = append(args, firstName)
		updates = append(updates, "first_name = ?")
	}
	if lastName := employeePartial.LastName; lastName != nil {
		args = append(args, lastName)
		updates = append(updates, "last_name = ?")
	}
	if len(updates) <= 0 || len(args) <= 0 {
		return nil, meta.ErrEmployeeNotUpdated
	}
	args = append(args, lastUpdatedBy)
	updates = append(updates, "last_updated_by = ?")
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=?;", tableEmployees, strings.Join(updates, ","))
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		switch err := err.(type) {
		default:
			return nil, err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return nil, err
			case 1364:
				return nil, meta.ErrEmployeeConflictUpdate
			}
		}
	}
	if err := rowsAffected(result, meta.ErrEmployeeNotUpdated); err != nil {
		return nil, err
	}
	employee, err := employeeRead(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *mysql) EmployeeDelete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableEmployees)
	result, err := m.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return rowsAffected(result, meta.ErrEmployeeNotFound)
}

func (m *mysql) EmployeesRead(ctx context.Context, search data.EmployeeSearch) ([]*data.Employee, error) {
	var searchParameters []string
	var args []interface{}
	var query string

	if ids := search.IDs; len(ids) > 0 {
		var parameters []string
		for _, id := range ids {
			args = append(args, id)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("employee_id IN(%s)", strings.Join(parameters, ",")))
	}
	switch {
	case search.FirstName != nil:
		searchParameters = append(searchParameters, "first_name = ?")
		args = append(args, search.FirstName)
	case len(search.FirstNames) > 0:
		var parameters []string
		for _, firstName := range search.FirstNames {
			parameters = append(parameters, "?")
			args = append(args, firstName)
		}
		searchParameters = append(searchParameters, fmt.Sprintf("first_name IN(%s)", strings.Join(parameters, ",")))
	}
	switch {
	case search.LastName != nil:
		searchParameters = append(searchParameters, "last_name = ?")
		args = append(args, search.LastName)
	case len(search.LastNames) > 0:
		var parameters []string
		for _, lastName := range search.LastNames {
			parameters = append(parameters, "?")
			args = append(args, lastName)
		}
		searchParameters = append(searchParameters, fmt.Sprintf("last_name IN(%s)", strings.Join(parameters, ",")))
	}
	switch {
	case search.EmailAddress != nil:
		searchParameters = append(searchParameters, "email_address = ?")
		args = append(args, search.EmailAddress)
	case len(search.EmailAddresses) > 0:
		var parameters []string
		for _, emailAddress := range search.EmailAddresses {
			parameters = append(parameters, "?")
			args = append(args, emailAddress)
		}
		searchParameters = append(searchParameters, fmt.Sprintf("first_name IN(%s)", strings.Join(parameters, ",")))
	}
	if len(searchParameters) > 0 {
		query = fmt.Sprintf(`SELECT employee_id, first_name, last_name, email_address,
		version, last_updated, last_updated_by FROM %s WHERE %s`,
			tableEmployeesV1, strings.Join(searchParameters, " AND "))
	} else {
		query = fmt.Sprintf(`SELECT employee_id, first_name, last_name, email_address,
		version, last_updated, last_updated_by FROM %s`, tableEmployeesV1)
	}
	rows, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrEmployeeNotFound
		case errors.Is(err, &driver_mysql.MySQLError{}):
			err := err.(*driver_mysql.MySQLError)
			switch err.Number {
			default:
				return nil, err
			}
		}
	}
	defer rows.Close()
	var employees []*data.Employee
	for rows.Next() {
		employee, err := employeeScan(rows.Scan)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}
