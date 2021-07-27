package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNode(t *testing.T) {
	handlerNode := NewFnode("handler.go")
	httpNode := NewDnode("http", WithSubNodes(handlerNode))
	serverNode := NewDnode("server", WithSubNodes(httpNode))
	l := New("layout", serverNode)

	testCases := []struct {
		desc     string
		loc      string
		expected Node
	}{
		{
			desc:     "returns server/http node",
			loc:      "server/http",
			expected: httpNode,
		},
		{
			desc:     "returns handler node",
			loc:      "server/http/handler.go",
			expected: handlerNode,
		},
		{
			desc: "returns nil when node does not exist",
			loc:  "server/http/other.go",
		},
		{
			desc: "returns nil for root",
			loc:  ".",
		},
		{
			desc: "returns nil when location does not exist",
			loc:  "",
		},
		{
			desc: "returns nil when location does not exist",
			loc:  "server/grpc",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			node := l.findNode(tC.loc)
			assert.Equal(t, tC.expected, node)
		})
	}
}

func TestRegisterComponent(t *testing.T) {
	t.Run("returns error when provided location does not exist", func(t *testing.T) {
		dnode := NewDnode("handler")
		l := New("layout", dnode)

		componentName := "handler"
		assert.NotContains(t, l.components, componentName)

		err := l.RegisterComponent(componentName, "other/handler", nil)
		assert.EqualError(t, err, `"other/handler" does not exist`)
		assert.NotContains(t, l.components, componentName)
	})

	t.Run("returns error when provided location is not a directory", func(t *testing.T) {
		dnode := NewFnode("handler")
		l := New("layout", dnode)

		componentName := "handler"
		assert.NotContains(t, l.components, componentName)

		err := l.RegisterComponent(componentName, "handler", nil)
		assert.EqualError(t, err, `"handler" not a directory`)
		assert.NotContains(t, l.components, componentName)
	})

	t.Run("returns error when template is nil", func(t *testing.T) {})

	t.Run("registers component", func(t *testing.T) {
		dnode := NewDnode("handler")
		l := New("layout", dnode)

		componentName := "http_hander"
		assert.NotContains(t, l.components, componentName)

		err := l.RegisterComponent("http_hander", "handler", nil)
		require.NoError(t, err)
		assert.Contains(t, l.components, componentName)

		// registers other component to the same location
		componentName = "grpc_hander"
		assert.NotContains(t, l.components, componentName)

		err = l.RegisterComponent("grpc_hander", "handler", nil)
		require.NoError(t, err)
		assert.Contains(t, l.components, componentName)
	})

	t.Run("overrides an existing component", func(t *testing.T) {

	})

	t.Run("registers component at the root of layout", func(t *testing.T) {})
}
