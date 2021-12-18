package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
	memory "github.com/antonio-alexander/go-bludgeon/meta/memory"
	"github.com/pkg/errors"
)

type metaMemory interface {
	meta.Owner
	meta.Serializer
	meta.Timer
	meta.TimeSlice
}

type file struct {
	sync.RWMutex        //mutex for threadsafe functionality
	file         string //file to read/write to/from
	metaMemory          //
}

func New() interface {
	Owner
	meta.Owner
	meta.Serializer
	meta.Timer
	meta.TimeSlice
} {

	return &file{
		metaMemory: memory.New(),
	}
}

//write will serialize and write the current in-memory data to
// file
func (m *file) write() error {
	serializedData := m.SerializedDataRead()
	bytes, err := json.MarshalIndent(&serializedData, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.file, bytes, os.ModePerm)
}

func (m *file) read() error {
	var serializedData meta.SerializedData

	bytes, err := ioutil.ReadFile(m.file)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, &serializedData); err != nil {
		return err
	}
	m.SerializedDataWrite(serializedData)
	return nil
}

//Initialize
func (m *file) Initialize(config *Configuration) (err error) {
	m.Lock()
	defer m.Unlock()

	var folder string

	if config == nil {
		return errors.New("config is nil")
	}
	if err = config.Validate(); err != nil {
		return
	}
	//store file
	m.file = config.File
	//attempt to read the file
	if err = m.read(); err != nil {
		//get the folder to create
		folder = filepath.Dir(m.file)
		//attempt to make the folder
		err = os.MkdirAll(folder, os.ModePerm)
	}

	return
}

//Shutdown
func (m *file) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	//set internal configuration to defaults
	m.file = ""
	//close internal pointers
	//set internal pointers to nil
	m.metaMemory.Shutdown()

	return
}

//Serialize will attempt to commit current data
func (m *file) Serialize() (err error) {
	m.Lock()
	defer m.Unlock()

	//commit in-memory data
	err = m.write()

	return
}

//Deserialize will attempt to read current data in-memory
func (m *file) DeSerialize() (err error) {
	m.Lock()
	defer m.Unlock()

	//attempt to read serialized data from file and
	// incorporate into pointer
	err = m.read()

	return
}

//MetaTimerWrite
func (m *file) TimerWrite(timerID string, timer data.Timer) (err error) {
	m.Lock()
	defer m.Unlock()

	if err = m.metaMemory.TimerWrite(timerID, timer); err != nil {
		return
	}
	return m.write()
}

//MetaTimerDelete
func (m *file) TimerDelete(timerID string) (err error) {
	m.Lock()
	defer m.Unlock()

	if err = m.metaMemory.TimerDelete(timerID); err != nil {
		return
	}
	return m.write()
}

//MetaTimeSliceWrite
func (m *file) TimeSliceWrite(timeSliceID string, timeSlice data.TimeSlice) (err error) {
	m.Lock()
	defer m.Unlock()

	if err = m.metaMemory.TimeSliceWrite(timeSliceID, timeSlice); err != nil {
		return
	}
	return m.write()
}

//MetaTimeSliceDelete
func (m *file) TimeSliceDelete(timeSliceID string) (err error) {
	m.Lock()
	defer m.Unlock()

	if err = m.metaMemory.TimeSliceDelete(timeSliceID); err != nil {
		return
	}
	return m.write()
}
