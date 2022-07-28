package file

import (
	timers "github.com/antonio-alexander/go-bludgeon/timers/meta"

	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

//File combines all the methods implemented by the underlying pointer
type File interface {
	timers.Timer
	timers.TimeSlice
	internal_meta.Owner

	//Initialize will configure and prepare the underlying pointer to
	// execute its business logic
	Initialize(config *internal_file.Configuration) (err error)
}
