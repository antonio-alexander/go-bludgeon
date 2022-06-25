package data

//route constants
const (
	RouteBase             string = "/api/v1"
	RouteTimers           string = RouteBase + "/timers"
	RouteTimersSearch     string = RouteTimers + "/search"
	RouteTimersID         string = RouteTimers + "/{id}"
	RouteTimersIDStart    string = RouteTimersID + "/start"
	RouteTimersIDStop     string = RouteTimersID + "/stop"
	RouteTimersIDSubmit   string = RouteTimersID + "/submit"
	RouteTimersIDComment  string = RouteTimersID + "/comment"
	RouteTimersIDArchive  string = RouteTimersID + "/archive"
	RouteTimersIDf        string = RouteTimers + "/%s"
	RouteTimersIDStartf   string = RouteTimersIDf + "/start"
	RouteTimersIDStopf    string = RouteTimersIDf + "/stop"
	RouteTimersIDSubmitf  string = RouteTimersIDf + "/submit"
	RouteTimersIDCommentf string = RouteTimersIDf + "/comment"
	RouteTimersIDArchivef string = RouteTimersIDf + "/archive"
	RouteTimeSlices       string = RouteBase + "/time_slices"
	RouteTimeSlicesSearch string = RouteTimeSlices + "/search"
	RouteTimeSlicesID     string = RouteTimeSlices + "/{id}"
	RouteTimeSlicesIDf    string = RouteTimeSlices + "/%s"
)

//path constants
const PathID string = "id"

//parameter constants
const (
	ParameterIDs         string = "ids"
	ParameterEmployeeID  string = "employee_id"
	ParameterEmployeeIDs string = "employee_ids"
	ParameterCompleted   string = "completed"
	ParameterArchived    string = "archived"
	ParameterTimerID     string = "timer_id"
	ParameterTimerIDs    string = "timer_ids"
)

//Contract is used for requests that don't have a
// solid data type to communicate data in the body
// of a request
type Contract struct {
	//Finish provides the finish time for a
	// timer (or time slice)
	// example: 1653719229
	Finish int64 `json:"finish_time,omitempty"`
}
