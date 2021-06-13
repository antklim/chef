package layout

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLayout struct {
	nodes []Node
}

func (l testLayout) Nodes() []Node {
	return l.nodes
}

func (testLayout) Schema() string {
	return "testLayout"
}

var _ Layout = testLayout{}

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	l := testLayout{nodes: []Node{newdnode("server"), srvMain}}

	err = Builder(tmpDir, "XYZ", l)
	require.NoError(t, err)

	d, err := os.ReadDir(path.Join(tmpDir, "XYZ"))
	require.NoError(t, err)
	assert.Len(t, d, 2)
}

func TestDnode(t *testing.T) {
	f1 := fnode{node: node{name: "test_file_1", permissions: 0644}}
	f2 := fnode{node: node{name: "test_file_2", permissions: 0644}}
	f3 := fnode{node: node{name: "test_file_3", permissions: 0644}}
	d1 := dnode{node: node{name: "test_dir_1", permissions: 0755}}

	t.Run("has default directory permissions and no children when created", func(t *testing.T) {
		n := newdnode("test_dir")
		expected := dnode{node: node{name: "test_dir", permissions: 0755}}
		assert.Equal(t, expected, n)
	})

	t.Run("has custom directory permissions when created with permission option", func(t *testing.T) {
		n := newdnode("test_dir", withPermissions(0700))
		expected := dnode{node: node{name: "test_dir", permissions: 0700}}
		assert.Equal(t, expected, n)
	})

	t.Run("has non empty children list when created with children option", func(t *testing.T) {
		n := newdnode("test_dir", withSubNodes(f1, d1))
		expected := dnode{
			node: node{
				name:        "test_dir",
				permissions: 0755,
			},
			subnodes: []Node{
				fnode{node: node{name: "test_file_1", permissions: 0644}},
				dnode{node: node{name: "test_dir_1", permissions: 0755}},
			},
		}
		assert.Equal(t, expected, n)
	})

	t.Run("adds children using AddChildren", func(t *testing.T) {
		n := newdnode("test_dir", withSubNodes(f1, f2))

		n.addSubNodes([]Node{f3})
		n.addSubNodes([]Node{d1})

		expected := []Node{
			fnode{node: node{name: "test_file_1", permissions: 0644}},
			fnode{node: node{name: "test_file_2", permissions: 0644}},
			fnode{node: node{name: "test_file_3", permissions: 0644}},
			dnode{node: node{name: "test_dir_1", permissions: 0755}},
		}
		assert.Equal(t, expected, n.SubNodes())
	})
}
