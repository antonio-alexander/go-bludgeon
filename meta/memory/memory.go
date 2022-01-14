package memory

import (
	"errors"
	"fmt"
	"sync"

	data "github.com/antonio-alexander/go-bludgeon/data"
	meta "github.com/antonio-alexander/go-bludgeon/meta"
)

type memory struct {
	sync.RWMutex                            //mutex for threadsafe functionality
	timers       map[string]*data.Timer     //map to store timers
	timeSlices   map[string]*data.TimeSlice //active time slices indexed by timer id
}

func New() interface {
	meta.Owner
	meta.Serializer
	meta.Timer
	meta.TimeSlice
} {

	return &memory{
		timers:     make(map[string]*data.Timer),
		timeSlices: make(map[string]*data.TimeSlice),
	}
}

func (m *memory) SerializedDataWrite(s meta.SerializedData) {
	m.Lock()
	defer m.Unlock()
	if len(m.timers) > 0 {
		m.timers = make(map[string]*data.Timer)
	}
	if len(m.timeSlices) > 0 {
		m.timeSlices = make(map[string]*data.TimeSlice)
	}
	for id, timer := range s.Timers {
		m.timers[id] = &timer
	}
	for id, timeSlice := range s.TimeSlices {
		m.timeSlices[id] = &timeSlice
	}
}

func (m *memory) SerializedDataRead() meta.SerializedData {
	m.RLock()
	defer m.RUnlock()

	timers := make(map[string]data.Timer)
	for id, timer := range m.timers {
		timers[id] = *timer
	}
	timeSlices := make(map[string]data.TimeSlice)
	for id, timeSlice := range m.timeSlices {
		timeSlices[id] = *timeSlice
	}
	return meta.SerializedData{
		Timers:     timers,
		TimeSlices: timeSlices,
	}
}

func (m *memory) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	m.timers = nil
	m.timeSlices = nil

	return
}

//MetaTimerWrite
func (m *memory) TimerWrite(timerID string, timer data.Timer) (err error) {
	m.Lock()
	defer m.Unlock()

	m.timers[timerID] = &timer
	return
}

//MetaTimerDelete
func (m *memory) TimerDelete(timerID string) (err error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.timers[timerID]; !ok {
		return errors.New(data.ErrBadTimerID)
	}
	delete(m.timers, timerID)
	return
}

//MetaTimerRead
func (m *memory) TimerRead(timerID string) (data.Timer, error) {
	m.Lock()
	defer m.Unlock()

	timer, ok := m.timers[timerID]
	if !ok {
		return data.Timer{}, fmt.Errorf(meta.ErrTimerNotFoundf, timerID)
	}
	return *timer, nil
}

//MetaTimeSliceWrite
func (m *memory) TimeSliceWrite(timeSliceID string, timeSlice data.TimeSlice) error {
	m.Lock()
	defer m.Unlock()

	m.timeSlices[timeSliceID] = &timeSlice
	return nil
}

//MetaTimeSliceDelete
func (m *memory) TimeSliceDelete(timeSliceID string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.timeSlices[timeSliceID]; !ok {
		return errors.New(data.ErrBadTimerID)
	}
	delete(m.timeSlices, timeSliceID)
	return nil
}

//MetaTimeSliceRead
func (m *memory) TimeSliceRead(timeSliceID string) (data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()

	timeSlice, ok := m.timeSlices[timeSliceID]
	if !ok {
		return data.TimeSlice{}, fmt.Errorf(meta.ErrTimeSliceNotFoundf, timeSliceID)
	}
	return *timeSlice, nil
}
