package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	memory "github.com/antonio-alexander/go-bludgeon/meta/memory"

	"github.com/pkg/errors"
)

type file struct {
	sync.RWMutex  //mutex for threadsafe functionality
	logger.Logger //logger
	meta.Owner
	meta.Serializer
	meta.Employee
	meta.Timer
	meta.TimeSlice
	config *Configuration
}

func New(parameters ...interface{}) interface {
	Owner
	meta.Owner
	meta.Timer
	// meta.TimeSlice
	meta.Employee
} {
	memory := memory.New(parameters)
	m := &file{
		Owner:      memory,
		Serializer: memory,
		Employee:   memory,
		Timer:      memory,
		TimeSlice:  memory,
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

//write will serialize and write the current in-memory data to
// file
func (m *file) write() error {
	serializedData, err := m.Serialize()
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(&serializedData, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.config.File, bytes, os.ModePerm)
}

func (m *file) read() error {
	serializedData := &meta.SerializedData{}
	bytes, err := ioutil.ReadFile(m.config.File)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, serializedData); err != nil {
		return err
	}
	return m.Deserialize(serializedData)
}

//Initialize
func (m *file) Initialize(config *Configuration) error {
	m.Lock()
	defer m.Unlock()

	var folder string

	if config == nil {
		return errors.New("config is nil")
	}
	if err := config.Validate(); err != nil {
		return err
	}
	//store file
	m.config = config
	//attempt to read the file
	if err := m.read(); err != nil {
		//get the folder to create
		folder = filepath.Dir(m.config.File)
		//attempt to make the folder
		return os.MkdirAll(folder, os.ModePerm)
	}
	return nil
}

//Shutdown
func (m *file) Shutdown() {
	m.Lock()
	defer m.Unlock()
	if err := m.write(); err != nil {
		m.Error("error while shutting down: %s", err.Error())
		return
	}
	//set internal configuration to defaults
	m.config = nil
	//close internal pointers
	//set internal pointers to nil
	m.Owner.Shutdown()
}

func (m *file) EmployeeCreate(e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeCreate(e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}
func (m *file) EmployeeUpdate(id string, e data.EmployeePartial) (*data.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, err := m.Employee.EmployeeUpdate(id, e)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (m *file) EmployeeDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Employee.EmployeeDelete(id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}

func (m *file) TimerCreate(t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerCreate(t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerUpdate(id string, t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerUpdate(id, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.Timer.TimerDelete(id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}

func (m *file) TimerStart(id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerStart(id)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimerStop(id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.Timer.TimerStop(id)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timer, nil
}

func (m *file) TimeSliceCreate(t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	timeSlice, err := m.TimeSlice.TimeSliceCreate(t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func (m *file) TimeSliceUpdate(id string, t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	timeSlice, err := m.TimeSlice.TimeSliceUpdate(id, t)
	if err != nil {
		return nil, err
	}
	if err := m.write(); err != nil {
		return nil, err
	}
	return timeSlice, nil
}

func (m *file) TimeSliceDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if err := m.TimeSlice.TimeSliceDelete(id); err != nil {
		return err
	}
	if err := m.write(); err != nil {
		return err
	}
	return nil
}
