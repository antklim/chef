package project

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectOptions(t *testing.T) {
	testCases := []struct {
		desc     string
		opts     []Option
		expected projectOptions
	}{
		{
			desc: "project created with default options",
			expected: projectOptions{
				root: "",
				cat:  CategoryService,
				srv:  ServerNone,
			},
		},
		{
			desc: "project created with the custom root",
			opts: []Option{WithRoot("/r")},
			expected: projectOptions{
				root: "/r",
				cat:  CategoryService,
				srv:  ServerNone,
			},
		},
		{
			desc: "project created with custom category",
			opts: []Option{WithCategory(CategoryCLI)},
			expected: projectOptions{
				root: "",
				cat:  CategoryCLI,
				srv:  ServerNone,
			},
		},
		{
			desc: "project created with custom server",
			opts: []Option{WithServer(ServerHTTP)},
			expected: projectOptions{
				root: "",
				cat:  CategoryService,
				srv:  ServerHTTP,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := New("test", tC.opts...)
			assert.Equal(t, tC.expected, p.opts)
		})
	}
}

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
