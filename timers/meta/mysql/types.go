package mysql

import (
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"
)

//Configuration provides parameters to control the functionality of the
// meta pointer at run time, although it exposes the configuration of the
// underlying pointer, it prevents having to directly reference it
type Configuration struct {
	internal_mysql.Configuration
}

//MySQL combines all the methods implemented by the underlying pointer
type MySQL interface {
	meta.Timer
	meta.TimeSlice
	meta.Shutdown

	//Initialize will configure and prepare the underlying pointer to
	// execute its business logic
	Initialize(config *Configuration) (err error)
}
