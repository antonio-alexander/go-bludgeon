package restclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"
)

type client struct {
	*http.Client
	logger.Logger
	config *Configuration
}

func New(parameters ...interface{}) interface {
	Client
} {
	var config *Configuration
	c := &client{
		Client: new(http.Client),
	}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			config = p
		case logger.Logger:
			c.Logger = p
		}
	}
	if config != nil {
		if err := c.Initialize(config); err != nil {
			panic(err)
		}
	}
	return c
}

func (c *client) Initialize(config *Configuration) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if err := config.Validate(); err != nil {
		return err
	}
	c.config = config
	c.Client.Timeout = c.config.Timeout
	return nil
}

func (c *client) DoRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	default:
		if len(data) > 0 {
			return nil, errors.New(string(data))
		}
		return nil, fmt.Errorf("failure: %d", response.StatusCode)
	case http.StatusOK:
		return data, nil
	case http.StatusNoContent:
		return nil, nil
	}
}
