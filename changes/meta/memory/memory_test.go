package memory_test

import (
	"testing"

	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/tests"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"
)

func TestMetaMemory(t *testing.T) {
	logger := logger.New()
	m := memory.New()
	m.SetParameters(logger)
	defer func() {
		m.Shutdown()
	}()
	t.Run("Change CRUD", tests.TestChangeCRUD(m))
	t.Run("Changes Read", tests.TestChangeSearch(m))
	t.Run("Registration CRUD", tests.TestRegistrationCRUD(m))
	t.Run("Change Registrations", tests.TestRegistrationChanges(m))
}
