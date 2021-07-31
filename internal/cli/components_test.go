package cli

import (
	"errors"
	"testing"

	"github.com/antklim/chef/internal/cli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddCmdRunner(t *testing.T) {
	t.Run("returns an error when failed to add component to a project", func(t *testing.T) {
		p := mocks.FailedProject(errors.New("some add component error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "could not add project component: some add component error")
	})

	t.Run("returns no errors when when successfully added a component to a project", func(t *testing.T) {
		p := mocks.Project{}
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.NoError(t, err)
	})
}
