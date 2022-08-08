package memory_test

import (
	"context"
	"testing"

	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	metamemory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"

	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"
)

func TestMetaMemory(t *testing.T) {
	ctx := context.TODO()
	m := metamemory.New(
		logger.New(),
	)
	t.Run("Timer CRUD", tests.TestTimerCRUD(ctx, m))
	t.Run("Timers Read", tests.TestTimersRead(ctx, m))
	t.Run("Timer Logic", tests.TestTimerLogic(ctx, m))
	m.Shutdown()
}
