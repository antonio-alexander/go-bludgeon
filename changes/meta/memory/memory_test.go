package memory_test

import (
	"testing"

	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/tests"
	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/stretchr/testify/assert"
)

type memoryMetaTest struct {
	metaMemory common.Initializer
	*tests.Fixture
}

func newMemoryMetaTest() *memoryMetaTest {
	logger := logger.New()
	metaMemory := memory.New()
	metaMemory.SetUtilities(logger)
	return &memoryMetaTest{
		metaMemory: metaMemory,
		Fixture:    tests.NewFixture(metaMemory),
	}
}

func (m *memoryMetaTest) Initialize(t *testing.T) {
	err := m.metaMemory.Initialize()
	if !assert.Nil(t, err) {
		assert.FailNow(t, "unable to initialize persistence")
	}
}

func (m *memoryMetaTest) Shutdown(t *testing.T) {
	m.metaMemory.Shutdown()
}

func testMetaMemory(t *testing.T) {
	m := newMemoryMetaTest()
	m.Initialize(t)
	defer m.Shutdown(t)

	t.Run("Change CRUD", m.TestChangeCRUD)
	t.Run("Changes Read", m.TestChangeSearch)
	t.Run("Registration CRUD", m.TestRegistrationCRUD)
	t.Run("Change Registrations", m.TestRegistrationChanges)
}

func TestMetaMemory(t *testing.T) {
	testMetaMemory(t)
}
