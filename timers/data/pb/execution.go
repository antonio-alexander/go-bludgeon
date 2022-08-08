package pb

import "github.com/antonio-alexander/go-bludgeon/timers/data"

func FromTimerPartial(t *data.TimerPartial) *TimerPartial {
	if t == nil {
		return nil
	}
	TimerPartial := &TimerPartial{}
	if t.Completed != nil {
		TimerPartial.CompletedOneof = &TimerPartial_Completed{
			Completed: *t.Completed,
		}
	}
	if t.Archived != nil {
		TimerPartial.ArchivedOneof = &TimerPartial_Archived{
			Archived: *t.Archived,
		}
	}
	if t.EmployeeID != nil {
		TimerPartial.EmployeeIdOneof = &TimerPartial_EmployeeId{
			EmployeeId: *t.EmployeeID,
		}
	}
	if t.Comment != nil {
		TimerPartial.CommentOneof = &TimerPartial_Comment{
			Comment: *t.Comment,
		}
	}
	if t.Finish != nil {
		TimerPartial.FinishOneof = &TimerPartial_Finish{
			Finish: *t.Finish,
		}
	}
	return TimerPartial
}

func ToTimerPartial(t *TimerPartial) *data.TimerPartial {
	if t == nil {
		return nil
	}
	TimerPartial := &data.TimerPartial{}
	if t.CompletedOneof != nil {
		s := t.GetCompleted()
		TimerPartial.Completed = &s
	}
	if t.ArchivedOneof != nil {
		s := t.GetArchived()
		TimerPartial.Archived = &s
	}
	if t.EmployeeIdOneof != nil {
		s := t.GetEmployeeId()
		TimerPartial.EmployeeID = &s
	}
	if t.CommentOneof != nil {
		s := t.GetComment()
		TimerPartial.Comment = &s
	}
	if t.FinishOneof != nil {
		s := t.GetFinish()
		TimerPartial.Finish = &s
	}
	return TimerPartial
}

func FromTimer(t *data.Timer) *Timer {
	if t == nil {
		return nil
	}
	return &Timer{
		Completed:         t.Completed,
		Archived:          t.Archived,
		Start:             t.Start,
		Finish:            t.Finish,
		ElapsedTime:       t.ElapsedTime,
		EmployeeId:        t.EmployeeID,
		ActiveTimeSliceId: t.ActiveTimeSliceID,
		Id:                t.ID,
		Comment:           t.Comment,
		LastUpdated:       t.LastUpdated,
		LastUpdatedBy:     t.LastUpdatedBy,
		Version:           int32(t.Version),
	}
}

func ToTimer(t *Timer) *data.Timer {
	if t == nil {
		return nil
	}
	return &data.Timer{
		Completed:         t.GetCompleted(),
		Archived:          t.GetArchived(),
		Start:             t.GetStart(),
		Finish:            t.GetFinish(),
		ElapsedTime:       t.GetElapsedTime(),
		EmployeeID:        t.GetEmployeeId(),
		ActiveTimeSliceID: t.GetActiveTimeSliceId(),
		ID:                t.GetId(),
		Comment:           t.GetComment(),
		LastUpdated:       t.GetLastUpdated(),
		LastUpdatedBy:     t.GetLastUpdatedBy(),
		Version:           int(t.GetVersion()),
	}
}

func FromTimers(e []*data.Timer) []*Timer {
	var Timers []*Timer
	for _, e := range e {
		Timers = append(Timers, FromTimer(e))
	}
	return Timers
}

func ToTimers(e []*Timer) []*data.Timer {
	var Timers []*data.Timer
	for _, e := range e {
		Timers = append(Timers, ToTimer(e))
	}
	return Timers
}

func FromTimerSearch(t *TimerSearch) *data.TimerSearch {
	if t == nil {
		return nil
	}
	TimerSearch := &data.TimerSearch{
		IDs:         t.GetIds(),
		EmployeeIDs: t.GetEmployeeIds(),
	}
	if t.EmployeeIdOneof != nil {
		s := t.GetEmployeeId()
		TimerSearch.EmployeeID = &s
	}
	if t.CompletedOneof != nil {
		s := t.GetCompleted()
		TimerSearch.Completed = &s
	}
	if t.ArchivedOneof != nil {
		s := t.GetArchived()
		TimerSearch.Archived = &s
	}
	return TimerSearch
}

func ToTimerSearch(t *data.TimerSearch) *TimerSearch {
	if t == nil {
		return nil
	}
	TimerSearch := &TimerSearch{
		Ids:         t.IDs,
		EmployeeIds: t.EmployeeIDs,
	}
	if t.EmployeeID != nil {
		TimerSearch.EmployeeIdOneof = &TimerSearch_EmployeeId{
			EmployeeId: *t.EmployeeID,
		}
	}
	if t.Completed != nil {
		TimerSearch.CompletedOneof = &TimerSearch_Completed{
			Completed: *t.Completed,
		}
	}
	if t.Archived != nil {
		TimerSearch.ArchivedOneof = &TimerSearch_Archived{
			Archived: *t.Archived,
		}
	}
	return TimerSearch
}

