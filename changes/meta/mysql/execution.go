package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
)

// rowsAffected can be used to return a pre-determined error via errorString in the event
// no rows are affected; this function assumes that in the event no error is returned and
// rows were supposed to be affected, an error will be returned
func rowsAffected(result sql.Result, errIfNoRowsAffected error) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n <= 0 {
		return errIfNoRowsAffected
	}
	return nil
}

func changeRead(ctx context.Context, db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}, id interface{}) (*data.Change, error) {
	var action, dataType, service, changedBy sql.NullString
	var whenChanged sql.NullFloat64
	var condition string

	switch id.(type) {
	case string:
		condition = "change_id = ?"
	case int64:
		condition = fmt.Sprintf("change_id = (SELECT id FROM %s WHERE aux_id = ?)", tableChanges)
	}
	query := fmt.Sprintf(`SELECT change_id, data_id, version, type,
		service, action, when_changed, changed_by FROM %s WHERE %s;`,
		tableChangesV1, condition)
	row := db.QueryRowContext(ctx, query, id)
	change := &data.Change{}
	if err := row.Scan(
		&change.Id,
		&change.DataId,
		&change.DataVersion,
		&dataType,
		&service,
		&action,
		&whenChanged,
		&changedBy,
	); err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrChangeNotFound
		}
	}
	change.DataType, change.DataServiceName = dataType.String, service.String
	change.ChangedBy, change.WhenChanged = changedBy.String, int64(whenChanged.Float64*1000)
	change.DataAction = action.String
	return change, nil
}
