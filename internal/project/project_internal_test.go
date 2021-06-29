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
				cat:  "srv",
				srv:  "",
			},
		},
		{
			desc: "project created with the custom root",
			opts: []Option{WithRoot("/r")},
			expected: projectOptions{
				root: "/r",
				cat:  "srv",
				srv:  "",
			},
		},
		{
			desc: "project created with custom category",
			opts: []Option{WithCategory("cli")},
			expected: projectOptions{
				root: "",
				cat:  "cli",
				srv:  "",
			},
		},
		{
			desc: "project created with custom server",
			opts: []Option{WithServer("http")},
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "http",
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
			p:      New("test", WithCategory("srv"), WithServer("http")),
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
		assert.EqualError(t, err, "not found layout for category test")
		assert.Nil(t, l)
	})
}
