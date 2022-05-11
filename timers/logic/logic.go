package logic

import (
	"sync"
	"time"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
)

type logic struct {
	sync.RWMutex
	logger.Logger
	meta.Timer
	started bool
}

type Owner interface {
	Start() (err error)
	Stop() (err error)
}

type Logic interface {
	//meta.Timer
	//KIM: these are implemented by meta already
	TimerCreate(timer data.TimerPartial) (*data.Timer, error)
	TimerRead(id string) (*data.Timer, error)
	TimersRead(search data.TimerSearch) ([]*data.Timer, error)
	TimerStart(id string) (*data.Timer, error)
	TimerStop(id string) (*data.Timer, error)
	TimerUpdateComment(id, comment string) (*data.Timer, error)
	TimerArchive(id string, archive bool) (*data.Timer, error)
	TimerDelete(id string) error
	//these are implemented by the logic
	TimerSubmit(timerID string, finishTime *time.Time) (*data.Timer, error)
}

func New(parameters ...interface{}) interface {
	Logic
	Owner
} {
	l := &logic{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case meta.Timer:
			l.Timer = p
		case logger.Logger:
			l.Logger = p
		}
	}
	if l.Timer == nil {
		panic("no meta found for timer")
	}
	return l
}

func (l *logic) Start() (err error) {
	l.Lock()
	defer l.Unlock()

	if l.started {
		return
	}
	l.started = true

	return
}

func (l *logic) Stop() (err error) {
	l.Lock()
	defer l.Unlock()

	if !l.started {
		return
	}
	l.started = false

	return
}

func (l *logic) TimerSubmit(id string, submitTime *time.Time) (*data.Timer, error) {
	if submitTime == nil || submitTime.UnixNano() <= 0 {
		submitTime = new(time.Time)
		*submitTime = time.Now()
	}
	timer, err := l.Timer.TimerSubmit(id, submitTime.UnixNano())
	if err != nil {
		return nil, err
	}
	return timer, nil
}

func (l *logic) TimerUpdateComment(id, comment string) (*data.Timer, error) {
	timer, err := l.Timer.TimerUpdate(id, data.TimerPartial{
		Comment: &comment,
	})
	if err != nil {
		return nil, err
	}
	return timer, nil
}

func (l *logic) TimerArchive(id string, archived bool) (*data.Timer, error) {
	timer, err := l.Timer.TimerUpdate(id, data.TimerPartial{
		Archived: &archived,
	})
	if err != nil {
		return nil, err
	}
	return timer, nil
}
