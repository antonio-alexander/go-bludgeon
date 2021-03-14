package bludgeonmetamysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

const timeout time.Duration = 5 * time.Second

type mysql struct {
	sync.RWMutex                   //mutex for threadsafe functionality
	sync.WaitGroup                 //waitgroup to manage goroutines
	started        bool            //whether or not started
	config         Configuration   //configuration
	stopper        chan struct{}   //stopper for go routines
	db             *sql.DB         //pointer to the database
	ctx            context.Context //context
}

func NewMetaMySQL() interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
} {
	//create internal pointers
	//create mysql pointer
	return &mysql{
		stopper: make(chan struct{}),
	}
}

//queryNoResult is used to perform a query and return an error and ignore the result
// it's used in the API to allow "how" queries are run be located in one place and if
// necessary, this code can use a switch case to run differently depending on the type of
// database configured, this will not return a result
func (m *mysql) queryNoResult(query string, v ...interface{}) (err error) {
	var tx *sql.Tx

	//check to see if the pointer is nil, if so, exit immediately
	if m.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}
	//create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//begin the transaction
	if tx, err = m.db.BeginTx(ctx, &sql.TxOptions{Isolation: DatabaseIsolation}); err != nil {
		return
	}
	defer tx.Rollback()
	//if no error starting the transaction, attempt to execute it
	if _, err = tx.ExecContext(ctx, query, v...); err != nil {
		return
	}
	//if no error, commit the changes
	if err = tx.Commit(); err != nil {
		//if there is an error, attempt to rollback the changes
		tx.Rollback()
	}

	return
}

//queryResult is used to perform a query and return an error and ignore the result
// it's used in the API to allow "how" queries are run be located in one place and if
// necessary, this code can use a switch case to run differently depending on the type of
// database configured, this will return a result
func (m *mysql) queryResult(query string, v ...interface{}) (result sql.Result, err error) {
	var tx *sql.Tx

	//check to see if the pointer is nil, if so, exit immediately
	if m.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}
	//create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//begin the transaction
	if tx, err = m.db.BeginTx(ctx, &sql.TxOptions{Isolation: DatabaseIsolation}); err != nil {
		return
	}
	defer tx.Rollback()
	//if no error starting the transaction, attempt to execute it
	if result, err = tx.ExecContext(ctx, query, v...); err != nil {
		return
	}
	//if no error, commit the changes
	if err = tx.Commit(); err != nil {
		//if there is an error, attempt to rollback the changes
		if err = tx.Rollback(); err != nil {
			return
		}
		//return an error if the commit fails and the rollback doesn't fail
		err = fmt.Errorf(ErrQueryFailed, query)
	}

	return
}

//ensure that mysql implements Owner
var _ bludgeon.MetaOwner = &mysql{}

func (m *mysql) Initialize(element interface{}) (err error) {
	m.Lock()
	defer m.Unlock()

	var config Configuration

	//attempt to cast element into configuration
	if config, err = castConfiguration(element); err != nil {
		return
	}
	//connect
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Username, config.Password, config.Hostname, config.Port, config.Database, config.ParseTime)
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//user:password@tcp(localhost:5555)/dbname?charset=utf8
	//create a connection to the database
	if m.db, err = sql.Open("mysql", dataSourceName); err != nil {
		return
	}
	//attempt to ping the database to verify valid connectivity
	err = m.db.Ping()

	return
}

//Close
func (m *mysql) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	close(m.stopper)
	m.Wait()
	//only close if it's nil
	if m.db != nil {
		err = m.db.Close()
	}
	//set internal configuration to defaults
	//set internal pointers to nil

	return
}

//ensure that mysql implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaTimer = &mysql{}

