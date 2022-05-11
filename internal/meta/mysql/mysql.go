package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

//common constants
const (
	DatabaseIsolation = sql.LevelSerializable
	LogAlias          = "Database"
)

type Owner interface {
	Initialize(config *Configuration) (err error)
	Shutdown()
}

type DB struct {
	*sql.DB                      //pointer to the database
	logger.Logger                //logger
	config        *Configuration //configuration
}

func New(parameters ...interface{}) *DB {
	m := &DB{}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *DB) Initialize(config *Configuration) error {
	if config == nil {
		return errors.New("configuration is nil")
	}
	if err := config.Validate(); err != nil {
		return err
	}
	//EXAMPLE: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// user:password@tcp(localhost:5555)/dbname?charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Username, config.Password, config.Hostname, config.Port, config.Database, config.ParseTime)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	m.DB, m.config = db, config
	return nil
}

func (m *DB) Shutdown() {
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			m.Error("MySQL: ", err)
		}
	}
	m.config.Default()
}
