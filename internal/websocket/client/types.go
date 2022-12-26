package client

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

const (
	logAlias     string = "[websocket_client] "
	NotConnected string = "not connected"
)

var ErrNotConnected = errors.New(NotConnected)

type Client interface {
	Connect(ctx context.Context, url string, requestHeader http.Header) (*http.Response, error)
	Write(item interface{}) error
	Read(item interface{}) error
}
