package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmdRunner(t *testing.T) {
	t.Run("returns an error when init failed", func(t *testing.T) {
		p := FailedInit(errors.New("some init error"))
		err := initCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("returns an error when failed to build project", func(t *testing.T) {
		p := FailedBuild(errors.New("some build error"))
		err := initCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some build error")
	})

	t.Run("returns no errors when successfully inited a project", func(t *testing.T) {
		p := projMock{}
		err := initCmdRunner(p)
		assert.NoError(t, err)
	})
}
