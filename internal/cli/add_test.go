package cli_test

import (
	"errors"
	"testing"

	"github.com/antklim/chef/internal/cli"
	"github.com/antklim/chef/internal/cli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddCmdRunner(t *testing.T) {
	t.Run("returns an error when failed to add component to a project", func(t *testing.T) {
		p := mocks.FailedProject(errors.New("some add component error"))
		err := cli.AddCmdRunner(p)
		assert.EqualError(t, err, "unable to add to a project: some add component error")
	})

	t.Run("returns no errors when when successfully added a component to a project", func(t *testing.T) {
		p := mocks.Project{}
		err := cli.AddCmdRunner(p)
		assert.NoError(t, err)
	})
}
