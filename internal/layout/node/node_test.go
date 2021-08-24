package node_test

import (
	"os"
	"path"
	"testing"

	"github.com/antklim/chef/internal/layout/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDnodeGetSubNode(t *testing.T) {
	fnode := node.NewFnode("file.txt")
	dnode := node.NewDnode("dnode", node.WithSubNodes(fnode))

	testCases := []struct {
		desc     string
		name     string
		expected node.Node
	}{
		{
			desc:     "returns sub node by name",
			name:     "file.txt",
			expected: fnode,
		},
		{
			desc:     "returns nil when node not found",
			name:     "other-file.txt",
			expected: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := dnode.Get(tC.name)
			assert.Equal(t, tC.expected, n)
		})
	}
}

func TestDnodeAdd(t *testing.T) {
	fnode := node.NewFnode("file.txt")
	dnode := node.NewDnode("dnode", node.WithSubNodes(fnode))

	t.Run("returns an error when existing sub node has same name as the new", func(t *testing.T) {
		subnodesBefore := len(dnode.Nodes())

		newNode := node.NewDnode("file.txt")
		err := dnode.Add(newNode)
		assert.EqualError(t, err, `node "file.txt" already exists`)

		subnodesAfter := len(dnode.Nodes())
		assert.Equal(t, subnodesBefore, subnodesAfter)
	})

	t.Run("adds a new subnode", func(t *testing.T) {
		subnodesBefore := len(dnode.Nodes())

		newNode := node.NewFnode("file2.txt")
		err := dnode.Add(newNode)
		assert.NoError(t, err)

		subnodesAfter := len(dnode.Nodes())
		assert.Equal(t, subnodesBefore+1, subnodesAfter)
	})
}

func TestDnodeBuild(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates node directory in a provided location", func(t *testing.T) {
		n := node.NewDnode("test_dir_1")
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name()))
		assert.NoError(t, err)
	})

	t.Run("creates a directory subnode", func(t *testing.T) {
		sn := node.NewDnode("sub_test_dir_2")
		n := node.NewDnode("test_dir_2", node.WithSubNodes(sn))
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})

	t.Run("creates a file subnode", func(t *testing.T) {
		sn := node.NewFnode("test_file_1", node.WithNewTemplate("test", "package foo"))
		n := node.NewDnode("test_dir_3", node.WithSubNodes(sn))
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadFile(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})
}

func TestFnodeBuild(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("returns an error when does not have template", func(t *testing.T) {
		f := node.NewFnode("test_file_1")
		err := f.Build(tmpDir, "module_name")
		assert.EqualError(t, err, "node template is nil")

		_, err = os.ReadFile(path.Join(tmpDir, f.Name()))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("creates a file using node template", func(t *testing.T) {
		f := node.NewFnode("test_file_2", node.WithNewTemplate("test", "package foo"))
		err := f.Build(tmpDir, "module_name")
		require.NoError(t, err)

		expected := "package foo"

		data, err := os.ReadFile(path.Join(tmpDir, f.Name()))
		require.NoError(t, err)
		assert.Equal(t, expected, string(data))
	})
}
