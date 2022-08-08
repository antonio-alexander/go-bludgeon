package client

import (
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	internal_client "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"
)

type Configuration struct {
	internal_client.Configuration
}

//Client is an interface that provides all of the functions
// that should only be accessible by the entity that instantiates
// the rest client
type Client interface {
	logic.Logic

	//Initialize can be used to ready the underlying pointer for use
	Initialize(*Configuration) error
}
