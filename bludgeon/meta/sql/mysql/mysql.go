package bludgeonmetamysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	meta "github.com/antonio-alexander/go-bludgeon/bludgeon/meta"

	"github.com/go-sql-driver/mysql"
)

type metaMySQL struct {
	sync.RWMutex                   //mutex for threadsafe functionality
	sync.WaitGroup                 //waitgroup to manage goroutines
	started        bool            //whether or not started
	config         Configuration   //configuration
	stopper        chan struct{}   //stopper for go routines
	chDisconnect   chan struct{}   //disconnect channel
	db             *sql.DB         //pointer to the database
	ctx            context.Context //context
}

func NewMetaMySQL() interface {
	Owner
	// Manage
	meta.MetaTimer
	meta.MetaTimeSlice
} {
	//create internal pointers
	//create metaMySQL pointer
	return &metaMySQL{
		chDisconnect: make(chan struct{}),
	}
}

func (m *metaMySQL) createTables() (err error) {
	// CreateTableTimer will create a table using a query for the configured driver

	//create timer table
	if err = m.queryNoResult(QueryTimerCreateTable); err != nil {
		return
	}
	//create time slice table
	if err = m.queryNoResult(QueryTimeSliceCreateTable); err != nil {
		return
	}

	return
}

//queryNoResult is used to perform a query and return an error and ignore the result
// it's used in the API to allow "how" queries are run be located in one place and if
// necessary, this code can use a switch case to run differently depending on the type of
// database configured, this will not return a result
func (m *metaMySQL) queryNoResult(query string, v ...interface{}) (err error) {
	//check to see if the pointer is nil, if so, exit immediately
	if m.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}
	//check if transactions enabled
	if m.config.UseTransactions {
		var tx *sql.Tx

		//create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
		defer cancel()
		//begin the transaction
		if tx, err = m.db.BeginTx(ctx, &sql.TxOptions{Isolation: DatabaseIsolation}); err != nil {
			return
		}
		//if no error starting the transaction, attempt to execute it
		if _, err = tx.ExecContext(ctx, query, v...); err != nil {
			return
		}
		//if no error, commit the changes
		if err = tx.Commit(); err != nil {
			//if there is an error, attempt to rollback the changes
			tx.Rollback()
		}
	} else {
		_, err = m.db.Exec(query, v...)
	}
	//TODO: create feedback to handle when to disconnect

	return
}

//queryResult is used to perform a query and return an error and ignore the result
// it's used in the API to allow "how" queries are run be located in one place and if
// necessary, this code can use a switch case to run differently depending on the type of
// database configured, this will return a result
func (m *metaMySQL) queryResult(query string, v ...interface{}) (result sql.Result, err error) {
	//check to see if the pointer is nil, if so, exit immediately
	if m.db == nil {
		err = errors.New(ErrDatabaseNil)
		return
	}
	//check if transactions enabled
	if m.config.UseTransactions {
		var tx *sql.Tx

		//create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
		defer cancel()
		//begin the transaction
		if tx, err = m.db.BeginTx(ctx, &sql.TxOptions{Isolation: DatabaseIsolation}); err != nil {
			return
		}
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
	} else {
		result, err = m.db.Exec(query, v...)
	}
	//TODO: create feedback to handle when to disconnect

	return
}

//Owner
type Owner interface {
	//Close
	Close()

	//
	Connect(config Configuration) (err error)

	//
	Disconnect() (err error)
}

//ensure that metaMySQL implements Owner
var _ Owner = &metaMySQL{}

//Close
func (m *metaMySQL) Close() {
	m.Lock()
	defer m.Unlock()

	//only close if it's nil
	if m.db != nil {
		m.db.Close()
	}
	//set internal configuration to defaults
	//close internal pointers
	close(m.chDisconnect)
	//set internal pointers to nil
	m.stopper, m.chDisconnect = nil, nil

	return
}

