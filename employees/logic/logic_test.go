package logic_test

import (
	"testing"

	"github.com/antonio-alexander/go-bludgeon/employees/logic"
	"github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestLogic(t *testing.T) {
	logger := logger.New()
	meta := memory.New(logger)
	logic := logic.New(logger, meta)
	assert.NotNil(t, logic)
}
