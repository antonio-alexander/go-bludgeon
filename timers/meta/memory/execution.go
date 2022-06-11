package memory

import (
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"

	"github.com/google/uuid"
)

func generateID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func elapsedTime(timer *data.Timer, timeSlices []*data.TimeSlice) *data.Timer {
	timer.ElapsedTime = int64(0)
	for _, timeSlice := range timeSlices {
		if timeSlice.Finish <= 0 {
			timer.ElapsedTime += time.Now().UnixNano() - timeSlice.Start
			continue
		}
		timer.ElapsedTime += timeSlice.Finish - timeSlice.Start
	}
	return timer
}

func copyTimer(t *data.Timer) *data.Timer {
	return &data.Timer{
		LastUpdated:       t.LastUpdated,
		LastUpdatedBy:     t.LastUpdatedBy,
		Version:           t.Version,
		Completed:         t.Completed,
		Archived:          t.Archived,
		Start:             t.Start,
		Finish:            t.Finish,
		ElapsedTime:       t.ElapsedTime,
		EmployeeID:        t.EmployeeID,
		ActiveTimeSliceID: t.ActiveTimeSliceID,
		ID:                t.ID,
		Comment:           t.Comment,
	}
}

func copyTimeSlice(t *data.TimeSlice) *data.TimeSlice {
	return &data.TimeSlice{
		Completed:     t.Completed,
		Start:         t.Start,
		Finish:        t.Finish,
		ElapsedTime:   t.ElapsedTime,
		ID:            t.ID,
		TimerID:       t.TimerID,
		LastUpdated:   t.LastUpdated,
		LastUpdatedBy: t.LastUpdatedBy,
		Version:       t.Version,
	}
}

func validateTimeSlice(t data.TimeSlice) error {
	if !t.Validate() {
		if t.TimerID == "" {
			return ErrTimeSliceNoTimerID
		}
		if t.Start <= 0 {
			return ErrTimeSliceStartZero
		}
		if t.Finish != 0 && t.Finish <= t.Start {
			return ErrTimeSliceFinishLessStart
		}
	}
	return nil
}
