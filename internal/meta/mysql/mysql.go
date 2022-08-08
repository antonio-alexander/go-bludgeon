package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

//common constants
const (
	DatabaseIsolation = sql.LevelSerializable
	LogAlias          = "MySQL"
)

type Owner interface {
	Initialize(config *Configuration) (err error)
	Shutdown()
}

type DB struct {
	*sql.DB
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	stopper     chan struct{}
	config      *Configuration
	initialized bool
}

func New(parameters ...interface{}) *DB {
	m := &DB{}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	if m.Logger == nil {
		m.Logger = logger.New()
	}
	return m
}

func (m *DB) launchPing() {
	m.Add(1)
	started := make(chan struct{})
	go func() {
		defer m.Done()

		pingFx := func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), m.config.ConnectTimeout)
			defer cancel()
			if err := m.PingContext(ctx); err != nil {
				m.Error("%s %s", LogAlias, err)
				return false
			}
			m.Debug("%s successfully connected", LogAlias)
			return true
		}
		tConnect := time.NewTicker(time.Second)
		defer tConnect.Stop()
		close(started)
		if pingFx() {
			return
		}
		for {
			select {
			case <-m.stopper:
				return
			case <-tConnect.C:
				if pingFx() {
					return
				}
			}
		}
	}()
	<-started
}

func (m *DB) Initialize(config *Configuration) error {
	m.Lock()
	defer m.Unlock()

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
	m.stopper = make(chan struct{})
	m.DB, m.config = db, config
	m.initialized = true
	m.launchPing()
	return nil
}

func (m *DB) Shutdown() {
	m.Lock()
	defer m.Unlock()

	if !m.initialized {
		return
	}
	close(m.stopper)
	m.Wait()
	if err := m.DB.Close(); err != nil {
		m.Error("%s %s", LogAlias, err)
	}
	m.config.Default()
}
