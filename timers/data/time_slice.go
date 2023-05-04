package data

import "encoding/json"

// swagger:model TimeSlice
// TimeSlice is the basic unit of "time", the idea is that a task may span over multiple slices that
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

	//LastUpdated represents the last time (unix nano) something was mutated
	// example: 1652417242000
	LastUpdated int64 `json:"last_updated"`

	//LastUpdatedBy will identify the last someone who mutated something
	// example: bludgeon_employee_memory
	LastUpdatedBy string `json:"last_updated_by"`

	//Version is an integer that's atomically incremented each time something i smutated
	// example: 1
	Version int `json:"version"`
}

func (t *TimeSlice) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TimeSlice) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, t)
}

func (t *TimeSlice) Validate() bool {
	if t.TimerID == "" {
		return false
	}
	if t.Start <= 0 {
		return false
	}
	if t.Finish != 0 && t.Finish <= t.Start {
		return false
	}
	return true
}

func (t *TimeSlice) Contains(tC TimeSlice) bool {
	//both time slices must have a non zero start and if
	// the finish is non zero it must be greater or equal to start
	if t.Start == 0 || (t.Finish > 0 && t.Finish <= t.Start) ||
		tC.Start == 0 || (tC.Finish > 0 && tC.Finish <= tC.Start) {
		return false
	}
	switch {
	default:
		//start must be greater than the finish
		contains := tC.Start <= t.Finish
		return contains
	case t.Finish == 0:
		//finish and start must be less than start
		contains := tC.Start >= t.Start || tC.Finish >= t.Start
		return contains
	case tC.Finish == 0:
		//start must be greater than finish and start
		contains := tC.Start <= t.Finish || tC.Start <= t.Start
		return contains
	}
}

// swagger:model TimeSlicePartial
// TimeSlicePartial can be used to update fields of a time slice
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

// TimeSliceByStart implements sort.Interface
// var _ sort.Interface = TimeSliceByStart{}
type TimeSliceByStart []*TimeSlice

// Len is the number of elements in the collection.
func (t TimeSliceByStart) Len() int {
	return len(t)
}

// Swap swaps the elements with indexes i and j.
func (t TimeSliceByStart) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (t TimeSliceByStart) Less(i, j int) bool {
	return t[i].Start < t[j].Start
}
