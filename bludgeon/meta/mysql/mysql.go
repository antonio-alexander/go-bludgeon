package bludgeonmetamysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

const timeout time.Duration = 5 * time.Second

type mysql struct {
	sync.RWMutex                 //mutex for threadsafe functionality
	sync.WaitGroup               //waitgroup to manage goroutines
	started        bool          //whether or not started
	config         Configuration //configuration
	stopper        chan struct{} //stopper for go routines
	db             *sql.DB       //pointer to the database
}

func NewMetaMySQL() interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
} {
	config := &Configuration{}
	config.Default()
	//create internal pointers
	//create mysql pointer
	return &mysql{
		stopper: make(chan struct{}),
		config:  *config,
	}
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
	m.RLock()
	defer m.RUnlock()

	//REVIEW: does it make sense to read all of the items within a transaction and then
	// roll it back?

	var row *sql.Row
	var query string
	var tx *sql.Tx

	//start a transaction (to be rolled back), then do the following:
	// (1) Query standalone attributes from the timer table
	// (2) Query any active timeslice
	// (3) Query elapsed time
	if tx, err = m.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	query = fmt.Sprintf(`SELECT timer_uuid, timer_start, timer_finish, timer_comment, 
			timer_archived, timer_billed, timer_completed
		FROM %s WHERE timer_uuid = ?`, TableTimer)
	row = tx.QueryRow(query, timerUUID)
	if err = row.Scan(
		&timer.UUID, &timer.Start, &timer.Finish, &timer.Comment,
		&timer.Archived, &timer.Billed, &timer.Completed,
	); err != nil {
		if err == sql.ErrNoRows {
			err = errors.Errorf(ErrTimerNotFoundf, timerUUID)
		}

		return
	}
	query = fmt.Sprintf(`SELECT slice_uuid FROM %s
			INNER JOIN %s ON %s.slice_id=%s.slice_id
			INNER JOIN %s ON %s.timer_id=%s.timer_id
		WHERE timer_uuid=?`,
		TableSlice,
		TableTimerSliceActive, TableSlice, TableTimerSliceActive,
		TableTimer, TableTimer, TableTimerSliceActive,
	)
	row = tx.QueryRow(query, timer.UUID)
	if err = row.Scan(&timer.ActiveSliceUUID); err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = nil
	}
	query = fmt.Sprintf(`SELECT COALESCE(sum(slice_elapsed_time),0) 
		FROM %s INNER JOIN %s ON %s.timer_id = %s.timer_id WHERE timer_uuid=?`,
		TableSlice, TableTimer, TableSlice, TableTimer)
	row = tx.QueryRow(query, timer.UUID)
	if err = row.Scan(&timer.ElapsedTime); err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = nil
	}
	err = tx.Rollback()

	return
}

//MetaTimerWrite
func (m *mysql) TimerWrite(timerUUID string, timer bludgeon.Timer) (err error) {
	m.RLock()
	defer m.RUnlock()

	//REVIEW: how would we handle a uuid collision here?

	var tx *sql.Tx

	//Start a transaction, then do the following:
	// (1) Attempt to upsert the standalone timer attributes
	// (2) Set or unset the active timeslice as provided
	if tx, err = m.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	query := fmt.Sprintf(`INSERT INTO %s (timer_uuid, timer_start, timer_finish, timer_comment, 
			timer_archived, timer_billed, timer_completed)
		VALUES(?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY 
			UPDATE timer_start=VALUES(timer_start), timer_finish=VALUES(timer_finish), timer_comment=VALUES(timer_comment), 
			timer_billed=VALUES(timer_billed), timer_completed=VALUES(timer_completed), timer_archived=VALUES(timer_archived)`, TableTimer)
	if _, err = tx.Exec(query,
		timerUUID, timer.Start, timer.Finish, timer.Comment,
		timer.Archived, timer.Billed, timer.Completed,
	); err != nil {
		return
	}
	//KIM: No need to check if rows were affected since none of the above may
	// have changed
	if timer.ActiveSliceUUID != "" {
		//REVIEW: insert ignore is dangerous...but this doesn't do a whole lot
		query = fmt.Sprintf(`INSERT IGNORE INTO %s (timer_id, slice_id) values((%s),(%s))`,
			TableTimerSliceActive,
			fmt.Sprintf("SELECT timer_id FROM %s WHERE timer_uuid=\"%s\"", TableTimer, timer.UUID),
			fmt.Sprintf("SELECT slice_id FROM %s WHERE slice_uuid=\"%s\"", TableSlice, timer.ActiveSliceUUID),
		)
		if _, err = tx.Exec(query); err != nil {
			return
		}
	} else {
		query = fmt.Sprintf(`DELETE FROM %s WHERE timer_id=(%s)`,
			TableTimerSliceActive,
			fmt.Sprintf("SELECT timer_id FROM %s WHERE timer_uuid=\"%s\"", TableTimer, timer.UUID),
		)
		if _, err = tx.Exec(query); err != nil {
			return
		}
	}
	err = tx.Commit()

	return
}

