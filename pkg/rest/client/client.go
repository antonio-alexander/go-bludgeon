package restclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/pkg/config"
	"github.com/antonio-alexander/go-bludgeon/pkg/logger"
)

type client struct {
	*http.Client
	logger.Logger
	config *Configuration
}

func New() interface {
	common.Configurer
	common.Parameterizer
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
	var cc *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case *Configuration:
			cc = v
		}
	}
	if c == nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	c.config = cc
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
