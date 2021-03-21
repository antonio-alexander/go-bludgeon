package bludgeonmetamysql

import (
	"database/sql"
	"errors"
)

//rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func rowsAffected(result sql.Result, errorString string) (err error) {
	var n int64

	//get the number of rows affected
	if n, err = result.RowsAffected(); err != nil {
		return
	}
	//return error string
	if n <= 0 {
		err = errors.New(errorString)
	}

	return
}

//lastInsertID will return the last id after an insert operation or an error if present
func lastInsertID(result sql.Result) (id int64, err error) {
	//get the last inserted id
	id, err = result.LastInsertId()

	return
}
