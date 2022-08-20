package mysql

import (
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"
)

const (
	// logAlias                 string = "[meta_mysql] "
	tableChanges             string = "changes"
	tableRegistrations       string = "registrations"
	tableRegistrationChanges string = "registration_changes"
	tableRegistrationsV1     string = "registrations_v1"
	tableChangesV1           string = "changes_v1"
)

type Owner interface {
	//Initialize will configure and prepare the underlying pointer to
	// execute its business logic
	Initialize(config *internal_mysql.Configuration) (err error)
}
