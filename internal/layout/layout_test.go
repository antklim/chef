package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/stretchr/testify/assert"
)

func TestNewLayout(t *testing.T) {
	n := newTestNode("foo")
	nodes := []node.Node{n}
	l := layout.New(nodes...)
	assert.Equal(t, n, l.FindNode("foo"))
}

func TestLayoutBuild(t *testing.T) {
	testCases := []struct {
		desc string
		node *testNode
		loc  string
		err  string
	}{
		{
			desc: "builds layout nodes",
			node: newTestNode("bar"),
			loc:  "/tmp/foo/bar",
		},
		{
			desc: "fails when node build fails",
			node: newTestNode("bar"),
			loc:  "/error/bar",
			err:  "node build error",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := layout.New(tC.node)
			assert.False(t, tC.node.WasBuild())

			err := l.Build(tC.loc, "module_name")
			if tC.err != "" {
				assert.EqualError(t, err, tC.err)
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, tC.node.WasBuild())
			assert.Equal(t, tC.loc, tC.node.BuiltAt())
		})
	}
}

func TestLayoutAddNodeFails(t *testing.T) {
	/* Test layout:
	  .
		+- dir
		   +- file.txt
	*/

	f := node.NewFnode("file.txt")
	d := node.NewDnode("dir", node.WithSubNodes(f))
	l := layout.New(d)

	testCases := []struct {
		desc string
		node node.Node
		loc  string
		err  string
	}{
		{
			desc: "when nested level is a file",
			node: node.NewFnode("new_file.txt"),
			loc:  "dir/file.txt",
			err:  `node "dir/file.txt" does not support adding subnodes`,
		},
		{
			desc: "when nested level not found in layout",
			node: node.NewFnode("new_file.txt"),
			loc:  "other",
			err:  `node "other" not found in layout`,
		},
		{
			desc: "when adding existing node",
			node: node.NewDnode("dir"),
			loc:  layout.Root,
			err:  `node "." already has subnode "dir"`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := l.AddNode(tC.node, tC.loc)
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestLayoutAddNode(t *testing.T) {
	t.Run("adds nodes to the root level of layout nodes", func(t *testing.T) {
		d := node.NewDnode("dir")
		l := layout.New()

		assert.Nil(t, l.FindNode("dir"))

		err := l.AddNode(d, layout.Root)
		assert.NoError(t, err)

		assert.NotNil(t, l.FindNode("dir"))
	})

	t.Run("adds nodes to a nested level in layout", func(t *testing.T) {
		f := node.NewFnode("file.txt")
		d := node.NewDnode("dir")
		l := layout.New(d)

		assert.Empty(t, d.Nodes())
		assert.Nil(t, l.FindNode("dir/file.txt"))

		err := l.AddNode(f, "dir")
		assert.NoError(t, err)

		assert.NotEmpty(t, d.Nodes())
		assert.NotNil(t, l.FindNode("dir/file.txt"))
	})
}

func TestLayoutFindNode(t *testing.T) {
	/* Test layout:
	  .
		+- dir
		   +- file1.txt
		   +- subdir
		      +- file2.txt
	*/
	f1, f2 := node.NewFnode("file1.txt"), node.NewFnode("file2.txt")
	sd := node.NewDnode("subdir", node.WithSubNodes(f2))
	d := node.NewDnode("dir", node.WithSubNodes(f1, sd))
	l := layout.New(d)

	testCases := []struct {
		desc string
		loc  string
		node node.Node
	}{
		{
			desc: "finds subdir node by location",
			loc:  "dir/subdir",
			node: sd,
		},
		{
			desc: "finds subdir node by location prefixed with root location",
			loc:  "./dir/subdir",
			node: sd,
		},
		{
			desc: "finds file node by location",
			loc:  "dir/subdir/file2.txt",
			node: f2,
		},
		{
			desc: "returns nil when no node found in location",
			loc:  "foo/bar",
		},
		{
			desc: "returns nil when subnode is a file",
			loc:  "dir/file1.txt/foo",
		},
		{
			desc: "returns nil for empty location",
			loc:  "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := l.FindNode(tC.loc)
			if tC.node != nil {
				assert.Equal(t, tC.node, n)
			} else {
				assert.Nil(t, n)
			}
		})
	}

	t.Run("returns root node", func(t *testing.T) {
		n := l.FindNode(layout.Root)
		assert.NotNil(t, n)
	})
}
