package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/gorilla/websocket"
)

type client struct {
	sync.RWMutex
	*websocket.Conn
	websocket.Dialer
	logger.Logger
	connected  bool
	configured bool
	config     *Configuration
}

func New() interface {
	Client
	internal.Configurer
	internal.Parameterizer
	internal.Closer
} {
	return &client{
		Logger: logger.NewNullLogger(),
		Dialer: *websocket.DefaultDialer,
	}
}

func (c *client) pingHandler(ping string) error {
	if !c.connected {
		return ErrNotConnected
	}
	c.Trace(logAlias+"ping received: %s", ping)
	return nil
}

func (c *client) pongHandler(pong string) error {
	if !c.connected {
		return ErrNotConnected
	}
	c.Trace(logAlias+"pong received: %s", pong)
	return nil
}

func (c *client) SetParameters(parameters ...interface{}) {
	//use this to set common utilities/parameters
}

func (c *client) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			c.Logger = p
		}
	}
}

func (c *client) Configure(items ...interface{}) error {
	c.Lock()
	defer c.Unlock()

	var envs map[string]string
	var configuration *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			envs = v
		case *Configuration:
			configuration = v
		}
	}
	if configuration == nil {
		configuration = new(Configuration)
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	c.config = configuration
	c.configured = true
	return nil
}

func (c *client) Connect(ctx context.Context, url string, requestHeader http.Header) (*http.Response, error) {
	c.Lock()
	defer c.Unlock()
	conn, response, err := c.DialContext(ctx, url, requestHeader)
	if err != nil {
		return response, err
	}
	c.Conn = conn
	c.SetPingHandler(c.pingHandler)
	c.SetPongHandler(c.pongHandler)
	c.connected = true
	return response, nil
}

func (c *client) IsConnected() bool {
	c.RLock()
	defer c.RUnlock()
	return c.connected
}

func (c *client) Write(item interface{}) error {
	c.Lock()
	defer c.Unlock()
	if !c.connected {
		return ErrNotConnected
	}
	if c.config.WriteTimeout > 0 {
		deadline := time.Now().Add(c.config.WriteTimeout)
		if err := c.SetWriteDeadline(deadline); err != nil {
			c.connected = false
			return err
		}
	}
	if err := c.Conn.WriteJSON(item); err != nil {
		c.connected = false
		return err
	}
	return nil
}

func (c *client) Read(item interface{}) error {
	//KIM: you can safely read concurrently
	if !c.connected {
		return ErrNotConnected
	}
	if c.config.ReadTimeout > 0 {
		deadline := time.Now().Add(c.config.ReadTimeout)
		if err := c.SetReadDeadline(deadline); err != nil {
			c.connected = false
			return err
		}
	}
	if err := c.Conn.ReadJSON(item); err != nil {
		c.connected = false
		return err
	}
	return nil
}

func (c *client) Close() {
	if !c.connected {
		return
	}
	if err := c.Conn.Close(); err != nil {
		c.Error("error while closing connection: %s", err)
	}
	c.connected = false
}
