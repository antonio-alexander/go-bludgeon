package metamysql

//error constants
const (
	ErrTimerNotFoundf     string = "timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "timeSlice with id, \"%s\", not found locally"
)

//query constants
const (
	tableEmployees    string = "employees"
	tableTimers       string = "timers"
	tableTimeSlices   string = "time_slices"
	tableTimersV1     string = "timers_v1"
	tableTimeSlicesV1 string = "time_slices_v1"
)
