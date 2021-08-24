package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLayout(t *testing.T) {
	n := newTestNode("foo")
	nodes := []node.Node{n}
	l := layout.New(nodes...)
	assert.Equal(t, n, l.FindNode("foo"))
}

func TestLayoutBuild(t *testing.T) {
	testCases := []struct {
		desc   string
		node   *testNode
		loc    string
		assert func(*testing.T, error)
	}{
		{
			desc: "builds layout nodes",
			node: newTestNode("bar"),
			loc:  "/tmp/foo/bar",
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			desc: "fails to build layout when node build fails",
			node: newTestNode("bar"),
			loc:  "/error/bar",
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "node build error")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := layout.New(tC.node)
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
		dnode := node.NewDnode("subdir")
		l := layout.New()

		err := l.AddNode(dnode, layout.Root)
		assert.NoError(t, err)
		assert.NotNil(t, l.FindNode("subdir"))
	})

	t.Run("adds nodes to a nested level in layout", func(t *testing.T) {
		fnode := node.NewFnode("file.txt")
		dnode := node.NewDnode("dnode", node.WithSubNodes(fnode))
		l := layout.New(dnode)

		err := l.AddNode(node.NewFnode("new_file.txt"), "dnode")
		assert.NoError(t, err)
		assert.Len(t, dnode.Nodes(), 2)
	})

	t.Run("returns error when nested level is a file", func(t *testing.T) {
		fnode := node.NewFnode("file.txt")
		dnode := node.NewDnode("dnode", node.WithSubNodes(fnode))
		l := layout.New(dnode)

		err := l.AddNode(node.NewFnode("new_file.txt"), "dnode/file.txt")
		assert.EqualError(t, err, `node "dnode/file.txt" does not support adding subnodes`)
	})

	t.Run("returns error when nested level not found in layout", func(t *testing.T) {
		fnode := node.NewFnode("file.txt")
		dnode := node.NewDnode("dnode", node.WithSubNodes(fnode))
		l := layout.New(dnode)

		err := l.AddNode(node.NewFnode("new_file.txt"), "other")
		assert.EqualError(t, err, `node "other" not found in layout`)
	})

	t.Run("returns error when adding existing node", func(t *testing.T) {
		nodes := []node.Node{node.NewDnode("subdir"), node.NewFnode("file.txt")}
		l := layout.New(nodes...)

		err := l.AddNode(node.NewFnode("file.txt"), layout.Root)
		assert.EqualError(t, err, `node "." already has subnode "file.txt"`)
	})
}

func TestLayoutFindNode(t *testing.T) {
	fileNode := node.NewFnode("file.txt")
	subdNode := node.NewDnode("subdir", node.WithSubNodes(fileNode))
	baseNode := node.NewDnode("base", node.WithSubNodes(subdNode))
	l := layout.New(baseNode)

	testCases := []struct {
		desc string
		loc  string
		node node.Node
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
