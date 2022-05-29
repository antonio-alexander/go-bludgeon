package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

type Owner interface {
	Initialize(config *Configuration) (err error)
}

type File interface {
	Write(item interface{}) error
	Read(item interface{}) error
}

type file struct {
	sync.RWMutex
	*flock.Flock
	logger.Logger
	config *Configuration
}

func New(parameters ...interface{}) interface {
	File
	Owner
} {
	m := &file{
		config: new(Configuration),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
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

//write will serialize and write the current in-memory data to
// File
func (m *file) Write(item interface{}) error {
	m.Lock()
	defer m.Unlock()
	bytes, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.config.File, bytes, os.ModePerm)
}

func (m *file) Read(item interface{}) error {
	m.RLock()
	defer m.RUnlock()
	bytes, err := ioutil.ReadFile(m.config.File)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, item); err != nil {
		return err
	}
	return nil
}

//Initialize
func (m *file) Initialize(config *Configuration) error {
	m.RLock()
	defer m.RUnlock()
	if config == nil {
		return errors.New("config is nil")
	}
	if err := config.Validate(); err != nil {
		return err
	}
	m.config = config
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
