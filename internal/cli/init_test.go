package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmdRunner(t *testing.T) {
	t.Run("fails when project init failed", func(t *testing.T) {
		p := FailedInit(errors.New("some init error"))
		err := initCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some init error")
	})

	t.Run("fails when project build failed", func(t *testing.T) {
		p := FailedBuild(errors.New("some build error"))
		err := initCmdRunner(p)
		assert.EqualError(t, err, "init project failed: some build error")
	})

	t.Run("successfully inits a project", func(t *testing.T) {
		var buf bytes.Buffer
		printout = &buf

		p := projMock{loc: "project_location"}
		err := initCmdRunner(p)
		assert.NoError(t, err)

		assert.Contains(t, buf.String(), "project successfully inited at project_location\n")
	})
}
