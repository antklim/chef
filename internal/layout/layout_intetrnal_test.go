package layout

import (
	"testing"
	"text/template"

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
			loc:  Root,
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
	tmpl := template.Must(template.New("test").Parse("package foo"))

	t.Run("new layout does not have registered components", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))
		assert.Empty(t, l.components)
	})

	t.Run("adds component to the list of components", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))

		componentName := "http_handler"
		err := l.RegisterComponent(componentName, "handler", tmpl)
		require.NoError(t, err)
		assert.Contains(t, l.components, componentName)

		componentName = "main.go"
		err = l.RegisterComponent(componentName, Root, tmpl)
		require.NoError(t, err)
		assert.Contains(t, l.components, componentName)
	})

	t.Run("returns error when location does not exist", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))
		componentName := "handler"
		err := l.RegisterComponent(componentName, "other/handler", tmpl)
		assert.EqualError(t, err, `"other/handler" does not exist`)
		assert.NotContains(t, l.components, componentName)
	})

	t.Run("returns error when location is not a directory", func(t *testing.T) {
		l := New("layout", NewFnode("handler"))
		componentName := "handler"
		err := l.RegisterComponent(componentName, "handler", tmpl)
		assert.EqualError(t, err, `"handler" not a directory`)
		assert.NotContains(t, l.components, componentName)
	})

	t.Run("returns error when template is nil", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))
		componentName := "handler"
		err := l.RegisterComponent(componentName, "handler", nil)
		require.EqualError(t, err, "component template is nil")
		assert.NotContains(t, l.components, componentName)
	})

	t.Run("registers components", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))
		componentName := "http_handler"
		err := l.RegisterComponent(componentName, "handler", tmpl)
		require.NoError(t, err)

		// registers other component to the same location
		componentName = "grpc_hander"
		err = l.RegisterComponent(componentName, "handler", tmpl)
		require.NoError(t, err)
	})

	t.Run("overrides an existing component", func(t *testing.T) {
		l := New("layout", NewDnode("handler"))
		componentName := "http_handler"
		err := l.RegisterComponent(componentName, "handler", tmpl)
		require.NoError(t, err)

		otherTmpl := template.Must(template.New("test2").Parse("package bar"))
		err = l.RegisterComponent(componentName, "handler", otherTmpl)
		require.NoError(t, err)

		cmp := l.components[componentName]
		assert.Equal(t, otherTmpl, cmp.template)
	})
}
