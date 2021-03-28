package metamysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func castConfiguration(element interface{}) (c Configuration, err error) {
	switch v := element.(type) {
	case json.RawMessage:
		err = json.Unmarshal(v, &c)
	case *Configuration:
		c = *v
	case Configuration:
		c = v
	default:
		err = fmt.Errorf("unsupported type: %t", element)
	}

	return
}

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
