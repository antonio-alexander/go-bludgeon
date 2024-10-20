package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
)

//rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func RowsAffected(result sql.Result, errorString string) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return errors.New(errorString)
	}
	return nil
}
