package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql" //import for driver support
)

//common constants
const (
	DatabaseIsolation = sql.LevelSerializable
	LogAlias          = "MySQL"
)

type DB struct {
	*sql.DB
	sync.RWMutex
	sync.WaitGroup
	logger.Logger
	stopper     chan struct{}
	initialized bool
	configured  bool
	config      *Configuration
}

func New() *DB {
	return &DB{Logger: logger.NewNullLogger()}
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

func (m *DB) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (m *DB) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *DB) Configure(items ...interface{}) error {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()

	var envs map[string]string
	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		c = new(Configuration)
		c.Default()
		c.FromEnv(envs)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	m.config = c
	m.configured = true
	return nil
}

func (m *DB) Initialize() error {
	m.Lock()
	defer m.Unlock()

	if !m.configured {
		return errors.New("not configured")
	}
	//EXAMPLE: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// user:password@tcp(localhost:5555)/dbname?charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		m.config.Username, m.config.Password, m.config.Hostname, m.config.Port, m.config.Database, m.config.ParseTime)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	m.stopper = make(chan struct{})
	m.DB = db
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
	m.configured = false
}
