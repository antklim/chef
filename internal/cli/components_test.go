package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentsEmployCmdRunner(t *testing.T) {
	t.Run("returns an error when failed to employ component", func(t *testing.T) {
		p := FailedEmployComponent(errors.New("some employ component error"))
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.EqualError(t, err, "employ component failed: some employ component error")
	})

	t.Run("returns no errors when when successfully employed a component", func(t *testing.T) {
		p := projMock{}
		err := componentsEmployCmdRunner(p, "handler", "health")
		assert.NoError(t, err)
	})
}
