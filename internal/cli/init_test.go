package cli

import (
	"errors"
	"testing"

	"github.com/antklim/chef/internal/cli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInitCmdRunner(t *testing.T) {
	t.Run("returns an error when init failed", func(t *testing.T) {
		p := mocks.FailedProject(errors.New("some init error"))
		err := initCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("returns no errors when successfully inited a project", func(t *testing.T) {
		p := mocks.Project{}
		err := initCmdRunner(p)
		assert.NoError(t, err)
	})
}
