package metamysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/data"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/meta"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

type mysql struct {
	sync.RWMutex                  //mutex for threadsafe functionality
	sync.WaitGroup                //waitgroup to manage goroutines
	*sql.DB                       //pointer to the database
	logger.Logger                 //logger
	started        bool           //whether or not started
	config         *Configuration //configuration
	stopper        chan struct{}  //stopper for go routines
}

func New(parameters ...interface{}) interface {
	Owner
	meta.Owner
	meta.Timer
	meta.TimeSlice
	meta.Employee
} {
	m := &mysql{
		stopper: make(chan struct{}),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *mysql) Initialize(config *Configuration) error {
	m.Lock()
	defer m.Unlock()

	if m.started {
		return errors.New(ErrStarted)
	}
	if config == nil {
		return errors.New("configuration is nil")
	}
	if err := config.Validate(); err != nil {
		return err
	}
	//EXAMPLE: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// user:password@tcp(localhost:5555)/dbname?charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Username, config.Password, config.Hostname, config.Port, config.Database, config.ParseTime)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	m.DB, m.config = db, config
	m.started = true
	return nil
}

func (m *mysql) Shutdown() {
	m.Lock()
	defer m.Unlock()

	if !m.started {
		m.Error("MySQL: ", errors.New(ErrNotStarted))
		return
	}
	close(m.stopper)
	m.Wait()
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			m.Error("MySQL: ", errors.New(ErrNotStarted))
		}
	}
	m.config.Default()
	m.started = false
}

func (m *mysql) EmployeeCreate(employeePartial data.EmployeePartial) (*data.Employee, error) {
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
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", tableEmployees, strings.Join(columns, ","), strings.Join(values, ","))
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	employee, err := employeeRead(tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *mysql) EmployeeRead(id string) (*data.Employee, error) {
	return employeeRead(m, id)
}

func (m *mysql) EmployeeUpdate(id string, employeePartial data.EmployeePartial) (*data.Employee, error) {
	var args []interface{}
	var updates []string

	if firstName := employeePartial.FirstName; firstName != nil {
		args = append(args, firstName)
		updates = append(updates, "first_name = ?")
	}
	if lastName := employeePartial.LastName; lastName != nil {
		args = append(args, lastName)
		updates = append(updates, "last_name = ?")
	}
	if len(updates) <= 0 || len(args) <= 0 {
		return nil, errors.New("nothing to update")
	}
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE uuid=?;", tableEmployees, strings.Join(updates, ","))
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	if err := rowsAffected(result, fmt.Sprintf(ErrEmployeeNotFoundf, id)); err != nil {
		return nil, err
	}
	employee, err := employeeRead(tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *mysql) EmployeeDelete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableEmployees)
	result, err := m.Exec(query, id)
	if err != nil {
		return err
	}
	return rowsAffected(result, fmt.Sprintf(ErrEmployeeNotFoundf, id))
}

