package rest

import (
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	client "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
)

type Configuration struct {
	client.Configuration
}

//Client is an interface that provides all of the functions
// that should only be accessible by the entity that instantiates
// the rest client
type Client interface {
	logic.Logic

	//Initialize can be used to ready the underlying pointer for use
	Initialize(*Configuration) error
}
