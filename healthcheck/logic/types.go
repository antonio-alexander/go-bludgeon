package logic

import (
	"context"

	"github.com/antonio-alexander/go-bludgeon/healthcheck/data"
)

type Logic interface {
	HealthCheck(ctx context.Context) (*data.HealthCheck, error)
}
