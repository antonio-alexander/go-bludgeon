package swagger

import (
	"github.com/antonio-alexander/go-bludgeon/healthcheck/data"
)

// swagger:route GET /healthcheck healthcheck
// Executes a healthcheck function.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
// responses:
//   200: HealthCheckResponseOk

// swagger:response HealthCheckResponseOk
type HealthCheckResponseOk struct {
	// in:body
	Body data.HealthCheck
}
