package cli

import (
	"errors"
	"testing"

	"github.com/antklim/chef/internal/cli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestComponentsEmployCmdRunner(t *testing.T) {
	t.Run("returns an error when failed to employ component", func(t *testing.T) {
		p := mocks.FailedProject(errors.New("some employ component error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "project employ component failed: some employ component error")
	})

	t.Run("returns no errors when when successfully employed a component", func(t *testing.T) {
		p := mocks.Project{}
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.NoError(t, err)
	})
}
