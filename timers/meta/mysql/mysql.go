package mysql

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

// query constants
const (
	//REVIEW: figure out why this was originally here
	// tableEmployees    string = "employees"
	tableTimers       string = "timers"
	tableTimeSlices   string = "time_slices"
	tableTimersV1     string = "timers_v1"
	tableTimeSlicesV1 string = "time_slices_v1"
)

type mysql struct {
	sync.RWMutex
	sync.WaitGroup
	*internal_mysql.DB
	logger.Logger
}

// New will instante a concrete implementation of the MySQL
// pointer that implements the meta abstraction for interacting
// with timers, if Configuration provided as a parameter, it will
// attempt to Initialize (and panic on error)
func New() interface {
	meta.Timer
	meta.TimeSlice
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	return &mysql{
		Logger: logger.NewNullLogger(),
		DB:     internal_mysql.New(),
	}
}

func (m *mysql) SetParameters(parameters ...interface{}) {
	m.DB.SetParameters(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case *internal_mysql.DB:
			m.DB = p
		}
	}
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

// TimerCreate can be used to create a timer, although
// all fields are available, the only fields that will
// actually be set are: timer_id and comment
func (m *mysql) TimerCreate(ctx context.Context, timerValues data.TimerPartial) (*data.Timer, error) {
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
		//REVIEW: figure out why this was originally here?
		// values = append(values, fmt.Sprintf("COALESCE((SELECT id FROM %s WHERE id=?), 0)", tableEmployees))
		values = append(values, "?")
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
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	timerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	timer, err := timerRead(ctx, tx, timerID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimerRead can be used to read the current value of a given
// timer, values such as start/finish and elapsed time are
// "calculated" values rather than values that can be set
func (m *mysql) TimerRead(ctx context.Context, id string) (*data.Timer, error) {
	return timerRead(ctx, m, id)
}

// TimerUpdate can be used to update values a given timer
// not associated with timer operations, values such as:
// comment, archived and completed
func (m *mysql) TimerUpdate(ctx context.Context, id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	timer, err := timerUpdate(ctx, tx, id, timerPartial)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimerDelete can be used to delete a timer if it exists
func (m *mysql) TimerDelete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableTimers)
	result, err := m.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return rowsAffected(result, meta.ErrTimerNotFound)
}

// TimersRead can be used to read one or more timers depending
// on search values provided
func (m *mysql) TimersRead(ctx context.Context, search data.TimerSearch) ([]*data.Timer, error) {
	var searchParameters []string
	var timers []*data.Timer
	var args []interface{}
	var query string

	if ids := search.IDs; len(ids) > 0 {
		var parameters []string
		for _, changeId := range search.IDs {
			args = append(args, changeId)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("timer_id IN(%s)", strings.Join(parameters, ",")))
	}
	switch {
	case search.EmployeeID != nil:
		searchParameters = append(searchParameters, "employee_id = ?")
		args = append(args, search.EmployeeID)
	case len(search.EmployeeIDs) > 0:
		var parameters []string
		for _, changeId := range search.EmployeeIDs {
			args = append(args, changeId)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("employee_id IN(%s)", strings.Join(parameters, ",")))
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
	rows, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		timer, err := timerScan(rows.Scan)
		if err != nil {
			return nil, err
		}
		timers = append(timers, timer)
	}
	return timers, nil
}

// TimerStart can be used to start a given timer or do nothing
// if the timer is already started
func (m *mysql) TimerStart(ctx context.Context, id string) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	start := time.Now().UnixNano()
	if _, err = timeSliceCreate(ctx, tx, data.TimeSlicePartial{
		TimerID: &id,
		Start:   &start,
	}); err != nil {
		//KIM: this will fail if an active time slice already exists
		return nil, err
	}
	timer, err := timerRead(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimerStop can be used to stop a given timer or do nothing
// if the timer is not started
func (m *mysql) TimerStop(ctx context.Context, id string) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	timer, err := timerStop(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimerSubmit can be used to stop a timer and set completed to true
func (m *mysql) TimerSubmit(ctx context.Context, id string, finishTime int64) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := timerStop(ctx, tx, id); err != nil {
		return nil, err
	}
	completed := true
	timer, err := timerUpdate(ctx, tx, id, data.TimerPartial{
		Completed: &completed,
	})
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

// TimeSliceCreate can be used to create a single time
// slice
func (m *mysql) TimeSliceCreate(ctx context.Context, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	return timeSliceCreate(ctx, m, timeSlicePartial)
}

// TimeSliceRead can be used to read an existing time slice
func (m *mysql) TimeSliceRead(ctx context.Context, timeSliceID string) (*data.TimeSlice, error) {
	return timeSliceRead(ctx, m, timeSliceID)
}

// TimeSliceUpdate can be used to update an existing time slice
func (m *mysql) TimeSliceUpdate(ctx context.Context, timeSliceID string, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	return timeSliceUpdate(ctx, m, timeSliceID, timeSlicePartial)
}

// TimeSliceDelete can be used to delete an existing time slice
func (m *mysql) TimeSliceDelete(ctx context.Context, timeSliceID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", tableTimeSlices)
	result, err := m.ExecContext(ctx, query, timeSliceID)
	if err != nil {
		return err
	}
	return rowsAffected(result, meta.ErrTimerNotFound)
}

// TimeSlicesRead can be used to read zero or more time slices depending on the
// search criteria
func (m *mysql) TimeSlicesRead(ctx context.Context, search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	var timeSlices []*data.TimeSlice
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
		searchParameters = append(searchParameters, "id IN(?)")
		args = append(args, search.IDs)
	}
	query := fmt.Sprintf(`SELECT time_slice_id, start, finish, completed, elapsed_time, timer_id,
		version, last_updated, last_updated_by FROM %s WHERE %s`,
		tableTimeSlicesV1, strings.Join(searchParameters, " AND "))
	rows, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		timeSlice, err := timeSliceScan(rows.Scan)
		if err != nil {
			return nil, err
		}
		timeSlices = append(timeSlices, timeSlice)
	}
	return timeSlices, nil
}
