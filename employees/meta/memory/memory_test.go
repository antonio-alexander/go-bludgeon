package memory_test

import (
	"testing"

	metamemory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

func TestMetaMemory(t *testing.T) {
	m := metamemory.New(
		logger.New(),
	)
	t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
	m.Shutdown()
}
