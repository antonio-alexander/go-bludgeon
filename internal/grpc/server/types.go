package grpcserver

type RegisterFx func()

type Owner interface {
	Initialize(*Configuration, ...RegisterFx) error
	Shutdown()
}
