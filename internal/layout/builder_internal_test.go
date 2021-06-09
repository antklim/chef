package layout

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testProjectName = "XYZ"

func assertLayout(t *testing.T, root string) {
	{
		// root of the project should include: server, and main.go
		de, err := os.ReadDir(path.Join(root, testProjectName))
		require.NoError(t, err)
		assert.Len(t, de, 2)
	}

	{
		// root/server should include .gitkeep
		de, err := os.ReadDir(path.Join(root, testProjectName, "server"))
		require.NoError(t, err)
		assert.Len(t, de, 1)
	}
}

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	server := newdnode("server")
	root := newdnode(testProjectName, withSubNodes(srvMain, server))

	err = Builder(tmpDir, root)
	require.NoError(t, err)
	assertLayout(t, tmpDir)
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

func TestLayoutSelector(t *testing.T) {
	t.Run("returns unknown category error when unknown category provided", func(t *testing.T) {})

	t.Run("returns default service layout for service category", func(t *testing.T) {})

	t.Run("returns http service layout for service category and http server", func(t *testing.T) {})

	t.Run("returns unknown server error for service category and unknown server", func(t *testing.T) {})
}
