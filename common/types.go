package common

type Configurer interface {
	Configure(items ...interface{}) error
}

type Initializer interface {
	Initialize() error
	Shutdowner
}

type Shutdowner interface {
	Shutdown()
}

type Closer interface {
	Close()
}

type Parameterizer interface {
	SetParameters(parameters ...interface{})
	SetUtilities(parameters ...interface{})
}
