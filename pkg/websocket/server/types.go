package server

import "github.com/pkg/errors"

const (
	logAlias     string = "[websocket_server] "
	NotConnected string = "not connected"
)

var ErrNotConnected = errors.New(NotConnected)

type Server interface {
	Write(item interface{}) error
	Read(item interface{}) error
}