func (m *mysql) EmployeesRead(search data.EmployeeSearch) ([]*data.Employee, error) {
	var searchParameters []string
	var args []interface{}
	var query string

	if ids := search.IDs; len(ids) > 0 {
		searchParameters = append(searchParameters, "uuid IN(?)")
		args = append(args, ids)
	}
	switch {
	case search.FirstName != nil:
		searchParameters = append(searchParameters, "first_name = ?")
		args = append(args, search.FirstName)
	case len(search.FirstNames) > 0:
		searchParameters = append(searchParameters, "first_name IN(?)")
		args = append(args, search.FirstNames)
	}
	switch {
	case search.LastName != nil:
		searchParameters = append(searchParameters, "last_name = ?")
		args = append(args, search.LastName)
	case len(search.LastNames) > 0:
		searchParameters = append(searchParameters, "last_name IN(?)")
		args = append(args, search.LastNames)
	}
	switch {
	case search.EmailAddress != nil:
		searchParameters = append(searchParameters, "email_address = ?")
		args = append(args, search.EmailAddress)
	case len(search.EmailAddresses) > 0:
		searchParameters = append(searchParameters, "email_address IN(?)")
		args = append(args, search.EmailAddresses)
	}
	if len(searchParameters) > 0 {
		query = fmt.Sprintf(`SELECT employee_id, first_name, last_name, email_address,
		version, last_updated, last_updated_by FROM %s WHERE %s`,
			tableEmployeesV1, strings.Join(searchParameters, " AND "))
	} else {
		query = fmt.Sprintf(`SELECT employee_id, first_name, last_name, email_address,
		version, last_updated, last_updated_by FROM %s`, tableEmployeesV1)
	}
	rows, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []*data.Employee
	for rows.Next() {
		employee := &data.Employee{}
		if err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.EmailAddress,
			&employee.Version,
			&employee.LastUpdated,
			&employee.LastUpdatedBy,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (m *mysql) TimerCreate(timerValues data.TimerPartial) (*data.Timer, error) {
	if !m.started {
		return nil, errors.New(ErrNotStarted)
	}
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	columns := make([]string, 0, 4)
	values := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)
	if employeeID := timerValues.EmployeeID; employeeID != nil {
		columns = append(columns, "employee_id")
		values = append(values, fmt.Sprintf("(SELECT id FROM %s WHERE uuid=?)", tableEmployees))
		args = append(args, employeeID)
	}
	if completed := timerValues.Completed; completed != nil {
		columns = append(columns, "completed")
		values = append(values, "?")
		args = append(args, completed)
	}
	if archived := timerValues.Archived; archived != nil {
		columns = append(columns, "archived")
		values = append(values, "?")
		args = append(args, archived)
	}
	if comment := timerValues.Comment; comment != nil {
		columns = append(columns, "comment")
		values = append(values, "?")
		args = append(args, comment)
	}
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", tableTimers, strings.Join(columns, ","), strings.Join(values, ","))
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	timerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	timer, err := timerRead(tx, timerID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *mysql) TimerRead(id string) (*data.Timer, error) {
	return timerRead(m, id)
}

func (m *mysql) TimerStart(id string) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	start := time.Now().UnixNano()
	if _, err = timeSliceCreate(tx, data.TimeSlicePartial{
		TimerID: &id,
		Start:   &start,
	}); err != nil {
		//KIM: this will fail if an active time slice already exists
		return nil, err
	}
	timer, err := timerRead(tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *mysql) TimerStop(id string) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	timer, err := timerRead(tx, id)
	if err != nil {
		return nil, err
	}
	if timer.ActiveTimeSliceID == "" {
		return timer, nil
	}
	finish := time.Now().UnixNano()
	if _, err := timeSliceUpdate(tx, timer.ActiveTimeSliceID, data.TimeSlicePartial{
		Finish: &finish,
	}); err != nil {
		return nil, err
	}
	if timer, err = timerRead(tx, id); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *mysql) TimerUpdate(id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	var args []interface{}
	var updates []string

	if comment := timerPartial.Comment; comment != nil {
		updates = append(updates, "comment = ?")
		args = append(args, comment)
	}
	if completed := timerPartial.Completed; completed != nil {
		updates = append(updates, "completed = ?")
		args = append(args, completed)
	}
	if employeeID := timerPartial.EmployeeID; employeeID != nil {
		updates = append(updates, "employee_id = ?")
		args = append(args, employeeID)
	}
	if archived := timerPartial.Archived; archived != nil {
		updates = append(updates, "archived = ?")
		args = append(args, archived)
	}
	if len(updates) <= 0 || len(args) <= 0 {
		return nil, errors.New("nothing to update")
	}
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	args = append(args, id)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE uuid = ?;`, tableTimers, strings.Join(updates, ","))
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	if err := rowsAffected(result, fmt.Sprintf(ErrTimerNotFoundf, id)); err != nil {
		return nil, err
	}
	timer, err := timerRead(tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *mysql) TimerDelete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableTimers)
	result, err := m.Exec(query, id)
	if err != nil {
		return err
	}
	return rowsAffected(result, fmt.Sprintf(ErrTimerNotFoundf, id))
}

func (m *mysql) TimersRead(search data.TimerSearch) ([]*data.Timer, error) {
	var searchParameters []string
	var args []interface{}
	var query string

	if ids := search.IDs; len(ids) > 0 {
		searchParameters = append(searchParameters, "timer_id IN(?)")
		args = append(args, ids)
	}
	switch {
	case search.EmployeeID != nil:
		searchParameters = append(searchParameters, "employee_id = ?")
		args = append(args, search.EmployeeID)
	case len(search.EmployeeIDs) > 0:
		searchParameters = append(searchParameters, "employee_id IN(?)")
		args = append(args, search.EmployeeIDs)
	}
	if completed := search.Completed; completed != nil {
		searchParameters = append(searchParameters, "completed = ?")
		args = append(args, completed)
	}
	if archived := search.Archived; archived != nil {
		searchParameters = append(searchParameters, "archived = ?")
		args = append(args, archived)
	}
	if len(searchParameters) > 0 {
		query = fmt.Sprintf(`SELECT timer_id, start, finish, elapsed_time, comment, archived, completed, 
		employee_id, active_time_slice_id, version, last_updated, last_updated_by FROM %s WHERE %s`,
			tableTimersV1, strings.Join(searchParameters, " AND "))
	} else {
		query = fmt.Sprintf(`SELECT timer_id, start, finish, elapsed_time, comment, archived, completed, 
		employee_id, active_time_slice_id, version, last_updated, last_updated_by FROM %s`,
			tableTimersV1)
	}
	rows, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timers []*data.Timer
	for rows.Next() {
		timer := &data.Timer{}
		employeeID, activeTimeSliceID := sql.NullString{}, sql.NullString{}
		start, finish := sql.NullInt64{}, sql.NullInt64{}
		elapsed_time := sql.NullInt64{}
		if err := rows.Scan(
			&timer.ID,
			&start,
			&finish,
			&elapsed_time,
			&timer.Comment,
			&timer.Archived,
			&timer.Completed,
			&employeeID,
			&activeTimeSliceID,
			&timer.Version,
			&timer.LastUpdated,
			&timer.LastUpdatedBy,
		); err != nil {
			return nil, err
		}
		timer.EmployeeID, timer.ActiveTimeSliceID = employeeID.String, activeTimeSliceID.String
		timer.Start, timer.Finish = start.Int64, finish.Int64
		timer.ElapsedTime = elapsed_time.Int64
		timers = append(timers, timer)
	}
	return timers, nil
}

func (m *mysql) TimeSliceCreate(timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	return timeSliceCreate(m, timeSlicePartial)
}

func (m *mysql) TimeSliceRead(timeSliceID string) (*data.TimeSlice, error) {
	return timeSliceRead(m, timeSliceID)
}

func (m *mysql) TimeSliceUpdate(timeSliceID string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	return timeSliceUpdate(m, timeSliceID, timeSlicePartial)
}

func (m *mysql) TimeSliceDelete(timeSliceID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uuid=?", tableTimeSlices)
	result, err := m.Exec(query, timeSliceID)
	if err != nil {
		return err
	}
	return rowsAffected(result, ErrTimeSliceNotFoundf)
}

func (m *mysql) TimeSlicesRead(search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	var searchParameters []string
	var args []interface{}

	if search.Completed != nil {
		searchParameters = append(searchParameters, "completed = ?")
		args = append(args, search.Completed)
	}
	if search.TimerID != nil {
		searchParameters = append(searchParameters, "timer_id = ?")
		args = append(args, search.TimerID)
	}
	if len(search.TimerIDs) > 0 {
		searchParameters = append(searchParameters, "timer_id IN(?)")
		args = append(args, search.TimerIDs)
	}
	if len(search.IDs) > 0 {
		searchParameters = append(searchParameters, "uuid IN(?)")
		args = append(args, search.IDs)
	}
	query := fmt.Sprintf(`SELECT time_slice_id, start, finish, completed, elapsed_time, timer_id,
		version, last_updated, last_updated_by FROM %s WHERE %s`,
		tableTimeSlicesV1, strings.Join(searchParameters, " AND "))
	rows, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timeSlices []*data.TimeSlice
	for rows.Next() {
		timeSlice := &data.TimeSlice{}
		if err := rows.Scan(
			&timeSlice.ID,
			&timeSlice.Start,
			&timeSlice.Finish,
			&timeSlice.Completed,
			&timeSlice.ElapsedTime,
			&timeSlice.TimerID,
			&timeSlice.Version,
			&timeSlice.LastUpdated,
			&timeSlice.LastUpdatedBy,
		); err != nil {
			return nil, err
		}
		timeSlices = append(timeSlices, timeSlice)
	}
	return timeSlices, nil
}
