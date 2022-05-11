package memory_test

import (
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	metamemory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"

	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"
)

func TestMetaMemory(t *testing.T) {
	m := metamemory.New(
		logger.New(),
	)
	t.Run("Timer CRUD", tests.TestTimerCRUD(m))
	t.Run("Timers Read", tests.TestTimersRead(m))
	t.Run("Timer Logic", tests.TestTimerLogic(m))
	m.Shutdown()
}
