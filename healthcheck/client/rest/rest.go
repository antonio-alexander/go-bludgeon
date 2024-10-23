package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	client "github.com/antonio-alexander/go-bludgeon/healthcheck/client"
	data "github.com/antonio-alexander/go-bludgeon/healthcheck/data"

	common "github.com/antonio-alexander/go-bludgeon/common"
	pkg_config "github.com/antonio-alexander/go-bludgeon/pkg/config"
	pkg_errors "github.com/antonio-alexander/go-bludgeon/pkg/errors"
	pkg_logger "github.com/antonio-alexander/go-bludgeon/pkg/logger"
	pkg_rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/client"
)

type restClient struct {
	pkg_logger.Logger
	client interface {
		common.Configurer
		common.Parameterizer
		pkg_rest.Client
	}
	config *Configuration
}

// New can be used to create a concrete instance of the rest client
// that implements the interfaces of logic.Logic and Owner
func New() interface {
	common.Configurer
	common.Parameterizer
	common.Initializer
	client.Client
} {
	return &restClient{
		Logger: pkg_logger.NewNullLogger(),
		client: pkg_rest.New(),
	}
}

func (r *restClient) doRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error) {
	bytes, statusCode, err := r.client.DoRequest(ctx, uri, method, data)
	if err != nil {
		return nil, err
	}
	switch statusCode {
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusNotModified, http.StatusConflict:
		return nil, pkg_errors.New(bytes)
	default:
		return bytes, nil
	}
}

func (r *restClient) SetParameters(parameters ...interface{}) {
	r.client.SetParameters(parameters...)
}

func (r *restClient) SetUtilities(parameters ...interface{}) {
	r.client.SetUtilities(parameters...)
	for _, p := range parameters {
		switch p := p.(type) {
		case pkg_logger.Logger:
			r.Logger = p
		}
	}
}

func (r *restClient) Configure(items ...interface{}) error {
	var configuration *Configuration
	var envs map[string]string

	for _, item := range items {
		switch v := item.(type) {
		case pkg_config.Envs:
			envs = v
		case *Configuration:
			configuration = v
		}
	}
	if configuration == nil {
		configuration = NewConfiguration()
		configuration.Default()
		configuration.FromEnv(envs)
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	r.config = configuration
	if err := r.client.Configure(r.config.Configuration); err != nil {
		return err
	}
	return nil
}

func (r *restClient) Initialize() error { return nil }

func (r *restClient) Shutdown() {}

func (r *restClient) HealthCheck(ctx context.Context) (*data.HealthCheck, error) {
	uri := fmt.Sprintf("http://%s:%s"+data.RouteHealthCheck, r.config.Address, r.config.Port)
	bytes, err := r.doRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	healthCheck := &data.HealthCheck{}
	if err = json.Unmarshal(bytes, healthCheck); err != nil {
		return nil, err
	}
	return healthCheck, nil
}
