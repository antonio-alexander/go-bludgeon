package memory_test

import (
	"context"
	"testing"

	metamemory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	tests "github.com/antonio-alexander/go-bludgeon/timers/meta/tests"

	"github.com/stretchr/testify/assert"
)

func TestMetaMemory(t *testing.T) {
	ctx, logger := context.TODO(), internal_logger.New()
	err := logger.Configure(internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "test_meta_memory",
	})
	assert.Nil(t, err)
	m := metamemory.New()
	m.SetUtilities(logger)
	defer m.Shutdown()

	t.Run("Timer CRUD", tests.TestTimerCRUD(ctx, m))
	t.Run("Timers Read", tests.TestTimersRead(ctx, m))
	t.Run("Timer Logic", tests.TestTimerLogic(ctx, m))
}
