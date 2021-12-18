package file

type Owner interface {
	Initialize(config *Configuration) (err error)
}
