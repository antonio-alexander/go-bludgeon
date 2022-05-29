package memory

import (
	"sync"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"

	"github.com/pkg/errors"
)

const lastUpdatedBy string = "bludgeon_meta_memory"

type memory struct {
	sync.RWMutex                             //mutex for threadsafe functionality
	logger.Logger                            //logger
	timers        map[string]*data.Timer     //map to store timers
	timeSlices    map[string]*data.TimeSlice //active time slices indexed by timer id
}

func New(parameters ...interface{}) interface {
	meta.Owner
	meta.Timer
	meta.TimeSlice
	meta.Serializer
} {
	m := &memory{
		timers:     make(map[string]*data.Timer),
		timeSlices: make(map[string]*data.TimeSlice),
	}
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
	return m
}

func (m *memory) validateTimeSlice(t data.TimeSlicePartial, id ...string) error {
	return nil
}

func (m *memory) timerStop(id string) (*data.Timer, error) {
	timer, ok := m.timers[id]
	if !ok {
		return nil, meta.ErrTimerNotFound
	}
	timeSlices, err := m.timeSlicesRead(data.TimeSliceSearch{
		TimerID: &id,
	})
	if err != nil {
		return nil, err
	}
	timer = elapsedTime(timer, timeSlices)
	if timer.ActiveTimeSliceID == "" {
		return copyTimer(timer), nil
	}
	timer.Finish = time.Now().UnixNano()
	timeSlice, err := m.timeSliceUpdate(timer.ActiveTimeSliceID, data.TimeSlicePartial{
		Finish: &timer.Finish,
	})
	if err != nil {
		return nil, err
	}
	timeSlices[len(timeSlices)-1] = timeSlice
	timer = elapsedTime(timer, timeSlices)
	timer.ActiveTimeSliceID = ""
	timer.LastUpdated = time.Now().UnixNano()
	timer.Version++
	return copyTimer(timer), nil
}

