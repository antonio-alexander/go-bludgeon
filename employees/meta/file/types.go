package file

import (
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
)

//Configuration provides parameters to control the functionality of the
// meta pointer at run time, although it exposes the configuration of the
// underlying pointer, it prevents having to directly reference it
type Configuration struct {
	internal_file.Configuration
}

//File combines all the methods implemented by the underlying pointer
type File interface {
	meta.Owner
	meta.Employee

	//Initialize will configure and prepare the underlying pointer to
	// execute its business logic
	Initialize(config *Configuration) (err error)
}