//Connect will attempt to connect to the databse with the given driver and dataSourceName. If the
// connection is successful, it will attempt to ping the server
func (m *metaMySQL) Connect(config Configuration) (err error) {
	m.Lock()
	defer m.Unlock()

	//
	if m.config, err = validateConfiguration(config); err != nil {
		return
	}
	if m.config.Driver, m.config.DataSource, err = convertConfiguration(m.config); err != nil {
		return
	}
	//create a connection to the database
	if m.db, err = sql.Open(m.config.Driver, m.config.DataSource); err != nil {
		m.db = nil //set pointer to nil
		return
	}
	//attempt to ping the database to verify valid connectivity
	if err = m.db.Ping(); err != nil {
		return
	}
	//create the tables
	err = m.createTables()

	return
}

//Disconnect will close the connection to the database
func (m *metaMySQL) Disconnect() (err error) {
	//only close if it's nil
	if m.db != nil {
		//cancel the context??
		err = m.db.Close()
	}

	return
}

// type Manage interface {
// 	//
// 	Start(config Configuration) (err error)

// 	//
// 	Stop() (err error)
// }

// func (m *metaMySQL) Start(config Configuration) (err error) {
// 	m.Lock()
// 	defer m.Unlock()

// 	//check if started
// 	if m.started {
// 		return
// 	}
// 	//generate the driver and datasource depending on the configuration
// 	if m.config.Driver, m.config.DataSource, err = convertConfiguration(config); err != nil {
// 		return
// 	}
// 	//start the goRoutines
// 	m.LaunchManage()
// 	//set started to true
// 	m.started = true

// 	return
// }

// func (m *metaMySQL) Stop() (err error) {
// 	m.Lock()
// 	defer m.Unlock()

// 	//check if not started
// 	if !m.started {
// 		return
// 	}
// 	//close stopper
// 	close(m.stopper)
// 	//wait on goRoutines to complete
// 	m.Wait()
// 	//attempt to disconnect the database
// 	err = m.disconnect()
// 	//set started to false
// 	m.started = false

// 	return

// }

// //ensure that metaMySQL implements meta.Functional
// var _ meta.Functional = &metaMySQL{}

// type Functional interface {
// 	//LaunchManage
// 	LaunchManage()
// }

// func (m *metaMySQL) LaunchManage() {
// 	started := make(chan struct{})
// 	m.Add(1)
// 	go m.goManage(started)
// 	<-started
// }

// func (m *metaMySQL) goManage(started chan<- struct{}) {
// 	defer m.Done()

// 	var connected, tablesCreated bool

// 	//create ticker
// 	tConnect := time.NewTicker(10 * time.Second)
// 	//attempt to connect with given configuration
// 	connected, tablesCreated = m.manageLogic(connected, tablesCreated, true)
// 	close(started)
// 	//start the business logic
// 	for {
// 		select {
// 		case <-tConnect.C:
// 			connected, tablesCreated = m.manageLogic(connected, tablesCreated, true)
// 		case <-m.chDisconnect:
// 			//set both to false
// 			connected, tablesCreated = false, false
// 		case <-m.stopper:
// 			return
// 		}
// 	}
// }

// func (m *metaMySQL) manageLogic(connectedIn, tablesCreatedIn, firstCall bool) (connectedOut, tablesCreatedOut bool) {
// 	m.Lock()
// 	defer m.Unlock()

// 	var err error

// 	//store out > in
// 	connectedOut, tablesCreatedOut = connectedIn, tablesCreatedIn
// 	//only attempt to connect if not connected
// 	if !connectedOut || firstCall {
// 		//attempt to connect
// 		if err = m.connect(); err != nil {
// 			fmt.Println(err)
// 		} else {
// 			connectedOut = true
// 		}
// 	}
// 	//only perform
// 	if connectedOut && !tablesCreatedOut {
// 		//attempt to create tables
// 		if err = m.createTables(); err != nil {
// 			fmt.Println(err)
// 		} else {
// 			tablesCreatedOut = true
// 		}
// 	}

// 	return
// }

//ensure that metaMySQL implements meta.MetaTimer
var _ meta.MetaTimer = &metaMySQL{}

