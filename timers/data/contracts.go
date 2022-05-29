package data

const (
	RouteBase            string = "/api/v1"
	RouteTimers          string = RouteBase + "/timers"
	RouteTimersSearch    string = RouteTimers + "/search"
	RouteTimersID        string = RouteTimers + "/{id}"
	RouteTimersIDStart   string = RouteTimersID + "/start"
	RouteTimersIDStop    string = RouteTimersID + "/stop"
	RouteTimersIDSubmit  string = RouteTimersID + "/submit"
	RouteTimersIDf       string = RouteTimers + "/%s"
	RouteTimersIDStartf  string = RouteTimersIDf + "/start"
	RouteTimersIDStopf   string = RouteTimersIDf + "/stop"
	RouteTimersIDSubmitf string = RouteTimersIDf + "/submit"
	RouteTimeSlices      string = RouteBase + "/time_slices"
	RouteTimeSlicesID    string = RouteTimeSlices + "/{id}"
	RouteTimeSlicesIDf   string = RouteTimeSlices + "/%s"
)

const PathID string = "id"

const (
	ParameterIDs         string = "ids"
	ParameterEmployeeID  string = "employee_id"
	ParameterEmployeeIDs string = "employee_ids"
	ParameterCompleted   string = "completed"
	ParameterArchived    string = "archived"
	ParameterTimerID     string = "timer_id"
	ParameterTimerIDs    string = "timer_ids"
)

type Contract struct {
	ID     string `json:"id,omitempty"`
	Finish int64  `json:"finish_time,omitempty"`
}
