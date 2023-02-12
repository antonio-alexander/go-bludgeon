package meta

import (
	"context"
	"errors"

	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal_errors "github.com/antonio-alexander/go-bludgeon/internal/errors"
)

// these constants are used to generate change specific errors
const (
	ChangeNotFound               string = "change not found"
	ChangeNotWritten             string = "change not written; data id not provided"
	ChangeConflictWrite          string = "cannot write change"
	ChangeNotDeletedConflict     string = "change not deleted, not fully acknowledged"
	RegistrationNotFound         string = "registration not found"
	RegistrationNotWritten       string = "registration not written; id not provided"
	RegistrationChangeNotWritten string = "registration change not written; change id not provided"
)

// these are error variables used within the change meta
var (
	ErrChangeNotFound               = internal_errors.NewNotFound(errors.New(ChangeNotFound))
	ErrChangeNotWritten             = internal_errors.NewNotUpdated(errors.New(ChangeNotWritten))
	ErrChangeConflictWrite          = internal_errors.NewConflict(errors.New(ChangeConflictWrite))
	ErrChangeNotDeletedConflict     = internal_errors.NewConflict(errors.New(ChangeNotDeletedConflict))
	ErrRegistrationNotFound         = internal_errors.NewNotFound(errors.New(RegistrationNotFound))
	ErrRegistrationNotWritten       = internal_errors.NewNotFound(errors.New(RegistrationNotWritten))
	ErrRegistrationChangeNotWritten = internal_errors.NewNotFound(errors.New(RegistrationChangeNotWritten))
)

// SerializedData provides a struct that describes the representation
// of the data when serialized
type SerializedData struct {
	Changes             map[string]data.Change         `json:"changes"`
	Registrations       map[string]struct{}            `json:"registrations"`
	RegistrationChanges map[string]map[string]struct{} `json:"registration_changes"`
}

// Serializer is an interface that can be used to convert the contents of
// meta into a scalar type
type Serializer interface {
	//Serialize can be used to convert all available metadata
	// to a single pointer to be used to serialize to bytes
	Serialize() (*SerializedData, error)

	//Deserialize can be used to provide metadata as a single pointer
	// once it's been deserialized from bytes
	Deserialize(data *SerializedData) error
}

// Change is an interface that groups functions to interact with one or more
// changes
type Change interface {
	ChangeCreate(ctx context.Context, change data.ChangePartial) (*data.Change, error)
	ChangeRead(ctx context.Context, changeId string) (*data.Change, error)
	ChangesDelete(ctx context.Context, changeIds ...string) error
	ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error)
}

type Registration interface {
	RegistrationUpsert(ctx context.Context, registrationId string) error
	RegistrationDelete(ctx context.Context, registrationId string) error
}

type RegistrationChange interface {
	RegistrationChangeUpsert(ctx context.Context, changeId string) error
	RegistrationChangesRead(ctx context.Context, registrationId string) (changeIds []string, err error)
	RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) (changeIdsToPrune []string, err error)
}
