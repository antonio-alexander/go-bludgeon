package data

import (
	"net/http"
)

const (
	PathChangeId       string = "change_id"
	PathRegistrationId string = "registration_id"
)

const (
	ParameterChangeIds       string = "change_ids"
	ParameterDataIds         string = "data_ids"
	ParameterTypes           string = "types"
	ParameterActions         string = "actions"
	ParameterServiceNames    string = "service_names"
	ParameterLatestVersion   string = "latest_version"
	ParameterSince           string = "since"
	ParameterRegistrationIds string = "registration_ids"
)

const (
	MethodChangeUpsert                  = http.MethodPatch
	MethodChangeRead                    = http.MethodGet
	MethodChangeDelete                  = http.MethodDelete
	MethodChangeRegister                = http.MethodPut
	MethodRegistrationChangeAcknowledge = http.MethodPut
	MethodRegistrationUpsert            = http.MethodPatch
	MethodRegistrationRead              = http.MethodGet
	MethodRegistrationDelete            = http.MethodDelete
)

const (
	RouteChanges                                  string = "/api/v1/changes"
	RouteChangesWebsocket                         string = RouteChanges + "/ws"
	RouteChangesSearch                            string = RouteChanges + "/search"
	RouteChangesParam                             string = RouteChanges + "/{" + PathChangeId + "}"
	RouteChangesParamf                            string = RouteChanges + "/%s"
	RouteChangesRegistration                      string = RouteChanges + "/registrations"
	RouteChangesRegistrationSearch                string = RouteChanges + "/registrations/search"
	RouteChangesRegistrationParam                 string = RouteChangesRegistration + "/{" + PathRegistrationId + "}"
	RouteChangesRegistrationParamf                string = RouteChangesRegistration + "/%s"
	RouteChangesRegistrationParamChanges          string = RouteChangesRegistrationParam + "/changes"
	RouteChangesRegistrationParamChangesf         string = RouteChangesRegistrationParamf + "/changes"
	RouteChangesRegistrationServiceIdAcknowledge  string = RouteChangesRegistration + "/{" + PathRegistrationId + "}/acknowledge"
	RouteChangesRegistrationServiceIdAcknowledgef string = RouteChangesRegistration + "/%s/acknowledge"
)
