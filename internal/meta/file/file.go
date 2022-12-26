package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

type File interface {
	Write(item interface{}) error
	Read(item interface{}) error
}

type file struct {
	sync.RWMutex
	*flock.Flock
	logger.Logger
	configured bool
	config     *Configuration
}

func New() interface {
	File
	internal.Configurer
	internal.Initializer
	internal.Parameterizer
} {
	return &file{Logger: logger.NewNullLogger()}
}

func (m *file) Lock() {
	m.RWMutex.Lock()
	if m.config.FileLocking {
		if err := m.Flock.Lock(); err != nil {
			panic(err)
		}
	}
}

func (m *file) Unlock() {
	m.RWMutex.Unlock()
	if m.config.FileLocking {
		if err := m.Flock.Unlock(); err != nil {
			panic(err)
		}
	}
}

func (m *file) RLock() {
	m.RWMutex.RLock()
	if m.config.FileLocking {
		//KIM: this is to simplify the implementation, errors should occur
		// catastrophically, and hence the panic is ok (I think)
		if err := m.Flock.RLock(); err != nil {
			panic(err)
		}
	}
}

func (m *file) RUnlock() {
	m.RWMutex.RUnlock()
	if m.config.FileLocking {
		if err := m.Flock.Unlock(); err != nil {
			panic(err)
		}
	}
}

func (m *file) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (m *file) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *file) Configure(items ...interface{}) error {
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

// write will serialize and write the current in-memory data to
// File
func (m *file) Write(item interface{}) error {
	m.Lock()
	defer m.Unlock()
	bytes, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.config.File, bytes, os.ModePerm)
}

func (m *file) Read(item interface{}) error {
	m.RLock()
	defer m.RUnlock()
	bytes, err := os.ReadFile(m.config.File)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, item); err != nil {
		return err
	}
	return nil
}

// Initialize
func (m *file) Initialize() error {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	if !m.configured {
		return errors.New("not configured")
	}
	if _, err := os.Stat(m.config.File); os.IsNotExist(err) {
		folder := filepath.Dir(m.config.File)
		return os.MkdirAll(folder, os.ModePerm)
	}
	if m.config.FileLocking {
		if m.Flock != nil {
			m.Flock.Close()
		}
		m.Flock = flock.New(m.config.LockFile)
	}
	return nil
}

func (m *file) Shutdown() {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()

	if m.Flock != nil {
		m.Flock.Close()
	}
	m.configured = false
}
