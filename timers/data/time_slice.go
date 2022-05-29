package data

// swagger:model TimeSlice
//TimeSlice is the basic unit of "time", the idea is that a task may span over multiple slices that
// aren't necessarily contiguous, but can be added together to get an altogether time.  This should
// reduce the overall error when you pause and restart timers (not time slices) from multiple locations
// time slices can be deleted/archived, but not "edited"
type TimeSlice struct {
	//Whether or not a timer has been completed
	// example: true
	Completed bool `json:"completed"`

	//The start time of the time slice
	// example: 1653720177
	Start int64 `json:"start"`

	//The finish time of the time slice
	// exmample: 1653720184
	Finish int64 `json:"finish"`

	//The elapsed time of the time slice
	// example: 7
	ElapsedTime int64 `json:"elapsed_time"`

	//The id of the time slice
	// example: "ff7e87af-e6c5-44c3-851f-8801a33ad888"
	ID string `json:"id"`

	//The ID of the associated timer (v4 UUID)
	// example: "7f583116-c7b8-457d-97e0-be0670e9e27e"
	TimerID string `json:"timer_id"`

	//Used for accounting of this unique time slice
	Audit
}

// swagger:model TimeSlicePartial
//TimeSlicePartial can be used to update fields of a time slice
// that can be mutated (contrast with Audit)
type TimeSlicePartial struct {
	//The ID of the associated timer (v4 UUID), this cannot be
	// changed post creation
	// example: "7f583116-c7b8-457d-97e0-be0670e9e27e"
	TimerID *string

	//Whether or not a timer has been completed
	// example: true
	Completed *bool

	//The start time of the time slice
	// example: 1653720177
	Start *int64

	//The finish time of the time slice
	// exmample: 1653720184
	Finish *int64
}

//TimeSliceByStart can be used to sort slices by
// start time using the sort.Sort function
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
