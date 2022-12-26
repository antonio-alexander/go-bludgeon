package restclient

import (
	"context"
)

type Client interface {
	DoRequest(ctx context.Context, uri, method string, data []byte) ([]byte, error)
}
