package logic

import (
	"context"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	errors "github.com/pkg/errors"
)

const (
	logAlias                          string = "[logic] "
	PanicChangeMetaNotSet             string = "change meta not set"
	PanicRegistrationMetaNotSet       string = "change meta not set"
	PanicRegistrationChangeMetaNotSet string = "change meta not set"
	ChangeIdNotProvided               string = "change id not provided"
	RegisterFilterHandlerNotProvided  string = "unable to register; neither fitler or handler not provided"
	DefaultQueueSize                  int    = 100
)

var (
	ErrChangeIdNotProvided              = errors.New(ChangeIdNotProvided)
	ErrRegisterFilterHandlerNotProvided = errors.New(RegisterFilterHandlerNotProvided)
	QueueSize                           = DefaultQueueSize
)

type HandlerFx func(ctx context.Context, handlerId string, changes []*data.Change) error

type Logic interface {
	//changes
	ChangeUpsert(ctx context.Context, change data.ChangePartial) (*data.Change, error)
	ChangeRead(ctx context.Context, changeId string) (*data.Change, error)
	ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error)
	ChangesDelete(ctx context.Context, changeIds ...string) error

	//registrations
	RegistrationUpsert(ctx context.Context, registrationId string) error
	RegistrationChangesRead(ctx context.Context, registrationId string) ([]*data.Change, error)
	RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) error
	RegistrationDelete(ctx context.Context, registrationId string) error

	//handlers
	HandlerCreate(ctx context.Context, handleFx HandlerFx) (handlerId string, err error)
	HandlerDelete(ctx context.Context, handlerId string) error
}
