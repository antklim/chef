package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
)

type testLayout struct{}

func (testLayout) Nodes() []layout.Node {
	return nil
}

func (testLayout) Schema() string {
	return "testLayout"
}

var _ layout.Layout = testLayout{}

func TestLayoutRegistry(t *testing.T) {
	t.Run("get returns nil when layout not registered", func(t *testing.T) {
		l := layout.Get("foo")
		assert.Nil(t, l)
	})

	t.Run("get returns layout by schema", func(t *testing.T) {
		tl := testLayout{}
		layout.Register(tl)
		l := layout.Get("testLayout")
		assert.Equal(t, tl, l)
	})

	t.Run("has predefined layouts", func(t *testing.T) {
		defs := []string{"srv", "srv_http"}
		for _, s := range defs {
			l := layout.Get(s)
			assert.NotNil(t, l)
		}
	})
}