//MetaTimerRead
func (m *mysql) TimerRead(timerUUID string) (timer bludgeon.Timer, err error) {
	m.Lock()
	defer m.Unlock()

	var rows *sql.Rows
	var timers []bludgeon.Timer

	//query rows for timer, this should only return a single element because timerID should be a primary column
	if rows, err = m.db.Query(QueryTimerSelectf, timerUUID); err != nil {
		return
	}
	//range over rows and get data
	for rows.Next() {
		var timer bludgeon.Timer

		//TODO: add completed and archived
		if err = rows.Scan(&timer.UUID, &timer.ActiveSliceUUID, &timer.Start, &timer.Finish, &timer.ElapsedTime, &timer.Comment); err != nil {
			break
		}
		//add timer to timers
		timers = append(timers, timer)
	}
	if err != nil {
		rows.Close()

		return
	}

	//check for errors
	err = rows.Err()
	//check timers
	if len(timers) > 0 {
		timer = timers[0]
	} else {
		err = fmt.Errorf(bludgeon.ErrTimerNotFoundf, timerUUID)
	}

	return
}

//MetaTimerWrite
func (m *mysql) TimerWrite(timerID string, timer bludgeon.Timer) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//upsert the timer
	if result, err = m.queryResult(QueryTimerUpsert, timer.UUID, timer.ActiveSliceUUID, timer.Start, timer.Finish, timer.ElapsedTime, timer.Comment,
		timer.UUID, timer.ActiveSliceUUID, timer.Start, timer.Finish, timer.ElapsedTime, timer.Comment); err != nil {
		return
	}
	err = rowsAffected(result, ErrUpdateFailed)

	return
}

//MetaTimerDelete
func (m *mysql) TimerDelete(timerUUID string) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//delete the timer
	if result, err = m.queryResult(QueryTimerDeletef, timerUUID); err != nil {
		return
	}
	//ensure that rows were affected
	err = rowsAffected(result, ErrDeleteFailed)

	return
}

//ensure that mysql implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaTimeSlice = &mysql{}

//MetaTimeSliceRead
func (m *mysql) TimeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	m.Lock()
	defer m.Unlock()

	var rows *sql.Rows
	var timeSlices []bludgeon.TimeSlice

	//query rows for timer, this should only return a single element because timerID should be a primary column
	if rows, err = m.db.Query(QueryTimeSliceSelectf, timeSliceID); err == nil {
		for rows.Next() {
			var timeSlice bludgeon.TimeSlice

			//TODO: add completed and archived
			if err = rows.Scan(&timeSlice.UUID, &timeSlice.TimerUUID, &timeSlice.Start, &timeSlice.Finish, &timeSlice.ElapsedTime); err != nil {
				break
			}

			//add timer to timers
			timeSlices = append(timeSlices, timeSlice)
		}
	}
	if err != nil {
		rows.Close()

		return
	}
	if err = rows.Err(); err != nil {
		return
	}
	//check timers
	if len(timeSlices) > 0 {
		timeSlice = timeSlices[0]
	} else {
		err = fmt.Errorf(bludgeon.ErrTimeSliceNotFoundf, timeSliceID)
	}

	return
}

//MetaTimeSliceWrite
func (m *mysql) TimeSliceWrite(timeSliceID string, timeSlice bludgeon.TimeSlice) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//upsert the tiem slice
	if result, err = m.queryResult(QueryTimeSliceUpsert, timeSlice.UUID, timeSlice.TimerUUID, timeSlice.Start, timeSlice.Finish, timeSlice.ElapsedTime,
		timeSlice.UUID, timeSlice.Start, timeSlice.Finish, timeSlice.ElapsedTime); err != nil {
		return
	}
	err = rowsAffected(result, ErrUpdateFailed)

	return
}

//MetaTimeSliceDelete
func (m *mysql) TimeSliceDelete(timeSliceID string) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//delete the timer
	if result, err = m.queryResult(QueryTimeSliceDeletef, timeSliceID); err != nil {
		return
	}
	//ensure that rows were affected
	err = rowsAffected(result, ErrDeleteFailed)

	return
}
