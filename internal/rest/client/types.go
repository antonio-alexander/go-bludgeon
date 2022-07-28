package restclient

import (
	"context"
)

type Client interface {
	Initialize(config *Configuration) error
	DoRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error)
}
