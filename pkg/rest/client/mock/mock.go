package mock

import (
	"context"
	"sync"

	rest "github.com/antonio-alexander/go-bludgeon/pkg/rest/client"
)

type mock struct {
	sync.RWMutex
	doRequest struct {
		bytes      []byte
		statusCode int
		err        error
	}
}

func New() interface {
	rest.Client
	Mock
} {
	return &mock{}
}

func (m *mock) MockDoRequest(bytes []byte, statusCode int, err error) {
	m.Lock()
	defer m.Unlock()

	m.doRequest.bytes = bytes
	m.doRequest.statusCode = statusCode
	m.doRequest.err = err
}

func (m *mock) DoRequest(ctx context.Context, uri, method string, data []byte) ([]byte, int, error) {
	m.RLock()
	defer m.RUnlock()

	return m.doRequest.bytes, m.doRequest.statusCode, m.doRequest.err
}
