package memory_test

import (
	"testing"

	metamemory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	tests "github.com/antonio-alexander/go-bludgeon/employees/meta/tests"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
)

func TestMetaMemory(t *testing.T) {
	meta, logger := metamemory.New(), internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "employee_memory_test",
		Level:  internal_logger.Trace,
	})
	meta.SetUtilities(logger)
	defer meta.Shutdown()

	t.Run("Employee CRUD", tests.TestEmployeeCRUD(meta))
}