//MetaTimerDelete
func (m *mysql) TimerDelete(timerUUID string) (err error) {
	m.RLock()
	defer m.RUnlock()

	var result sql.Result
	var query string
	var tx *sql.Tx

	//Start a transaction and do the following:
	// (1) Delete the timer from the timer table
	// (2) Delete any associated slices
	// (3) Delete
	if tx, err = m.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	query = fmt.Sprintf(`DELETE FROM %s WHERE timer_id=(SELECT timer_id from %s WHERE timer_uuid=?)`,
		TableTimerSliceActive, TableTimer)
	if _, err = tx.Exec(query, timerUUID); err != nil {
		return
	}
	query = fmt.Sprintf(`DELETE FROM %s WHERE timer_id=(SELECT timer_id from %s WHERE timer_uuid=?)`,
		TableSlice, TableTimer)
	if _, err = tx.Exec(query, timerUUID); err != nil {
		return
	}
	query = fmt.Sprintf(`DELETE FROM %s WHERE timer_uuid=?`, TableTimer)
	if result, err = tx.Exec(query, timerUUID); err != nil {
		return
	}
	if err = rowsAffected(result, ErrDeleteFailed); err != nil {
		return
	}
	err = tx.Commit()

	return
}

//ensure that mysql implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaTimeSlice = &mysql{}

//MetaTimeSliceRead
func (m *mysql) TimeSliceRead(timeSliceUUID string) (slice bludgeon.TimeSlice, err error) {
	m.RLock()
	defer m.RUnlock()

	var row *sql.Row
	var query string

	//query the slice attributes, also get teh timer_uuid via an inner join with the timer table
	// because a slice is dependent on a timer, this column can never be NULL (it's also a foreign
	// key)
	query = fmt.Sprintf(`SELECT slice_uuid, timer_uuid, slice_start, slice_finish, slice_archived, COALESCE(slice_elapsed_time,0)
		FROM %s	INNER JOIN %s ON %s.timer_id=%s.timer_id
		WHERE slice_uuid=?`,
		TableSlice, TableTimer, TableSlice, TableTimer)
	row = m.db.QueryRow(query, timeSliceUUID)
	if err = row.Scan(
		&slice.UUID,
		&slice.TimerUUID,
		&slice.Start,
		&slice.Finish,
		&slice.Archived,
		&slice.ElapsedTime,
	); err != nil {
		if err == sql.ErrNoRows {
			err = errors.Errorf(ErrTimeSliceNotFoundf, timeSliceUUID)
		}

		return
	}

	return
}

//MetaTimeSliceWrite
func (m *mysql) TimeSliceWrite(sliceUUID string, slice bludgeon.TimeSlice) (err error) {
	m.RLock()
	defer m.RUnlock()

	var result sql.Result
	var query string

	query = fmt.Sprintf(`INSERT INTO %s (slice_uuid, slice_start, slice_finish, slice_archived, timer_id) 
		VALUES(?, ?, ?, ?, (SELECT timer_id FROM %s WHERE timer_uuid="%s"))
			ON DUPLICATE KEY
		UPDATE slice_start=VALUES(slice_start), slice_finish=(slice_finish), slice_archived=VALUES(slice_archived)`,
		TableSlice, TableTimer, slice.TimerUUID)
	if result, err = m.db.Exec(query, sliceUUID, slice.Start, slice.Finish, slice.Archived); err != nil {
		return
	}
	if err = rowsAffected(result, ErrUpdateFailed); err != nil {
		return
	}

	return
}

//MetaTimeSliceDelete
func (m *mysql) TimeSliceDelete(timeSliceUUID string) (err error) {
	m.RLock()
	defer m.RUnlock()

	var result sql.Result
	var query string
	var tx *sql.Tx

	if tx, err = m.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	query = fmt.Sprintf("DELETE FROM %s WHERE slice_id=(SELECT slice_id FROM %s WHERE slice_uuid=?)", TableTimerSliceActive, TableSlice)
	if _, err = tx.Exec(query, timeSliceUUID); err != nil {
		return
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE slice_uuid = ?", TableSlice)
	if result, err = tx.Exec(query, timeSliceUUID); err != nil {
		return
	}
	if err = rowsAffected(result, ErrDeleteFailed); err != nil {
		return
	}
	err = tx.Commit()

	return
}
