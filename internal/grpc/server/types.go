package grpcserver

import "google.golang.org/grpc"

type RegisterFx func(server grpc.ServiceRegistrar)

type Registerer interface {
	Register(server grpc.ServiceRegistrar)
}
