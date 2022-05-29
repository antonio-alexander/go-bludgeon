package rest

//error constants specific to the rest service
const (
	ErrAddressEmpty string = "address is empty"
	ErrPortEmpty    string = "port is empty"
	ErrPortBadf     string = "port is a non-integer: %s"
	ErrTimeoutBadf  string = "timeout is lte to 0: %v"
)

//endpoint constants
const (
	ErrStarted    string = "already started"
	ErrNotStarted string = "not started"
)
