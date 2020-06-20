package data

const (
	RouteTimer         = "/api/timer"
	RouteTimerCreate   = RouteTimer + "/create"
	RouteTimerRead     = RouteTimer + "/read"
	RouteTimerUpdate   = RouteTimer + "/update"
	RouteTimerDelete   = RouteTimer + "/delete"
	RouteTimersRead    = RouteTimer + "/read"
	RouteTimerStart    = RouteTimer + "/start"
	RouteTimerPause    = RouteTimer + "/pause"
	RouteTimerSubmit   = RouteTimer + "/submit"
	RouteTimeSlice     = "/api/timeslice"
	RouteTimeSliceRead = RouteTimeSlice + "/read"
)

type Contract struct {
	ID         string    `json:"id,omitempty"`
	StartTime  int64     `json:"start_time,string,omitempty"`
	PauseTime  int64     `json:"pause_time,string,omitempty"`
	FinishTime int64     `json:"finish_time,string,omitempty"`
	Timer      Timer     `json:"timer,omitempty"`
	TimeSlice  TimeSlice `json:"time_slice,omitempty"`
}
