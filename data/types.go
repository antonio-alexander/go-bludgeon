package data

//--------------------------------------------------------------------------------------------
// types.go will contain basic types used by the package in general these will be the type sued
// by other functions, it'll include both their types and their methods, the types created here
// can be used elsewhere, unlike other types
//--------------------------------------------------------------------------------------------

//common constants
const (
	FmtTimeLong  = "Jan 2, 2006 at 3:04pm (MST)"
	FmtTimeShort = "2006-Jan-02"
	// DefaultFolder            = ".bludgeon"
	// DefaultConfigurationFile = "config/bludgeon_config.json"
	// DefaultCacheFile         = "data/bludgeon_cache.json"
)

//error constants
const (
	ErrBadProjectID        string = "projectID is invalid or missing"
	ErrBadEmployeeID       string = "employeeID is invalid or missing"
	ErrBadTaskID           string = "taskID is invalid or missing"
	ErrBadClientID         string = "clientID is invalid or missing"
	ErrBadTimerID          string = "timerID is invalid or missing"
	ErrBadTimeSliceID      string = "timeSlice is invalid or missing"
	ErrBadUnitID           string = "unitID is invalid or missing"
	ErrBadEmployeeIDTaskID string = "employeeID and/or TaskID is invalid or mising"
	ErrBadClientIDUnitID   string = "clientID and/or UnitID is invalid or missing"
	ErrTimerIsArchivedf    string = "timer with id, \"%s\", is archived"
	ErrNoActiveTimeSlicef  string = "timer with id, \"%s\", has no active slice"
)
