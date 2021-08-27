package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentsEmployCmdRunner(t *testing.T) {
	t.Run("fails when project init failed", func(t *testing.T) {
		p := FailedInit(errors.New("some init error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("failes when employ component failed", func(t *testing.T) {
		p := FailedEmployComponent(errors.New("some employ component error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "employ component failed: some employ component error")
	})

	t.Run("successfully employs a component", func(t *testing.T) {
		p := projMock{}
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.NoError(t, err)
	})
}
