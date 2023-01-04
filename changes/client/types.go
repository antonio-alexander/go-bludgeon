package client

import (
	"context"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
)

type Client interface {
	ChangeUpsert(ctx context.Context, changePartial data.ChangePartial) (*data.Change, error)
	ChangeRead(ctx context.Context, changeId string) (*data.Change, error)
	ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error)
	ChangeDelete(ctx context.Context, changeId string) error

	RegistrationUpsert(ctx context.Context, serviceName string) error
	RegistrationChangesRead(ctx context.Context, registrationId string) ([]*data.Change, error)
	RegistrationChangeAcknowledge(ctx context.Context, serviceName string, changeIds ...string) error
	RegistrationDelete(ctx context.Context, serviceName string) error
}

type Handler interface {
	HandlerCreate(handlerFx HandlerFx) (handlerId string, err error)
	HandlerConnected(handlerId string) (bool, error)
	HandlerDelete(handlerId string) (err error)
}

type HandlerFx func(...*data.Change) error
