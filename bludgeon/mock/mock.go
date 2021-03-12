package bludgeonmock

// bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"

type mock struct{}

func NewMockFunctional() interface {
	// bludgeon.FunctionalManage
	// bludgeon.FunctionalOwner
	// bludgeon.FunctionalTimer
	// bludgeon.FunctionalTimeSlice
} {
	return mock{}
}
