package restclient

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
)

type client struct {
	*http.Client
	logger.Logger
	config *Configuration
}

func New() interface {
	internal.Configurer
	internal.Parameterizer
	Client
} {
	return &client{
		Logger: logger.NewNullLogger(),
		Client: new(http.Client),
	}
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
	if c == nil {
		configuration = new(Configuration)
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	c.config = configuration
	c.Client.Timeout = c.config.Timeout
	return nil
}

func (c *client) DoRequest(ctx context.Context, uri, method string, data []byte) ([]byte, int, error) {
	request, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, -1, err
	}
	response, err := c.Do(request)
	if err != nil {
		return nil, -1, err
	}
	data, err = io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, -1, err
	}
	return data, response.StatusCode, nil
}
