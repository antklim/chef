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
	assert.Equal(t, node, l.FindNode("testNode"))
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
		assert.NotNil(t, l.FindNode("subdir"))
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
		assert.EqualError(t, err, `node "other" not found in layout`)
	})

	t.Run("returns error when adding existing node", func(t *testing.T) {
		nodes := []layout.Node{layout.NewDnode("subdir"), layout.NewFnode("file.txt")}
		l := layout.New("layout", nodes...)

		err := l.AddNode(layout.NewFnode("file.txt"), layout.Root)
		assert.EqualError(t, err, `node "." already has subnode "file.txt"`)
	})
}

func TestLayoutFindNode(t *testing.T) {
	fileNode := layout.NewFnode("file.txt")
	subdNode := layout.NewDnode("subdir", layout.WithSubNodes(fileNode))
	baseNode := layout.NewDnode("base", layout.WithSubNodes(subdNode))
	l := layout.New("layout", baseNode)

	testCases := []struct {
		desc string
		loc  string
		node layout.Node
	}{
		{
			desc: "finds subdir node by location",
			loc:  "base/subdir",
			node: subdNode,
		},
		{
			desc: "finds subdir node by location prefixed with root location",
			loc:  "./base/subdir",
			node: subdNode,
		},
		{
			desc: "finds file node by location",
			loc:  "base/subdir/file.txt",
			node: fileNode,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := l.FindNode(tC.loc)
			assert.Equal(t, tC.node, n)
		})
	}

	t.Run("returns root node", func(t *testing.T) {
		n := l.FindNode(layout.Root)
		assert.NotNil(t, n)
	})

	t.Run("returns nil when no node found in location", func(t *testing.T) {
		n := l.FindNode("foo/bar")
		assert.Nil(t, n)
	})

	t.Run("returns nil for empty location", func(t *testing.T) {
		n := l.FindNode("")
		assert.Nil(t, n)
	})
}

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
