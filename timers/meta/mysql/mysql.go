package metamysql

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

type mysql struct {
	sync.RWMutex
	sync.WaitGroup
	*internal_mysql.DB
	logger.Logger
}

func New(parameters ...interface{}) interface {
	internal_mysql.Owner
	meta.Owner
	meta.Timer
	meta.TimeSlice
} {
	m := &mysql{
		DB: internal_mysql.New(parameters...),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *mysql) TimerCreate(timerValues data.TimerPartial) (*data.Timer, error) {
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
		values = append(values, fmt.Sprintf("COALESCE((SELECT id FROM %s WHERE id=?), 0)", tableEmployees))
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

func (m *mysql) TimerUpdate(id string, timerPartial data.TimerPartial) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	timer, err := timerUpdate(tx, id, timerPartial)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *mysql) TimerDelete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableTimers)
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
	timer, err := timerStop(tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerSubmit can be used to stop a timer and set completed to true
func (m *mysql) TimerSubmit(id string, finishTime int64) (*data.Timer, error) {
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := timerStop(tx, id); err != nil {
		return nil, err
	}
	completed := true
	timer, err := timerUpdate(tx, id, data.TimerPartial{
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
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", tableTimeSlices)
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
		searchParameters = append(searchParameters, "id IN(?)")
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
