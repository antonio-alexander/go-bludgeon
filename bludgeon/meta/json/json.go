package bludgeonmetajson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

type metaJSON struct {
	sync.RWMutex                               //mutex for threadsafe functionality
	file         string                        //file to read/write to/from
	timers       map[string]bludgeon.Timer     //map to store timers
	timeSlices   map[string]bludgeon.TimeSlice //active time slices indexed by timer id
}

func NewMetaJSON() interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
} {

	//create internal pointers
	timers := make(map[string]bludgeon.Timer)
	timeSlices := make(map[string]bludgeon.TimeSlice)
	//create metaJSON pointer
	return &metaJSON{
		timers:     timers,
		timeSlices: timeSlices,
	}
}

//write will serialize and write the current in-memory data to
// file
func (m *metaJSON) write() (err error) {
	var bytes []byte

	//marshal serialized data into bytes
	if bytes, err = json.MarshalIndent(&SerializedData{
		Timers:     m.timers,
		TimeSlices: m.timeSlices,
	}, "", " "); err != nil {
		return
	}
	//write bytes to file
	err = ioutil.WriteFile(m.file, bytes, os.ModePerm)

	return
}

func (m *metaJSON) read() (err error) {
	var serializedData SerializedData
	var bytes []byte

	//read data from the file
	if bytes, err = ioutil.ReadFile(m.file); err != nil {
		return
	}
	//unmarshal the bytes into serialized data
	if err = json.Unmarshal(bytes, &serializedData); err != nil {
		return
	}
	//REVIEW: should in-memory take precedence or what's in the file?
	//we don't want to lose data, so we don't replace the maps in wholes, we just
	// copy over the data and let common keys get overwritten
	//store read serialized data
	//store timers
	for timerID, timer := range serializedData.Timers {
		m.timers[timerID] = timer
	}
	//store time slices
	for timeSliceID, timeSlice := range serializedData.TimeSlices {
		m.timeSlices[timeSliceID] = timeSlice
	}

	return
}

//ensure that metaJSON implements Owner
var _ bludgeon.MetaOwner = &metaJSON{}

//Initialize
func (m *metaJSON) Initialize(element interface{}) (err error) {
	m.Lock()
	defer m.Unlock()

	var config Configuration
	var folder string

	//attempt to cast element into configuration
	if config, err = castConfiguration(element); err != nil {
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
func (m *metaJSON) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	//set internal configuration to defaults
	m.file = ""
	//close internal pointers
	//set internal pointers to nil
	m.timers, m.timeSlices = nil, nil

	return
}

//ensure that metaJSON implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaSerialize = &metaJSON{}

//Serialize will attempt to commit current data
func (m *metaJSON) MetaSerialize() (err error) {
	m.Lock()
	defer m.Unlock()

	//commit in-memory data
	err = m.write()

	return
}

//Deserialize will attempt to read current data in-memory
func (m *metaJSON) MetaDeSerialize() (err error) {
	m.Lock()
	defer m.Unlock()

	//attempt to read serialized data from file and
	// incorporate into pointer
	err = m.read()

	return
}

//ensure that metaJSON implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaTimer = &metaJSON{}

//MetaTimerWrite
func (m *metaJSON) MetaTimerWrite(timerID string, timer bludgeon.Timer) (err error) {
	m.Lock()
	defer m.Unlock()

	//store timer into map
	m.timers[timerID] = timer
	//attempt to write data
	err = m.write()

	return
}

//MetaTimerDelete
func (m *metaJSON) MetaTimerDelete(timerID string) (err error) {
	m.Lock()
	defer m.Unlock()

	//store timer into map
	delete(m.timers, timerID)
	//attempt to write data
	err = m.write()

	return
}

//MetaTimerRead
func (m *metaJSON) MetaTimerRead(timerID string) (timer bludgeon.Timer, err error) {
	m.Lock()
	defer m.Unlock()

	var ok bool

	//if timer exists output it
	if timer, ok = m.timers[timerID]; !ok {
		err = fmt.Errorf(ErrTimerNotFoundf, timerID)

		return
	}

	return
}

//ensure that metaJSON implements bludgeon.MetaMetaTimer
var _ bludgeon.MetaTimeSlice = &metaJSON{}

//MetaTimeSliceWrite
func (m *metaJSON) MetaTimeSliceWrite(timeSliceID string, timeSlice bludgeon.TimeSlice) (err error) {
	m.Lock()
	defer m.Unlock()

	//store time slice
	m.timeSlices[timeSliceID] = timeSlice
	//write time slice
	err = m.write()

	return
}

//MetaTimeSliceDelete
func (m *metaJSON) MetaTimeSliceDelete(timeSliceID string) (err error) {
	m.Lock()
	defer m.Unlock()

	//delete time slice
	delete(m.timeSlices, timeSliceID)
	//write changes
	err = m.write()

	return
}

//MetaTimeSliceRead
func (m *metaJSON) MetaTimeSliceRead(timeSliceID string) (timeSlice bludgeon.TimeSlice, err error) {
	m.Lock()
	defer m.Unlock()

	var ok bool

	//if timer exists output it
	if timeSlice, ok = m.timeSlices[timeSliceID]; !ok {
		err = fmt.Errorf(ErrTimeSliceNotFoundf, timeSliceID)

		return
	}

	return
}
