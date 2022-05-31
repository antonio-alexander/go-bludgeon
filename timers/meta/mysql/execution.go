package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	"github.com/pkg/errors"
)

//rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func rowsAffected(result sql.Result, customErr error) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return customErr
	}
	return nil
}

func timeSliceCreate(db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	var columns, values []string
	var args []interface{}

	if timeSlicePartial.TimerID != nil {
		columns = append(columns, "timer_id")
		values = append(values, "?")
		args = append(args, timeSlicePartial.TimerID)
	}
	if timeSlicePartial.Start != nil {
		columns = append(columns, "start")
		values = append(values, "?")
		args = append(args, time.Unix(0, *timeSlicePartial.Start))
	}
	if timeSlicePartial.Finish != nil {
		columns = append(columns, "finish")
		values = append(values, "?")
		args = append(args, time.Unix(0, *timeSlicePartial.Finish))
	}
	if timeSlicePartial.Completed != nil {
		columns = append(columns, "completed")
		values = append(values, "?")
		args = append(args, timeSlicePartial.Completed)
	}
	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s);`, tableTimeSlices,
		strings.Join(columns, ","), strings.Join(values, ","))
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	activeSliceID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	timeSlice, err := timeSliceRead(db, activeSliceID)
	if err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func timeSliceUpdate(db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}, id interface{}, timeSlicePartial data.TimeSlicePartial) (*data.TimeSlice, error) {
	var args []interface{}
	var updates []string
	var column string

	switch id.(type) {
	case string:
		column = "id"
	case int64:
		column = "aux_id"
	}
	if completed := timeSlicePartial.Completed; completed != nil {
		updates = append(updates, "completed = ?")
		args = append(args, completed)
	}
	if start := timeSlicePartial.Start; start != nil {
		updates = append(updates, "start = ?")
		args = append(args, time.Unix(0, *start))
	}
	if finish := timeSlicePartial.Finish; finish != nil {
		updates = append(updates, "finish = ?")
		args = append(args, time.Unix(0, *finish))
	}
	if len(updates) <= 0 || len(args) <= 0 {
		return nil, errors.New("nothing to update")
	}
	args = append(args, id)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = ?;`, tableTimeSlices, strings.Join(updates, ","), column)
	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	timeSlice, err := timeSliceRead(db, id)
	if err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func timeSliceRead(db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}, id interface{}) (*data.TimeSlice, error) {
	var value string
	switch id.(type) {
	case string:
		value = "?"
	case int64:
		value = fmt.Sprintf("(SELECT id FROM %s WHERE aux_id = ?)", tableTimeSlices)
	}
	query := fmt.Sprintf(`SELECT time_slice_id, start, finish, completed, elapsed_time, timer_id,
		version, last_updated, last_updated_by FROM %s WHERE time_slice_id = %s`,
		tableTimeSlicesV1, value)
	row := db.QueryRow(query, id)
	timeSlice := &data.TimeSlice{}
	start, finish := sql.NullInt64{}, sql.NullInt64{}
	elapsed_time := sql.NullInt64{}
	if err := row.Scan(
		&timeSlice.ID,
		&start,
		&finish,
		&timeSlice.Completed,
		&elapsed_time,
		&timeSlice.TimerID,
		&timeSlice.Version,
		&timeSlice.LastUpdated,
		&timeSlice.LastUpdatedBy,
	); err != nil {
		return nil, err
	}
	timeSlice.Start, timeSlice.Finish = start.Int64, finish.Int64
	timeSlice.ElapsedTime = elapsed_time.Int64
	return timeSlice, nil
}

func timerRead(db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}, id interface{}) (*data.Timer, error) {
	var condition string

	switch id.(type) {
	case string:
		condition = "timer_id = ?"
	case int64:
		condition = fmt.Sprintf("timer_id = (SELECT id FROM %s WHERE aux_id = ?)", tableTimers)
	}
	query := fmt.Sprintf(`SELECT timer_id, start, finish, elapsed_time, comment, archived, completed, 
		employee_id, active_time_slice_id, version, last_updated, last_updated_by FROM %s WHERE %s;`,
		tableTimersV1, condition)
	row := db.QueryRow(query, id)
	timer := &data.Timer{}
	employeeID, activeTimeSliceID := sql.NullString{}, sql.NullString{}
	start, finish := sql.NullInt64{}, sql.NullInt64{}
	elapsed_time := sql.NullInt64{}
	if err := row.Scan(
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
	return timer, nil
}

func timerUpdate(db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}, id string, timerPartial data.TimerPartial) (*data.Timer, error) {
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
	if finish := timerPartial.Finish; finish != nil {
		updates = append(updates, "finish = ?")
		args = append(args, time.Unix(0, *finish))
	}
	if len(updates) <= 0 || len(args) <= 0 {
		return nil, errors.New("nothing to update")
	}
	args = append(args, id)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = ?;`, tableTimers, strings.Join(updates, ","))
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	if err := rowsAffected(result, meta.ErrTimerNotFound); err != nil {
		return nil, err
	}
	timer, err := timerRead(db, id)
	if err != nil {
		return nil, err
	}
	return timer, nil
}

func timerStop(db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}, id string) (*data.Timer, error) {
	timer, err := timerRead(db, id)
	if err != nil {
		return nil, err
	}
	if timer.ActiveTimeSliceID == "" {
		return timer, nil
	}
	finish := time.Now().UnixNano()
	if _, err := timeSliceUpdate(db, timer.ActiveTimeSliceID, data.TimeSlicePartial{
		Finish: &finish,
	}); err != nil {
		return nil, err
	}
	if timer, err = timerRead(db, id); err != nil {
		return nil, err
	}
	return timer, nil
}
