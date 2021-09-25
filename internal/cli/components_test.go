package cli

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/antklim/chef/internal/chef"
	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestComponentsListCmdRunner(t *testing.T) {
	t.Run("fails when project init failed", func(t *testing.T) {
		p := FailedInit(errors.New("some init error"))
		err := componentsListCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("shows a list of registered components", func(t *testing.T) {
		var buf bytes.Buffer
		printout = &buf

		p := projMock{components: []project.Component{
			{
				Name: "handler",
			},
			{
				Name: "test",
			},
		}}
		err := componentsListCmdRunner(p)
		assert.NoError(t, err)

		bufs := buf.String()
		assert.Contains(t, bufs, "handler")
		assert.Contains(t, bufs, "test")
	})
}

func TestComponentsEmployCmdRunner(t *testing.T) {
	t.Run("fails when project init failed", func(t *testing.T) {
		p := FailedInit(errors.New("some init error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("fails when employ component failed", func(t *testing.T) {
		p := FailedEmployComponent(errors.New("some employ component error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, `employ "handler" component failed: some employ component error`)
	})

	t.Run("successfully employs a component", func(t *testing.T) {
		var buf bytes.Buffer
		printout = &buf

		p := projMock{}
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.NoError(t, err)

		assert.Equal(t, "successfully added \"health\" as \"handler\" component\n", buf.String())
	})
}

func TestInitProjectFails(t *testing.T) {
	t.Run("when no notation file found in working directory", func(t *testing.T) {
		p, err := initProject()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open notation:")
		assert.Nil(t, p)
	})

	t.Run("when notation file corrupted", func(t *testing.T) {
		t.Cleanup(func() {
			os.Remove(chef.DefaultNotationFileName)
		})
		err := os.WriteFile(chef.DefaultNotationFileName, []byte(`foo`), 0600)
		assert.NoError(t, err)

		p, err := initProject()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read notation:")
		assert.Nil(t, p)
	})
}
