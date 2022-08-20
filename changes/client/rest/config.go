package restclient

import (
	"errors"

	internal_rest_client "github.com/antonio-alexander/go-bludgeon/internal/rest/client"
	internal_websocket_client "github.com/antonio-alexander/go-bludgeon/internal/websocket/client"
)

const (
	EnvNameAddress         string = "BLUDGEON_CHANGES_REST_ADDRESS"
	EnvNamePort            string = "BLUDGEON_CHANGES_REST_PORT"
	EnvNameAutoAcknowledge string = "BLUDGEON_CHANGES_AUTO_ACKNOWLEDGE"
	DefaultAddress         string = "localhost"
	DefaultPort            string = "8014"
	DefaultAutoAcknowledge bool   = true
)

type Configuration struct {
	Rest      internal_rest_client.Configuration
	Websocket internal_websocket_client.Configuration
}

func (c *Configuration) Default() {
	c.Rest.Address = DefaultAddress
	c.Rest.Port = DefaultPort
	c.Rest.Timeout = internal_rest_client.DefaultTimeout
	c.Websocket.ReadTimeout = internal_websocket_client.DefaultReadTimeout
	c.Websocket.WriteTimeout = internal_websocket_client.DefaultWriteTimeout
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if address, ok := envs[EnvNameAddress]; ok {
		c.Rest.Address = address
	}
	if port, ok := envs[EnvNamePort]; ok {
		c.Rest.Port = port
	}
}

func (c *Configuration) Validate() error {
	if c.Rest.Address == "" {
		return errors.New(internal_rest_client.ErrAddressEmpty)
	}
	if c.Rest.Port == "" {
		return errors.New(internal_rest_client.ErrPortEmpty)
	}
	return nil
}
