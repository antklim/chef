package cli_test

import (
	"errors"
	"testing"

	"github.com/antklim/chef/internal/cli"
	"github.com/antklim/chef/internal/cli/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapCmdRunner(t *testing.T) {
	t.Run("returns an error when bootstrap failed", func(t *testing.T) {
		p := mocks.FailedProject(errors.New("some bootstrap error"))
		err := cli.BootstrapCmdRunner(p)
		assert.EqualError(t, err, "unable to bootstrap project: some bootstrap error")
	})

	t.Run("returns no errors when successfully bootstrapped a project", func(t *testing.T) {
		p := mocks.Project{}
		err := cli.BootstrapCmdRunner(p)
		assert.NoError(t, err)
	})
}
