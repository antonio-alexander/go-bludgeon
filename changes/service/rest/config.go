package rest

import (
	internal_websocket "github.com/antonio-alexander/go-bludgeon/pkg/websocket/server"
)

type Configuration struct {
	internal_websocket.Configuration
}

func (c *Configuration) Default() {
	c.ReadTimeout = internal_websocket.DefaultReadTimeout
	c.PingInterval = internal_websocket.DefaultPingInterval
	c.PingTimeout = internal_websocket.DefaultPingTimeout
	c.WriteTimeout = internal_websocket.DefaultWriteTimeout
}

func (c *Configuration) FromEnv(envs map[string]string) {
	//
}

func (c *Configuration) Validate() error {
	//
	return nil
}
