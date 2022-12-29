package grpc

import (
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
)

// Client is an interface that provides all of the functions
// that should only be accessible by the entity that instantiates
// the rest client
type Client interface {
	logic.Logic
}
