package bludgeonmetamysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"

	config "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql/config"

	//shadow import for mysql driver support
	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	sync.RWMutex                        //mutex for threadsafe functionality
	sync.WaitGroup                      //waitgroup to manage goroutines
	started        bool                 //whether or not started
	config         config.Configuration //configuration
	stopper        chan struct{}        //stopper for go routines
	chDisconnect   chan struct{}        //disconnect channel
	db             *sql.DB              //pointer to the database
	ctx            context.Context      //context
}

func NewMetaMySQL() interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
} {
	//create internal pointers
	//create mysql pointer
	return &mysql{
		chDisconnect: make(chan struct{}),
	}
}

//Connect will attempt to connect to the databse with the given driver and dataSourceName. If the
// connection is successful, it will attempt to ping the server
func (m *mysql) connect(config config.Configuration) (err error) {

	//
	if m.config.Driver, m.config.DataSource, err = convertConfiguration(m.config); err != nil {
		return
	}
	//create a connection to the database
	if m.db, err = sql.Open(m.config.Driver, m.config.DataSource); err != nil {
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
func (m *mysql) disconnect() (err error) {
	//only close if it's nil
	if m.db != nil {
		//cancel the context??
		err = m.db.Close()
	}

	return
}

func (m *mysql) createTables() (err error) {
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
func (m *mysql) queryNoResult(query string, v ...interface{}) (err error) {
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
func (m *mysql) queryResult(query string, v ...interface{}) (result sql.Result, err error) {
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

//ensure that mysql implements Owner
var _ bludgeon.MetaOwner = &mysql{}

func (m *mysql) Initialize(element interface{}) (err error) {
	m.Lock()
	defer m.Unlock()

	var config config.Configuration

	//attempt to cast element into configuration
	if config, err = castConfiguration(element); err != nil {
		return
	}
	//connect
	err = m.connect(config)

	return
}

//Close
func (m *mysql) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	//attempt to disconnect
	err = m.disconnect()
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

// type Manage interface {
// 	//
// 	Start(config Configuration) (err error)

// 	//
// 	Stop() (err error)
// }

// func (m *mysql) Start(config Configuration) (err error) {
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

// func (m *mysql) Stop() (err error) {
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

// //ensure that mysql implements bludgeon.MetaFunctional
// var _ bludgeon.MetaFunctional = &mysql{}

// type Functional interface {
// 	//LaunchManage
// 	LaunchManage()
// }

// func (m *mysql) LaunchManage() {
// 	started := make(chan struct{})
// 	m.Add(1)
// 	go m.goManage(started)
// 	<-started
// }

// func (m *mysql) goManage(started chan<- struct{}) {
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

// func (m *mysql) manageLogic(connectedIn, tablesCreatedIn, firstCall bool) (connectedOut, tablesCreatedOut bool) {
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
