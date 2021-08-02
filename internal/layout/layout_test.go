package layout_test

import (
	"errors"
	"io/fs"
	"strings"
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testNode struct {
	buildCalled bool
	buildError  error
	loc         string
}

func (n *testNode) Build(loc, mod string) error {
	n.buildCalled = true
	n.loc = loc
	if strings.HasPrefix(loc, "/error") {
		n.buildError = errors.New("node build error")
		return n.buildError
	}
	return nil
}

func (testNode) Name() string {
	return "testNode"
}

func (testNode) Permissions() fs.FileMode {
	return 0400
}

func (n testNode) WasBuild() bool {
	return n.buildCalled == true
}

func (n testNode) BuiltAt() string {
	return n.loc
}

var _ layout.Node = (*testNode)(nil)

func TestNewLayout(t *testing.T) {
	schema := "testLayout"
	node := &testNode{}
	nodes := []layout.Node{node}
	l := layout.New(schema, nodes...)
	assert.Equal(t, schema, l.Schema())
	assert.Equal(t, node, l.GetNode("testNode", layout.Root))
}

func TestLayoutBuild(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		node   *testNode
		loc    string
		assert func(*testing.T, error)
	}{
		{
			desc: "builds layout nodes",
			name: "successBuild",
			node: &testNode{},
			loc:  "/tmp/foo/bar",
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			desc: "fails to build layout when node build fails",
			name: "errorBuild",
			node: &testNode{},
			loc:  "/error/bar",
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "node build error")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := layout.New(tC.name, tC.node)
			assert.False(t, tC.node.WasBuild())
			err := l.Build(tC.loc, "module_name")
			tC.assert(t, err)
			assert.True(t, tC.node.WasBuild())
			assert.Equal(t, tC.loc, tC.node.BuiltAt())
		})
	}
}

func TestLayoutAddNode(t *testing.T) {
	t.Run("adds nodes to the root level of layout nodes", func(t *testing.T) {
		dnode := layout.NewDnode("subdir")
		l := layout.New("layout")

		err := l.AddNode(dnode, layout.Root)
		assert.NoError(t, err)
		assert.NotNil(t, l.GetNode("subdir", layout.Root))
	})

	t.Run("adds nodes to a nested level in layout", func(t *testing.T) {
		fnode := layout.NewFnode("file.txt")
		dnode := layout.NewDnode("dnode", layout.WithSubNodes(fnode))
		l := layout.New("layout", dnode)

		err := l.AddNode(layout.NewFnode("new_file.txt"), "dnode")
		assert.NoError(t, err)
		assert.Len(t, dnode.Nodes(), 2)
	})

	t.Run("returns error when nested level is a file", func(t *testing.T) {
		fnode := layout.NewFnode("file.txt")
		dnode := layout.NewDnode("dnode", layout.WithSubNodes(fnode))
		l := layout.New("layout", dnode)

		err := l.AddNode(layout.NewFnode("new_file.txt"), "dnode/file.txt")
		assert.EqualError(t, err, `node "dnode/file.txt" does not support adding subnodes`)
	})

	t.Run("returns error when nested level not found in layout", func(t *testing.T) {
		fnode := layout.NewFnode("file.txt")
		dnode := layout.NewDnode("dnode", layout.WithSubNodes(fnode))
		l := layout.New("layout", dnode)

		err := l.AddNode(layout.NewFnode("new_file.txt"), "other")
		assert.EqualError(t, err, `path "other" not found in layout`)
	})

	t.Run("returns error when adding existing node", func(t *testing.T) {
		nodes := []layout.Node{layout.NewDnode("subdir"), layout.NewFnode("file.txt")}
		l := layout.New("layout", nodes...)

		err := l.AddNode(layout.NewFnode("file.txt"), layout.Root)
		assert.EqualError(t, err, `node file.txt already exists at "."`)
	})
}

func TestLayoutGetNode(t *testing.T) {
	fileNode := layout.NewFnode("file.txt")
	subdNode := layout.NewDnode("subdir", layout.WithSubNodes(fileNode))
	baseNode := layout.NewDnode("base", layout.WithSubNodes(subdNode))
	l := layout.New("layout", baseNode)

	testCases := []struct {
		desc     string
		node     string
		loc      string
		expected layout.Node
	}{
		{
			desc:     "returns node from root location",
			node:     "base",
			loc:      ".",
			expected: baseNode,
		},
		{
			desc:     "returns node from subdirectory",
			node:     "subdir",
			loc:      "base",
			expected: subdNode,
		},
		{
			desc:     "returns node from nested subdirectory",
			node:     "file.txt",
			loc:      "base/subdir",
			expected: fileNode,
		},
		{
			desc:     "returns nil when no node found",
			node:     "file.txt",
			loc:      ".",
			expected: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			node := l.GetNode(tC.node, tC.loc)
			assert.Equal(t, tC.expected, node)
		})
	}
}

// func TestAddComponent(t *testing.T) {
// 	tmpl := template.Must(template.New("test").Parse("package foo"))

// 	t.Run("returns error when adding unknown component", func(t *testing.T) {
// 		l := layout.New("layout")
// 		err := l.AddComponent("handler", "health")
// 		assert.EqualError(t, err, `unknown component "handler"`)
// 	})

// 	t.Run("returns error when component node with the provided name already exists", func(t *testing.T) {
// 		loc := "handler"
// 		componentName := "http_handler"
// 		fnode := layout.NewFnode("health.go")
// 		dnode := layout.NewDnode(loc, layout.WithSubNodes(fnode))
// 		l := layout.New("layout", dnode)
// 		err := l.RegisterComponent(componentName, loc, tmpl)
// 		require.NoError(t, err)

// 		err = l.AddComponent(componentName, "health.go")
// 		assert.EqualError(t, err, `http_handler "health.go" already exists`)
// 	})

// 	t.Run("adds a component node", func(t *testing.T) {
// 		loc := "handler"
// 		componentName := "http_handler"
// 		fnode := layout.NewFnode("health.go")
// 		dnode := layout.NewDnode(loc, layout.WithSubNodes(fnode))
// 		l := layout.New("layout", dnode)

// 		err := l.RegisterComponent(componentName, loc, tmpl)
// 		require.NoError(t, err)

// 		nodeName := "info.go"
// 		err = l.AddComponent(componentName, nodeName)
// 		require.NoError(t, err)
// 		assert.NotNil(t, l.GetNode(nodeName, loc))
// 	})
// }

func TestLayoutsRegistry(t *testing.T) {
	t.Run("get returns nil when layout not registered", func(t *testing.T) {
		l := layout.Get("foo")
		assert.Nil(t, l)
	})

	t.Run("get returns layout by schema", func(t *testing.T) {
		tl := layout.New("testLayout", nil)
		layout.Register(tl)
		l := layout.Get("testLayout")
		assert.Equal(t, tl, *l)
	})
}
