package mysql

import (
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"

	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
)

//MySQL combines all the methods implemented by the underlying pointer
type MySQL interface {
	internal_meta.Owner
	meta.Employee

	//Initialize will configure and prepare the underlying pointer to
	// execute its business logic
	Initialize(config *internal_mysql.Configuration) (err error)
}
