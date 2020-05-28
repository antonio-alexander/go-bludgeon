package bludgeondatabase

//--------------------------------------------------------------------------------------------
// api.go
//--------------------------------------------------------------------------------------------

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"

	"github.com/go-sql-driver/mysql"
)

//CreateTableEmployee will create a table for employees if it doesn't already exist
func (d *Database) CreateTableEmployee() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "sqlite":
		query = `CREATE TABLE IF NOT EXISTS ` + TableEmployee + ` (
			id INTEGER PRIMARY KEY,
			firstname TEXT,
			lastname TEXT
			);`
	case "mysql":
		query = `CREATE TABLE IF NOT EXISTS ` + TableEmployee + ` (
					id BIGINT NOT NULL AUTO_INCREMENT,
					firstname TEXT,
					lastname TEXT,
				
					PRIMARY KEY (id)
				)ENGINE=InnoDB;`
	//case "postgres":
	default:
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}
	err = d.queryNoResult(query)

	return
}


//CreateTableClient will create a table for client if it doesn't already exist
func (d *Database) CreateTableClient() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	//case "sqlite":
	//case "postgres":
	case "mysql":
		query = `CREATE TABLE IF NOT EXISTS ` + TableClient + ` (
					id BIGINT NOT NULL AUTO_INCREMENT,
					name TEXT,
					rate FLOAT,
				
					PRIMARY KEY (id)
				)ENGINE=InnoDB;`
	default:
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//CreateTableProject will create a table for projects if it doesn't already exist
func (d *Database) CreateTableProject() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	//case "sqlite":
	//case "postgres":
	case "mysql":
		query = `CREATE TABLE IF NOT EXISTS ` + TableProject + ` (
					id BIGINT NOT NULL AUTO_INCREMENT,
					client_id BIGINT NOT NULL,
					description TEXT,
					INDEX(id),
				
					PRIMARY KEY (id),
					FOREIGN KEY (client_id)
						REFERENCES client(id)
				)ENGINE=InnoDB;`
	default:
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//DropTableEmployee will drop the employee table if it exists
func (d *Database) DropTableEmployee() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql", "sqlite":
		query = `DROP TABLE IF EXISTS ` + TableEmployee + `;`
	default: //"postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//DropTableTimer will drop the timer table if it exists
func (d *Database) DropTableTimer() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DROP TABLE IF EXISTS ` + TableTimer + `;`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//DropTableClient will drop the client table if it exists
func (d *Database) DropTableClient() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DROP TABLE IF EXISTS ` + TableClient + `;`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//DropTableProject will drop the project table if it exists
func (d *Database) DropTableProject() (err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DROP TABLE IF EXISTS ` + TableProject + `;`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	err = d.queryNoResult(query)

	return
}

//EmployeeCreate allows you to write a new employee to the bludgeon database
func (d *Database) EmployeeCreate(employee bludgeon.Employee) (id int64, err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = "INSERT into " + TableEmployee + "(firstname, lastname) VALUES(?, ?)"
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	if result, err = d.queryResult(query, employee.FirstName, employee.LastName); err != nil {
		return
	}
	id, err = d.lastInsertID(result)

	return
}

//EmployeeRead allows you to read a single employee using an id
func (d *Database) EmployeeRead(id int64) (employee bludgeon.Employee, err error) {
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `SELECT * from ` + TableEmployee + ` WHERE id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	row := d.db.QueryRow(query, id)
	err = row.Scan(&employee.ID, &employee.FirstName, &employee.LastName)

	return
}

//EmployeesRead can be used to read all existing employees
func (d *Database) EmployeesRead() (employees []bludgeon.Employee, err error) {
	var rows *sql.Rows
	var query string
	var employee bludgeon.Employee

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `SELECT * from ` + TableEmployee
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	if rows, err = d.db.Query(query); err == nil {
		for rows.Next() {
			if err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName); err != nil {
				return
			}
			employees = append(employees, employee)
		}
		err = rows.Err()
	}

	return
}