//MetaTimerRead
func (m *metaMySQL) MetaTimerRead(timerUUID string) (timer bludgeon.Timer, err error) {
	m.Lock()
	defer m.Unlock()

	var rows *sql.Rows
	var timers []bludgeon.Timer

	//query rows for timer, this should only return a single element because timerID should be a primary column
	if rows, err = m.db.Query(fmt.Sprintf(QueryTimerSelectf, TableTimer, timerUUID)); err != nil {
		return
	}
	//range over rows and get data
	for rows.Next() {
		var timer bludgeon.Timer
		var start, finish mysql.NullTime

		//TODO: add completed and archived
		if err = rows.Scan(&timer.ID, &timer.ActiveSliceID, &timer.UUID, &timer.ActiveSliceUUID, &start, &finish, &timer.ElapsedTime); err != nil {
			break
		}
		//scan for time values
		start.Scan(timer.Start)
		finish.Scan(timer.Finish)
		//add timer to timers
		timers = append(timers, timer)
	}
	//check for errors
	err = rows.Err()
	//check timers
	if len(timers) > 0 {
		timer = timers[0]
	} else {
		//TODO: generate error??
	}

	return
}

//MetaTimerWrite
func (m *metaMySQL) MetaTimerWrite(timerID string, timer bludgeon.Timer) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	if result, err = m.queryResult(fmt.Sprintf(QueryTimerUpsertf, TableTimer), timer.ActiveSliceID, timer.UUID, timer.ActiveSliceUUID, timer.Start, timer.Finish, timer.ElapsedTime,
		timer.ActiveSliceID, timer.UUID, timer.ActiveSliceUUID, timer.Start, timer.Finish, timer.ElapsedTime); err != nil {
		return
	}
	err = rowsAffected(result, ErrUpdateFailed)

	return
}

//MetaTimerDelete
func (m *metaMySQL) MetaTimerDelete(timerUUID string) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//delete the timer
	if result, err = m.queryResult(fmt.Sprintf(QueryTimerDeletef, TableTimer), timerUUID); err != nil {
		return
	}
	//ensure that rows were affected
	err = rowsAffected(result, ErrDeleteFailed)

	return
}

//ensure that metaMySQL implements meta.MetaTimer
var _ meta.MetaTimeSlice = &metaMySQL{}

//MetaTimeSliceRead
func (m *metaMySQL) MetaTimeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	m.Lock()
	defer m.Unlock()

	var rows *sql.Rows
	var timeSlices []bludgeon.TimeSlice

	//query rows for timer, this should only return a single element because timerID should be a primary column
	if rows, err = m.db.Query(fmt.Sprintf(QueryTimeSliceSelectf, TableTimeSlice, timeSliceID)); err == nil {
		for rows.Next() {
			var timeSlice bludgeon.TimeSlice
			var start, finish mysql.NullTime

			//TODO: add completed and archived
			if err = rows.Scan(&timeSlice.ID, &timeSlice.TimerID, &timeSlice.UUID, &timeSlice.TimerUUID, &start, &finish, &timeSlice.ElapsedTime); err != nil {
				break
			}
			//scan for time values
			start.Scan(timeSlice.Start)
			finish.Scan(timeSlice.Finish)
			//add timer to timers
			timeSlices = append(timeSlices, timeSlice)
		}
	}
	err = rows.Err()
	//check timers
	if len(timeSlices) > 0 {
		timeSlice = timeSlices[0]
	} else {
		//TODO: generate error??
	}

	return
}

//MetaTimeSliceWrite
func (m *metaMySQL) MetaTimeSliceWrite(timeSliceID string, timeSlice bludgeon.TimeSlice) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	if result, err = m.queryResult(fmt.Sprintf(QueryTimeSliceUpsertf, TableTimeSlice), timeSlice.TimerID, timeSlice.UUID, timeSlice.TimerUUID, timeSlice.Start, timeSlice.Finish, timeSlice.ElapsedTime,
		timeSlice.TimerID, timeSlice.UUID, timeSlice.TimerUUID, timeSlice.Start, timeSlice.Finish, timeSlice.ElapsedTime); err != nil {
		return
	}
	err = rowsAffected(result, ErrUpdateFailed)

	return
}

//MetaTimeSliceDelete
func (m *metaMySQL) MetaTimeSliceDelete(timeSliceID string) (err error) {
	m.Lock()
	defer m.Unlock()

	var result sql.Result

	//delete the timer
	if result, err = m.queryResult(fmt.Sprintf(QueryTimeSliceDeletef, TableTimeSlice), timeSliceID); err != nil {
		return
	}
	//ensure that rows were affected
	err = rowsAffected(result, ErrDeleteFailed)

	return
}
