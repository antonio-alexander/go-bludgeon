package bludgeonclient

// "text/tabwriter"

//error constants
const (
	ErrStarted            string = "cache started"
	ErrNotStarted         string = "cache not started"
	ErrTimerNotFoundf     string = "Timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "TimeSlice with id, \"%s\", not found locally"
	ErrTimerIsArchivedf   string = "Timer with id, \"%s\", is archived"
	ErrNoActiveTimeSlicef string = "Timer with id, \"%s\", has no active slice"
)

//SerializedData
type SerializedData struct {
	LookupTimers     map[string]string `json:"LookupTimers,omit_empty"`
	LookupTimeSlices map[string]string `json:"LookupTimeSlices,omit_empty"`
}

type Configuration struct {
	ServerAddress string
	ServerPort    string
	ClientAddress string
	ClientPort    string
	Task          int64
	Employee      int64
}

//common constants
const (
	SQL_DRIVER = "sqlite3"
	HELP       = "Help Goes Here!"
	longForm   = "Jan 2, 2006 at 3:04pm (MST)"
	shortForm  = "2006-Jan-02"
	// tabwriterFlag = tabwriter.Debug //tabwriter.AlignRight | tabwriter.Debug
)
