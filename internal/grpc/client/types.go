package grpcclient

type Client interface {
	Initialize(*Configuration) error
	Shutdown()
}
