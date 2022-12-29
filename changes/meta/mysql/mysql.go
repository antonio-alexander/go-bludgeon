package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	driver_mysql "github.com/go-sql-driver/mysql" //import for driver support
	errors "github.com/pkg/errors"
)

type mysql struct {
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	*internal_mysql.DB
}

func New() interface {
	meta.Change
	meta.Registration
	meta.RegistrationChange
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	return &mysql{DB: internal_mysql.New()}
}

func (m *mysql) SetUtilities(parameters ...interface{}) {
	m.DB.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *mysql) SetParameters(parameters ...interface{}) {
	m.DB.SetParameters(parameters...)
}

func (m *mysql) ChangeCreate(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error) {
	var args []interface{}
	var columns []string
	var values []string

	if dataId := changePartial.DataId; dataId != nil {
		args = append(args, dataId)
		values = append(values, "?")
		columns = append(columns, "data_id")
	}
	if version := changePartial.DataVersion; version != nil {
		args = append(args, version)
		values = append(values, "?")
		columns = append(columns, "version")
	}
	if dataType := changePartial.DataType; dataType != nil {
		args = append(args, dataType)
		values = append(values, "?")
		columns = append(columns, "type")
	}
	if service := changePartial.DataServiceName; service != nil {
		args = append(args, service)
		values = append(values, "?")
		columns = append(columns, "service")
	}
	if action := changePartial.DataAction; action != nil {
		args = append(args, action)
		values = append(values, "?")
		columns = append(columns, "action")
	}
	if whenChanged := changePartial.WhenChanged; whenChanged != nil {
		//KIM: the column type is DATETIME, so providing an int64
		// without conversion won't work
		args = append(args, time.Unix(0, *whenChanged))
		values = append(values, "?")
		columns = append(columns, "when_changed")
	}
	if changedBy := changePartial.ChangedBy; changedBy != nil {
		args = append(args, changedBy)
		values = append(values, "?")
		columns = append(columns, "changed_by")
	}
	tx, err := m.Begin()
	if err != nil {
		return nil, err
	}
	//REVIEW: should we add "ON DUPLICATE DO NOTHING" to this to allow inserting?
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", tableChanges, strings.Join(columns, ","), strings.Join(values, ","))
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		switch err := err.(type) {
		default:
			return nil, err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return nil, err
			case 1364:
				return nil, errors.Wrap(err, meta.ChangeConflictWrite)
			}
		}
	}
	if err := rowsAffected(result, meta.ErrChangeNotWritten); err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	employee, err := changeRead(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *mysql) ChangeRead(ctx context.Context, changeId string) (*data.Change, error) {
	return changeRead(ctx, m, changeId)
}

func (m *mysql) ChangesDelete(ctx context.Context, changeIds ...string) error {
	var parameters []string
	var args []interface{}

	if len(changeIds) == 0 {
		return nil
	}
	for _, changeId := range changeIds {
		parameters = append(parameters, "?")
		args = append(args, changeId)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN(%s)", tableChanges, strings.Join(parameters, ","))
	_, err := m.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysql) ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error) {
	var searchParameters []string
	var args []interface{}
	var query string

	if changeIds := search.ChangeIds; len(changeIds) > 0 {
		var parameters []string
		for _, changeId := range changeIds {
			args = append(args, changeId)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("change_id IN(%s)", strings.Join(parameters, ",")))
	}
	if dataIds := search.DataIds; len(dataIds) > 0 {
		var parameters []string
		for _, dataId := range dataIds {
			args = append(args, dataId)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("data_id IN(%s)", strings.Join(parameters, ",")))
	}
	if dataTypes := search.Types; len(dataTypes) > 0 {
		var parameters []string
		for _, dataType := range dataTypes {
			args = append(args, dataType)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("type IN(%s)", strings.Join(parameters, ",")))
	}
	if actions := search.Actions; len(actions) > 0 {
		var parameters []string
		for _, action := range actions {
			args = append(args, action)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("action IN(%s)", strings.Join(parameters, ",")))
	}
	if serviceNames := search.ServiceNames; len(serviceNames) > 0 {
		var parameters []string
		for _, serviceName := range serviceNames {
			args = append(args, serviceName)
			parameters = append(parameters, "?")
		}
		searchParameters = append(searchParameters, fmt.Sprintf("service IN(%s)", strings.Join(parameters, ",")))
	}
	if search.Since != nil {
		searchParameters = append(searchParameters, "when_changed >= ?")
		args = append(args, search.Since)
	}
	if len(searchParameters) > 0 {
		query = fmt.Sprintf(`SELECT change_id, data_id, version, type,
		service, action, when_changed, changed_by FROM %s WHERE %s`,
			tableChangesV1, strings.Join(searchParameters, " AND "))
	} else {
		query = fmt.Sprintf(`SELECT change_id, data_id, version, type,
		service, action, when_changed, changed_by FROM %s`, tableChangesV1)
	}
	rows, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		default:
			return nil, err
		case err == sql.ErrNoRows:
			return nil, meta.ErrChangeNotFound
		case errors.Is(err, &driver_mysql.MySQLError{}):
			err := err.(*driver_mysql.MySQLError)
			switch err.Number {
			default:
				return nil, err
			case 1364:
				return nil, errors.Wrap(err, meta.ChangeConflictWrite)
			}
		}
	}
	defer rows.Close()
	var changes []*data.Change
	for rows.Next() {
		var action, dataType, service, changedBy sql.NullString
		var whenChanged sql.NullFloat64

		change := &data.Change{}
		if err := rows.Scan(
			&change.Id,
			&change.DataId,
			&change.DataVersion,
			&dataType,
			&service,
			&action,
			&whenChanged,
			&changedBy,
		); err != nil {
			return nil, err
		}
		change.DataType, change.DataServiceName = dataType.String, service.String
		change.ChangedBy, change.WhenChanged = changedBy.String, int64(whenChanged.Float64*1000)
		change.DataAction = action.String
		changes = append(changes, change)
	}
	return changes, nil
}

