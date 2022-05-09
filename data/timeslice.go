package data

//TimeSlice is the basic unit of "time", the idea is that a task may span over multiple slices that
// aren't necessarily contiguous, but can be added together to get an altogether time.  This should
// reduce the overall error when you pause and restart timers (not time slices) from multiple locations
// time slices can be deleted/archived, but not "edited"
type TimeSlice struct {
	Completed   bool   `json:"completed"`
	Start       int64  `json:"start"`
	Finish      int64  `json:"finish"`
	ElapsedTime int64  `json:"elapsed_time"`
	ID          string `json:"id"`
	TimerID     string `json:"timer_id"`
	Audit
}

type TimeSlicePartial struct {
	TimerID   *string
	Completed *bool
	Start     *int64
	Finish    *int64
}

type TimeSliceSearch struct {
	Completed *bool
	TimerID   *string
	TimerIDs  []string
	IDs       []string
}

type TimeSliceByStart []*TimeSlice

func (t TimeSliceByStart) Len() int {
	return len(t)
}
func (t TimeSliceByStart) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TimeSliceByStart) Less(i, j int) bool {
	return t[i].Start < t[j].Start
}