func FromTimeSlicePartial(t *data.TimeSlicePartial) *TimeSlicePartial {
	if t == nil {
		return nil
	}
	TimeSlicePartial := &TimeSlicePartial{}
	if t.TimerID != nil {
		TimeSlicePartial.TimerIdOneof = &TimeSlicePartial_TimerId{
			TimerId: *t.TimerID,
		}
	}
	if t.Completed != nil {
		TimeSlicePartial.CompletedOneof = &TimeSlicePartial_Completed{
			Completed: *t.Completed,
		}
	}
	if t.Start != nil {
		TimeSlicePartial.StartOneof = &TimeSlicePartial_Start{
			Start: *t.Start,
		}
	}
	if t.Finish != nil {
		TimeSlicePartial.FinishOneof = &TimeSlicePartial_Finish{
			Finish: *t.Finish,
		}
	}
	return TimeSlicePartial
}

func ToTimeSlicePartial(t *TimeSlicePartial) *data.TimeSlicePartial {
	if t == nil {
		return nil
	}
	TimeSlicePartial := &data.TimeSlicePartial{}
	if t.CompletedOneof != nil {
		s := t.GetCompleted()
		TimeSlicePartial.Completed = &s
	}
	if t.CompletedOneof != nil {
		s := t.GetCompleted()
		TimeSlicePartial.Completed = &s
	}
	if t.TimerIdOneof != nil {
		s := t.GetTimerId()
		TimeSlicePartial.TimerID = &s
	}
	if t.StartOneof != nil {
		s := t.GetStart()
		TimeSlicePartial.Start = &s
	}
	if t.FinishOneof != nil {
		s := t.GetFinish()
		TimeSlicePartial.Finish = &s
	}
	return TimeSlicePartial
}

func FromTimeSlice(t *data.TimeSlice) *TimeSlice {
	if t == nil {
		return nil
	}
	return &TimeSlice{
		Completed:     t.Completed,
		Start:         t.Start,
		Finish:        t.Finish,
		ElapsedTime:   t.ElapsedTime,
		Id:            t.ID,
		TimerId:       t.TimerID,
		LastUpdated:   t.LastUpdated,
		LastUpdatedBy: t.LastUpdatedBy,
		Version:       int32(t.Version),
	}
}

func ToTimeSlice(t *TimeSlice) *data.TimeSlice {
	if t == nil {
		return nil
	}
	return &data.TimeSlice{
		Completed:     t.GetCompleted(),
		Start:         t.GetStart(),
		Finish:        t.GetFinish(),
		ElapsedTime:   t.GetElapsedTime(),
		ID:            t.GetId(),
		TimerID:       t.GetTimerId(),
		LastUpdated:   t.GetLastUpdated(),
		LastUpdatedBy: t.GetLastUpdatedBy(),
		Version:       int(t.GetVersion()),
	}
}

func FromTimeSlices(e []*data.TimeSlice) []*TimeSlice {
	var TimeSlices []*TimeSlice
	for _, e := range e {
		TimeSlices = append(TimeSlices, FromTimeSlice(e))
	}
	return TimeSlices
}

func ToTimeSlices(e []*TimeSlice) []*data.TimeSlice {
	var TimeSlices []*data.TimeSlice
	for _, e := range e {
		TimeSlices = append(TimeSlices, ToTimeSlice(e))
	}
	return TimeSlices
}

func FromTimeSliceSearch(t *TimeSliceSearch) *data.TimeSliceSearch {
	if t == nil {
		return nil
	}
	TimeSliceSearch := &data.TimeSliceSearch{
		IDs:      t.GetIds(),
		TimerIDs: t.GetTimerIds(),
	}
	if t.CompletedOneof != nil {
		s := t.GetCompleted()
		TimeSliceSearch.Completed = &s
	}
	if t.TimerIdOneof != nil {
		s := t.GetTimerId()
		TimeSliceSearch.TimerID = &s
	}
	return TimeSliceSearch
}

func ToTimeSliceSearch(t *data.TimeSliceSearch) *TimeSliceSearch {
	if t == nil {
		return nil
	}
	TimeSliceSearch := &TimeSliceSearch{
		Ids:      t.IDs,
		TimerIds: t.TimerIDs,
	}
	if t.TimerID != nil {
		TimeSliceSearch.TimerIdOneof = &TimeSliceSearch_TimerId{
			TimerId: *t.TimerID,
		}
	}
	if t.Completed != nil {
		TimeSliceSearch.CompletedOneof = &TimeSliceSearch_Completed{
			Completed: *t.Completed,
		}
	}
	return TimeSliceSearch
}