func (m *memory) timeSliceCreate(t data.TimeSlicePartial) (*data.TimeSlice, error) {
	if err := m.validateTimeSlice(t); err != nil {
		return nil, err
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	timeSlice := &data.TimeSlice{
		ID: id,
		Audit: data.Audit{
			LastUpdated:   time.Now().UnixNano(),
			LastUpdatedBy: lastUpdatedBy,
			Version:       1,
		},
	}
	switch {
	default:
		timeSlice.TimerID = *t.TimerID
	case t.TimerID == nil || *t.TimerID == "":
		return nil, errors.New("timer id not provided")
	}
	if t.Completed != nil {
		timeSlice.Completed = *t.Completed
	}
	if t.Finish != nil {
		timeSlice.Finish = *t.Finish
	}
	if t.Start != nil {
		timeSlice.Start = *t.Start
	}
	m.timeSlices[id] = timeSlice
	return copyTimeSlice(timeSlice), nil
}

func (m *memory) timeSliceUpdate(id string, t data.TimeSlicePartial) (*data.TimeSlice, error) {
	if err := m.validateTimeSlice(t, id); err != nil {
		return nil, err
	}
	timeSlice, ok := m.timeSlices[id]
	if !ok {
		return nil, meta.ErrTimeSliceNotFound
	}
	if t.Completed != nil {
		timeSlice.Completed = *t.Completed
	}
	if t.Finish != nil {
		timeSlice.Finish = *t.Finish
	}
	if t.Start != nil {
		timeSlice.Start = *t.Start
	}
	timeSlice.LastUpdated = time.Now().UnixNano()
	timeSlice.Version++
	return copyTimeSlice(timeSlice), nil
}

func (m *memory) timeSlicesRead(search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	searchFx := func(t *data.TimeSlice) bool {
		//KIM: this is an inclusive search and is computationally expensive
		if len(search.IDs) > 0 {
			found := false
			for _, id := range search.IDs {
				if t.ID == id {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.TimerID != nil:
			if t.TimerID != *search.TimerID {
				return false
			}
		case len(search.TimerIDs) > 0:
			found := false
			for _, timerID := range search.TimerIDs {
				if t.TimerID == timerID {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if search.Completed != nil && t.Completed != *search.Completed {
			return false
		}
		return true
	}
	var timeSlices []*data.TimeSlice
	for _, timeSlice := range m.timeSlices {
		if searchFx(timeSlice) {
			timeSlices = append(timeSlices, copyTimeSlice(timeSlice))
		}
	}
	return timeSlices, nil
}

func (m *memory) Shutdown() {
	m.Lock()
	defer m.Unlock()
	m.timers = nil
	m.timeSlices = nil
}

func (m *memory) TimeSliceCreate(t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	return m.timeSliceCreate(t)
}

func (m *memory) TimeSliceRead(id string) (*data.TimeSlice, error) {
	m.RLock()
	defer m.RUnlock()
	timeSlice, ok := m.timeSlices[id]
	if !ok {
		return nil, meta.ErrTimeSliceNotFound
	}
	return copyTimeSlice(timeSlice), nil
}

func (m *memory) TimeSliceUpdate(id string, t data.TimeSlicePartial) (*data.TimeSlice, error) {
	m.Lock()
	defer m.Unlock()
	return m.timeSliceUpdate(id, t)
}

func (m *memory) TimeSliceDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.timeSlices[id]; !ok {
		return meta.ErrTimeSliceNotFound
	}
	delete(m.timeSlices, id)
	return nil
}

func (m *memory) TimeSlicesRead(search data.TimeSliceSearch) ([]*data.TimeSlice, error) {
	m.RLock()
	defer m.RUnlock()
	return m.timeSlicesRead(search)
}

func (m *memory) TimerCreate(t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	timer := &data.Timer{
		ID: id,
		Audit: data.Audit{
			LastUpdated:   time.Now().UnixNano(),
			LastUpdatedBy: lastUpdatedBy,
			Version:       1,
		},
	}
	if archived := t.Archived; archived != nil {
		timer.Archived = *archived
	}
	if completed := t.Completed; completed != nil {
		timer.Completed = *completed
	}
	if comment := t.Comment; comment != nil {
		timer.Comment = *comment
	}
	if employeeID := t.EmployeeID; employeeID != nil {
		timer.EmployeeID = *employeeID
	}
	m.timers[timer.ID] = timer
	return copyTimer(timer), nil
}

func (m *memory) TimerRead(id string) (*data.Timer, error) {
	m.RLock()
	defer m.RUnlock()
	timer, ok := m.timers[id]
	if !ok {
		return nil, meta.ErrTimerNotFound
	}
	timeSlices, err := m.timeSlicesRead(data.TimeSliceSearch{
		TimerID: &id,
	})
	if err != nil {
		return nil, err
	}
	timer = elapsedTime(timer, timeSlices)
	return copyTimer(timer), nil
}

func (m *memory) TimerUpdate(id string, t data.TimerPartial) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, ok := m.timers[id]
	if !ok {
		return nil, meta.ErrTimerNotFound
	}
	//REVIEW: should we give an error if nothing was
	// actually updated?
	if archived := t.Archived; archived != nil {
		timer.Archived = *archived
	}
	if completed := t.Completed; completed != nil {
		timer.Completed = *completed
	}
	if comment := t.Comment; comment != nil {
		timer.Comment = *comment
	}
	if employeeID := t.EmployeeID; employeeID != nil {
		timer.EmployeeID = *employeeID
	}
	timer.LastUpdated = time.Now().UnixNano()
	timer.Version++
	return copyTimer(timer), nil
}

func (m *memory) TimerDelete(id string) error {
	m.Lock()
	defer m.Unlock()
	_, ok := m.timers[id]
	if !ok {
		return meta.ErrTimerNotFound
	}
	for _, timeSlice := range m.timeSlices {
		if timeSlice.TimerID == id {
			delete(m.timeSlices, timeSlice.ID)
		}
	}
	delete(m.timers, id)
	return nil
}

func (m *memory) TimersRead(search data.TimerSearch) ([]*data.Timer, error) {
	m.RLock()
	defer m.RUnlock()
	searchFx := func(t *data.Timer) bool {
		//KIM: this is an inclusive search and is computationally expensive
		if len(search.IDs) > 0 {
			found := false
			for _, id := range search.IDs {
				if t.ID == id {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.EmployeeID != nil:
			if t.EmployeeID != *search.EmployeeID {
				return false
			}
		case len(search.EmployeeIDs) > 0:
			found := false
			for _, employeeID := range search.EmployeeIDs {
				if t.EmployeeID == employeeID {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if search.Completed != nil && t.Completed != *search.Completed {
			return false
		}
		if search.Archived != nil && t.Archived != *search.Archived {
			return false
		}
		return true
	}
	var timers []*data.Timer
	for _, timer := range m.timers {
		if searchFx(timer) {
			timers = append(timers, copyTimer(timer))
		}
	}
	return timers, nil
}

func (m *memory) TimerStart(id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, ok := m.timers[id]
	if !ok {
		return nil, meta.ErrTimerNotFound
	}
	timeSlices, err := m.timeSlicesRead(data.TimeSliceSearch{
		TimerID: &id,
	})
	if err != nil {
		return nil, err
	}
	timer = elapsedTime(timer, timeSlices)
	if timer.ActiveTimeSliceID != "" {
		return copyTimer(timer), nil
	}
	start := time.Now().UnixNano()
	if len(timeSlices) <= 0 {
		timer.Start = start
	}
	timeSlice, err := m.timeSliceCreate(data.TimeSlicePartial{
		TimerID: &timer.ID,
		Start:   &start,
	})
	if err != nil {
		return nil, err
	}
	timer.ActiveTimeSliceID = timeSlice.ID
	timer.LastUpdated = time.Now().UnixNano()
	timer.Version++
	return copyTimer(timer), nil
}

func (m *memory) TimerStop(id string) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	return m.timerStop(id)
}

//TimerSubmit can be used to stop a timer and set completed to true
func (m *memory) TimerSubmit(id string, finishTime int64) (*data.Timer, error) {
	m.Lock()
	defer m.Unlock()
	timer, err := m.timerStop(id)
	if err != nil {
		return nil, err
	}
	timer.Completed = true
	timer.Finish = finishTime
	return copyTimer(timer), nil
}

func (m *memory) Serialize() (*meta.SerializedData, error) {
	m.Lock()
	defer m.Unlock()
	serializedData := &meta.SerializedData{
		Timers:     make(map[string]data.Timer),
		TimeSlices: make(map[string]data.TimeSlice),
	}
	for id, timer := range m.timers {
		serializedData.Timers[id] = *timer
	}
	for id, timeslice := range m.timeSlices {
		serializedData.TimeSlices[id] = *timeslice
	}
	return serializedData, nil
}

func (m *memory) Deserialize(serializedData *meta.SerializedData) error {
	m.Lock()
	defer m.Unlock()
	if serializedData == nil {
		return errors.New("serialized data is nil")
	}
	m.timers = make(map[string]*data.Timer)
	for id := range serializedData.Timers {
		timer := serializedData.Timers[id]
		m.timers[id] = copyTimer(&timer)
	}
	m.timeSlices = make(map[string]*data.TimeSlice)
	for id := range serializedData.TimeSlices {
		timeSlice := serializedData.TimeSlices[id]
		m.timeSlices[id] = copyTimeSlice(&timeSlice)
	}
	return nil
}
