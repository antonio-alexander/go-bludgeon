package logic

import (
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
)

// Logic defines functions that describe the business logic
// of the timers micro service
type Logic interface {
	meta.TimeSlice
	meta.Timer

	// IsConnected can be used to determine whether or not
	// the underlying change handler is connected
	IsConnected() bool
}
