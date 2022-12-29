package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"

	"github.com/pkg/errors"
)

const secondToNanoSecond float64 = 1000000000

// rowsAffected can be used to return a pre-determined error via errorString in the event
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

func timerScan(scanFx func(...interface{}) error) (*data.Timer, error) {
	var employeeID, activeTimeSliceID sql.NullString

	var start, finish, elapsedTime, lastUpdated sql.NullFloat64

	timer := &data.Timer{}
	if err := scanFx(
		&timer.ID,
		&start,
		&finish,
		&elapsedTime,
		&timer.Comment,
		&timer.Archived,
		&timer.Completed,
		&employeeID,
		&activeTimeSliceID,
		&timer.Version,
		&lastUpdated,
		&timer.LastUpdatedBy,
	); err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrTimerNotFound
		}
	}
	timer.EmployeeID, timer.ActiveTimeSliceID = employeeID.String, activeTimeSliceID.String
	timer.Start, timer.Finish = int64(start.Float64*secondToNanoSecond), int64(finish.Float64*secondToNanoSecond)
	timer.ElapsedTime = int64(elapsedTime.Float64 * secondToNanoSecond)
	timer.LastUpdated = int64(lastUpdated.Float64 * secondToNanoSecond)
	return timer, nil
}

func timeSliceScan(scanFx func(...interface{}) error) (*data.TimeSlice, error) {
	var start, finish, elapsedTime, lastUpdated sql.NullFloat64

	timeSlice := &data.TimeSlice{}
	if err := scanFx(
		&timeSlice.ID,
		&start,
		&finish,
		&timeSlice.Completed,
		&elapsedTime,
		&timeSlice.TimerID,
		&timeSlice.Version,
		&lastUpdated,
		&timeSlice.LastUpdatedBy,
	); err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrTimeSliceNotFound
		}
	}
	timeSlice.Start, timeSlice.Finish = int64(start.Float64*1000), int64(finish.Float64*1000)
	timeSlice.ElapsedTime = int64(elapsedTime.Float64 * 1000)
	timeSlice.LastUpdated = int64(lastUpdated.Float64 * 1000)
	return timeSlice, nil
}

func timeSliceCreate(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	return timeSliceRead(ctx, db, activeSliceID)
}

func timeSliceUpdate(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return timeSliceRead(ctx, db, id)
}

func timeSliceRead(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	row := db.QueryRowContext(ctx, query, id)
	return timeSliceScan(row.Scan)
}

func timerRead(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	row := db.QueryRowContext(ctx, query, id)
	return timerScan(row.Scan)
}

func timerUpdate(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if err := rowsAffected(result, meta.ErrTimerNotFound); err != nil {
		return nil, err
	}
	return timerRead(ctx, db, id)
}

func timerStop(ctx context.Context, db interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}, id string) (*data.Timer, error) {
	timer, err := timerRead(ctx, db, id)
	if err != nil {
		return nil, err
	}
	if timer.ActiveTimeSliceID == "" {
		return timer, nil
	}
	finish := time.Now().UnixNano()
	if _, err := timeSliceUpdate(ctx, db, timer.ActiveTimeSliceID, data.TimeSlicePartial{
		Finish: &finish,
	}); err != nil {
		return nil, err
	}
	return timerRead(ctx, db, id)
}
