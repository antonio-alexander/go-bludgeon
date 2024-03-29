package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
	"github.com/antonio-alexander/go-bludgeon/employees/meta"
)

// rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func rowsAffected(result sql.Result, errIfNoRowsAffected error) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return errIfNoRowsAffected
	}
	return nil
}

func employeeScan(scanFx func(...interface{}) error) (*data.Employee, error) {
	var firstName, lastName sql.NullString
	var lastUpdated sql.NullFloat64

	employee := new(data.Employee)
	if err := scanFx(
		&employee.ID,
		&firstName,
		&lastName,
		&employee.EmailAddress,
		&employee.Version,
		&lastUpdated,
		&employee.LastUpdatedBy,
	); err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrEmployeeNotFound
		}
	}
	employee.FirstName, employee.LastName = firstName.String, lastName.String
	employee.LastUpdated = int64(lastUpdated.Float64 * 1000)
	return employee, nil
}

func employeeRead(ctx context.Context, db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}, id interface{}) (*data.Employee, error) {
	var condition string

	switch id.(type) {
	case string:
		condition = "employee_id = ?"
	case int64:
		condition = fmt.Sprintf("employee_id = (SELECT id FROM %s WHERE aux_id = ?)", tableEmployees)
	}
	query := fmt.Sprintf(`SELECT employee_id, first_name, last_name, email_address,
		version, last_updated, last_updated_by FROM %s WHERE %s;`,
		tableEmployeesV1, condition)
	row := db.QueryRowContext(ctx, query, id)
	employee, err := employeeScan(row.Scan)
	if err != nil {
		return nil, err
	}
	return employee, nil
}