//EmployeeUpdate can be used to update all values (except the employeeid) of an existing employee
func (d *Database) EmployeeUpdate(id int64, employee bludgeon.Employee) (err error) {
	var result sql.Result
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `UPDATE ` + TableEmployee + ` SET firstname = ?, lastname = ? where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
		return
	}

	if result, err = d.queryResult(query, employee.FirstName, employee.LastName, id); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrUpdateFailed)

	return
}

//EmployeeDelete can be used to delete an existing employee by using it's id
func (d *Database) EmployeeDelete(id int64) (err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DELETE FROM ` + TableEmployee + ` where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrDeleteFailed)

	return
}




//ProjectCreate can be used to create a new project and return its id
func (d *Database) ProjectCreate(project bludgeon.Project) (id int64, err error) {
	var result sql.Result
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `INSERT into ` + TableProject + `(clientid, description) VALUES(?, ?, ?)`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, project.ClientID, project.Description); err != nil {
		return
	}
	id, err = d.lastInsertID(result)

	return
}

//ProjectRead can be used to read an existing project using its id
func (d *Database) ProjectRead(id int64) (project bludgeon.Project, err error) {
	var row *sql.Row
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `SELECT * from ` + TableProject + ` WHERE id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	row = d.db.QueryRow(query, id)
	err = row.Scan(&project.ID, &project.ClientID, &project.Description)

	return
}

//ProjectsRead can be used to read one or more projects with the same ClientID or all Projects
func (d *Database) ProjectsRead(clientID int64) (projects []bludgeon.Project, err error) {
	var rows *sql.Rows
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		if clientID != 0 {
			query = `SELECT * from ` + TableProject + ` WHERE clientid = ?` + strconv.Itoa(int(clientID))
		} else {
			query = `SELECT * from ` + TableProject
		}
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if rows, err = d.db.Query(query); err == nil {
		for rows.Next() {
			var project bludgeon.Project

			if err = rows.Scan(&project.ID, &project.ClientID, &project.Description); err != nil {
				break
			}
			projects = append(projects, project)
		}
	}

	return
}

//ProjectUpdate can be used to update the values for an existing project using its id
func (d *Database) ProjectUpdate(id int64, project bludgeon.Project) (err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `UPDATE ` + TableProject + ` SET id = ?, unitid = ?, clientid = ?, description = ? where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, project.ID, project.ClientID, project.Description, id); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrUpdateFailed)

	return
}

//ProjectDelete can be used to remove a project using its id
func (d *Database) ProjectDelete(id int64) (err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DELETE FROM ` + TableProject + ` where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, id); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrDeleteFailed)

	return
}

//ClientCreate can be used to create a new Client and return its id
func (d *Database) ClientCreate(c bludgeon.Client) (id int64, err error) {
	var result sql.Result
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `INSERT into ` + TableClient + `(name) VALUES(?)`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, c.Name); err != nil {
		return
	}
	id, err = d.lastInsertID(result)

	return
}

//ClientRead can be used to return an existing client using its id
func (d *Database) ClientRead(id int64) (client bludgeon.Client, err error) {
	var row *sql.Row
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `SELECT * from ` + TableClient + ` WHERE id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	row = d.db.QueryRow(query, id)
	err = row.Scan(&client.ID, &client.Name)

	return
}

//ClientsRead read can be used to read all existing clients
func (d *Database) ClientsRead() (clients []bludgeon.Client, err error) {
	var rows *sql.Rows
	var query string

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `SELECT * from ` + TableClient
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if rows, err = d.db.Query(query); err == nil {
		for rows.Next() {
			var client bludgeon.Client

			if err = rows.Scan(&client.ID, &client.Name); err != nil {
				break
			}
			clients = append(clients, client)
		}
	}

	return
}

//ClientUpdate can be used to udpate an existing client using its id
func (d *Database) ClientUpdate(id int64, client bludgeon.Client) (err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `UPDATE ` + TableClient + ` SET id = ?, name = ? where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, id, client.Name); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrUpdateFailed)

	return
}

//ClientDelete can be used to delete an existing client
func (d *Database) ClientDelete(id int64) (err error) {
	var query string
	var result sql.Result

	//check to see if the pointer is nil, if so, exit immediately
	if d.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}

	switch d.config.Driver {
	case "mysql":
		query = `DELETE FROM ` + TableClient + ` where id = ?`
	default: //, "sqlite", "postgres"
		err = fmt.Errorf(ErrDriverUnsupported, d.config.Driver)
	}

	if result, err = d.queryResult(query, id); err != nil {
		return
	}
	err = d.rowsAffected(result, ErrDeleteFailed)

	return
}
