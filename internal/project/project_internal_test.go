package project

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	testCases := []struct {
		desc   string
		p      Project
		schema string
	}{
		{
			desc:   "returns default project layout",
			p:      New("test"),
			schema: layout.ServiceLayout,
		},
		{
			desc:   "returns http service layout",
			p:      New("test", WithCategory(CategoryService), WithServer(ServerHTTP)),
			schema: layout.HTTPServiceLayout,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l, err := tC.p.layout()
			require.NoError(t, err)
			assert.Equal(t, tC.schema, l.Schema())
		})
	}

	t.Run("returns error when unknown layout requested", func(t *testing.T) {
		p := New("test", WithCategory("test"))
		l, err := p.layout()
		assert.EqualError(t, err, "not found layout with name test")
		assert.Nil(t, l)
	})
}
