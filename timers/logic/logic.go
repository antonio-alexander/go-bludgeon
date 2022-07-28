package logic

import (
	"context"
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
	meta.TimeSlice
}

//New will instantiate a logic pointer that
// implements the Logic interface
func New(parameters ...interface{}) Logic {
	l := &logic{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case interface {
			meta.Timer
			meta.TimeSlice
		}:
			l.Timer = p
			l.TimeSlice = p
		case meta.Timer:
			l.Timer = p
		case meta.TimeSlice:
			l.TimeSlice = p
		case logger.Logger:
			l.Logger = p
		}
	}
	if l.Timer == nil {
		panic("no meta found for timer")
	}
	return l
}

//TimerSubmit can be used to stop a timer and set completed to true
func (l *logic) TimerSubmit(ctx context.Context, id string, submitTime *time.Time) (*data.Timer, error) {
	if submitTime == nil || submitTime.UnixNano() <= 0 {
		submitTime = new(time.Time)
		*submitTime = time.Now()
	}
	timer, err := l.Timer.TimerSubmit(ctx, id, submitTime.UnixNano())
	if err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerUpdateCommnet will only update the comment for timer with
// the provided id
func (l *logic) TimerUpdateComment(ctx context.Context, id, comment string) (*data.Timer, error) {
	timer, err := l.Timer.TimerUpdate(ctx, id, data.TimerPartial{
		Comment: &comment,
	})
	if err != nil {
		return nil, err
	}
	return timer, nil
}

//TimerArchive will only update the archive for timer with
// the provided id
func (l *logic) TimerArchive(ctx context.Context, id string, archived bool) (*data.Timer, error) {
	timer, err := l.Timer.TimerUpdate(ctx, id, data.TimerPartial{
		Archived: &archived,
	})
	if err != nil {
		return nil, err
	}
	return timer, nil
}
