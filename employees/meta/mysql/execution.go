package mysql

import (
	"database/sql"
	"fmt"

	"github.com/antonio-alexander/go-bludgeon/employees/data"
)

//rowsAffected can be used to return a pre-determined error via errorString in the event
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

func employeeRead(db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
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
	row := db.QueryRow(query, id)
	employee := &data.Employee{}
	firstName, lastName := sql.NullString{}, sql.NullString{}
	if err := row.Scan(
		&employee.ID,
		&firstName,
		&lastName,
		&employee.EmailAddress,
		&employee.Version,
		&employee.LastUpdated,
		&employee.LastUpdatedBy,
	); err != nil {
		return nil, err
	}
	employee.FirstName, employee.LastName = firstName.String, lastName.String
	return employee, nil
}