func (m *mysql) RegistrationUpsert(ctx context.Context, registrationId string) error {
	query := fmt.Sprintf("INSERT INTO %s(id) VALUES(?) ON DUPLICATE KEY UPDATE id=?;", tableRegistrations)
	result, err := m.ExecContext(ctx, query, registrationId, registrationId)
	if err != nil {
		switch err := err.(type) {
		default:
			return err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return err
			case 1364:
				return meta.ErrChangeConflictWrite
			}
		}
	}
	return rowsAffected(result, meta.ErrChangeNotWritten)
}

func (m *mysql) RegistrationDelete(ctx context.Context, registrationId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", tableRegistrations)
	result, err := m.ExecContext(ctx, query, registrationId)
	if err != nil {
		switch err := err.(type) {
		default:
			return err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return err
			case 1364:
				return meta.ErrChangeConflictWrite
			}
		}
	}
	return rowsAffected(result, meta.ErrChangeNotWritten)
}

func (m *mysql) RegistrationChangeUpsert(ctx context.Context, changeId string) error {
	query := fmt.Sprintf("INSERT INTO %s(registration_id, change_id) SELECT DISTINCT id AS registration_id, ? AS change_id FROM %s;",
		tableRegistrationChanges, tableRegistrations)
	if _, err := m.ExecContext(ctx, query, changeId); err != nil {
		switch err := err.(type) {
		default:
			return err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return err
			case 1364:
				return meta.ErrChangeConflictWrite
			}
		}
	}
	//KIM: this could affect no rows if there are no registrations so we shouldn't specicially
	// check to see if rows were affected
	return nil
}

func (m *mysql) RegistrationChangesRead(ctx context.Context, registrationId string) ([]string, error) {
	var changeIds []string

	query := fmt.Sprintf(`SELECT change_id FROM %s WHERE registration_id=?`, tableRegistrationChanges)
	rows, err := m.QueryContext(ctx, query, registrationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var changeId string

		if err := rows.Scan(
			&changeId,
		); err != nil {
			return nil, err
		}
		changeIds = append(changeIds, changeId)
	}
	return changeIds, nil
}

func (m *mysql) RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) ([]string, error) {
	var parameters []string

	args := []interface{}{registrationId}
	for _, changeId := range changeIds {
		args = append(args, changeId)
		parameters = append(parameters, "?")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE registration_id=? AND change_id IN(%s)", tableRegistrationChanges, strings.Join(parameters, ","))
	result, err := m.ExecContext(ctx, query, args...)
	if err != nil {
		switch err := err.(type) {
		default:
			return nil, err
		case *driver_mysql.MySQLError:
			switch err.Number {
			default:
				return nil, err
			case 1364:
				return nil, meta.ErrChangeConflictWrite
			}
		}
	}
	if err := rowsAffected(result, meta.ErrChangeNotWritten); err != nil {
		return nil, err
	}
	query = fmt.Sprintf("SELECT change_id FROM %s LEFT JOIN %s ON change_id WHERE %s.change_id IS NOT NULL;",
		tableChangesV1, tableRegistrationsV1, tableChangesV1)
	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var changeIdsToPrune []string
	for rows.Next() {
		var changeId string

		if err := rows.Scan(&changeId); err != nil {
			return nil, err
		}
		changeIdsToPrune = append(changeIdsToPrune, changeId)
	}
	return changeIdsToPrune, nil
}
