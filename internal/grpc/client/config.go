package grpcclient

import (
	"google.golang.org/grpc"
)

const (
	ErrAddressEmpty string = "address is empty"
	ErrPortEmpty    string = "port is empty"
	ErrPortBadf     string = "port is a non-integer: %s"
)

type Configuration struct {
	Address string            `json:"address"`
	Port    string            `json:"port"`
	Options []grpc.DialOption `json:"-"`
}
