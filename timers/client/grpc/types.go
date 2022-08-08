package client

import (
	internal_grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
)

//Client is an interface that provides all of the functions
// that should only be accessible by the entity that instantiates
// the rest client
type Client interface {
	logic.Logic

	//Initialize can be used to ready the underlying pointer for use
	Initialize(*internal_grpc.Configuration) error
}
