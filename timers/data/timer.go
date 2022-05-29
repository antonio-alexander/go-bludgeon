package data

// swagger:model Timer
//Timer is a high-level object that describes a single unit of time for a given "task". A timer may
// be started and paused many times, but only submitted once. Although there's an obvious desire to
// be able to edit a timer after its submitted, but before it's billed/invoiced. Prior to submission
// elapsed time should "always" be the sum of all the associated active timers and won't necessarily
// be the difference between finish and start (but may be if edited post submission).
type Timer struct {
	//Whether or not a timer has been completed
	// example: false
	Completed bool `json:"completed"`

	//Whether or not a timer has been archived
	// example: false
	Archived bool `json:"archived"`

	//The start time for the timer
	// example: 1653719208
	Start int64 `json:"start"`

	//The finish timer for the timer
	// example: 1653719229
	Finish int64 `json:"finish"`

	//The elapsed time for the timer, the elapsed time isn't
	// necessarily the difference between the finish and start
	// time
	// example: 21
	ElapsedTime int64 `json:"elasped_time"`

	//The ID of an employee (v4 UUID)
	// example: "2e3a4156-b415-4120-982f-399182e99588"
	EmployeeID string `json:"employee_id"`

	//The ID of the active time slice (v4 UUID)
	// example: "a33f813e-e9bc-46ad-9956-0c4b6c1367ab"
	ActiveTimeSliceID string `json:"active_time_slice_id"`

	//The ID of an employee (v4 UUID)
	// example: "24dfe1eb-26a7-41db-a647-fe6cc5e77ab8"
	ID string `json:"id"`

	//A comment describing the timer
	// example: "This is a timer for lunch"
	Comment string `json:"comment"`

	//Used for accounting of this unique timer
	Audit
}

// swagger:model TimerPartial
//TimerPartial represents the properties in timer that can be
// modified from the outside
type TimerPartial struct {
	//Whether or not a timer has been completed
	// example: true
	Completed *bool `json:"completed,omitempty"`

	//Whether or not a timer has been archived
	// example: true
	Archived *bool `json:"archived,omitempty"`

	//The ID of an employee (v4 UUID)
	// example: "24b32c23-e3a0-44d1-bdd4-9c370c050b29"
	EmployeeID *string `json:"employee_id,omitempty"`

	//A comment describing the timer
	// example: "This is a timer for breakfast"
	Comment *string `json:"comment,omitempty"`

	//The finish timer for the timer
	// example: 1653719229
	Finish *int64 `json:"finish,omitempty"`
}
