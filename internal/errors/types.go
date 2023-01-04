package errors

import "encoding/json"

//wanna be able to create a specific kind of error

// Error represents any error that can be generated by bludgeon
// swagger:model Error
type Error struct {
	err error

	// Error contains the text of any failed operation
	// example: employee not found
	ErrorMessage string `json:"error_message"`

	// Type contains the text of the type of error
	// example: not found
	ErrorType string `json:"error_type"`
}

func New(items ...interface{}) (e Error) {
	for _, item := range items {
		switch i := item.(type) {
		case []byte:
			_ = json.Unmarshal(i, &e)
		case string:
			e.ErrorMessage = i
		}
	}
	return
}

func (e Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.ErrorMessage
}

func (e Error) Type() string {
	return e.ErrorType
}

func (e Error) Is(err error) bool {
	if err, ok := err.(interface {
		Type() string
	}); ok {
		return e.ErrorType == err.Type()
	}
	return false
}

const (
	ErrorTypeNotFound string = "ERR_NOT_FOUND"
	ErrTypeNotCreated string = "ERR_NOT_CREATED"
	ErrTypeNotUpdated string = "ERR_NOT_UPDATED"
	ErrTypeConflict   string = "ERR_CONFLICT"
)

func NewNotFound(err error) error {
	return &Error{
		err:          err,
		ErrorMessage: err.Error(),
		ErrorType:    ErrorTypeNotFound,
	}
}

func NewNotCreated(err error) error {
	return &Error{
		err:          err,
		ErrorMessage: err.Error(),
		ErrorType:    ErrTypeNotCreated,
	}
}

func NewNotUpdated(err error) error {
	return &Error{
		err:          err,
		ErrorMessage: err.Error(),
		ErrorType:    ErrTypeNotUpdated,
	}
}

func NewConflict(err error) error {
	return &Error{
		err:          err,
		ErrorMessage: err.Error(),
		ErrorType:    ErrTypeConflict,
	}
}
