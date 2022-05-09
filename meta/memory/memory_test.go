package memory_test

import (
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	metamemory "github.com/antonio-alexander/go-bludgeon/meta/memory"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"
)

func TestMetaMemory(t *testing.T) {
	m := metamemory.New(
		logger.New(),
	)
	t.Run("Employee CRUD", tests.TestEmployeeCRUD(m))
	t.Run("Timer CRUD", tests.TestTimerCRUD(m))
	t.Run("Timers Read", tests.TestTimersRead(m))
	t.Run("Timer Logic", tests.TestTimerLogic(m))
	m.Shutdown()
}
