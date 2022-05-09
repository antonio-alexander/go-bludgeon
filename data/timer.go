package data

//Timer is a high-level object that describes a single unit of time for a given "task". A timer may
// be started and paused many times, but only submitted once. Although there's an obvious desire to
// be able to edit a timer after its submitted, but before it's billed/invoiced. Prior to submission
// elapsed time should "always" be the sum of all the associated active timers and won't necessarily
// be the difference between finish and start (but may be if edited post submission).
type Timer struct {
	Audit
	Completed         bool   `json:"completed"`
	Archived          bool   `json:"archived"`
	Start             int64  `json:"start"`
	Finish            int64  `json:"finish"`
	ElapsedTime       int64  `json:"elasped_time"`
	EmployeeID        string `json:"employee_id"`
	ActiveTimeSliceID string `json:"active_time_slice_id"`
	ID                string `json:"id"`
	Comment           string `json:"comment"`
}

type TimerPartial struct {
	Completed  *bool
	Archived   *bool
	EmployeeID *string
	Comment    *string
}

type TimerSearch struct {
	EmployeeID  *string
	EmployeeIDs []string
	Completed   *bool
	Archived    *bool
	IDs         []string
}
