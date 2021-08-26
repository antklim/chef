package node_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/antklim/chef/internal/layout/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDnodeGetSubNode(t *testing.T) {
	f := node.NewFnode("file.txt")
	d := node.NewDnode("dir", node.WithSubNodes(f))

	testCases := []struct {
		desc     string
		name     string
		expected node.Node
	}{
		{
			desc:     "returns sub node by name",
			name:     "file.txt",
			expected: f,
		},
		{
			desc:     "returns nil when node not found",
			name:     "other-file.txt",
			expected: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := d.Get(tC.name)
			assert.Equal(t, tC.expected, n)
		})
	}
}

func TestDnodeAdd(t *testing.T) {
	f := node.NewFnode("file.txt")
	d := node.NewDnode("dir", node.WithSubNodes(f))

	t.Run("returns an error when existing sub node has same name as the new", func(t *testing.T) {
		subnodesBefore := len(d.Nodes())

		newNode := node.NewDnode("file.txt")
		err := d.Add(newNode)
		assert.EqualError(t, err, `node "file.txt" already exists`)

		subnodesAfter := len(d.Nodes())
		assert.Equal(t, subnodesBefore, subnodesAfter)
	})

	t.Run("adds a new subnode", func(t *testing.T) {
		subnodesBefore := len(d.Nodes())

		newNode := node.NewFnode("file2.txt")
		err := d.Add(newNode)
		assert.NoError(t, err)

		subnodesAfter := len(d.Nodes())
		assert.Equal(t, subnodesBefore+1, subnodesAfter)
	})
}

func TestDnodeBuild(t *testing.T) {
	t.Run("creates node directory in a provided location", func(t *testing.T) {
		tmpDir := t.TempDir()
		d := node.NewDnode("dir")
		err := d.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, d.Name()))
		assert.NoError(t, err)
	})

	t.Run("creates a directory subnode", func(t *testing.T) {
		tmpDir := t.TempDir()
		sd := node.NewDnode("subdir")
		d := node.NewDnode("dir", node.WithSubNodes(sd))
		err := d.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, d.Name(), sd.Name()))
		require.NoError(t, err)
	})

	t.Run("creates a file subnode", func(t *testing.T) {
		tmpDir := t.TempDir()
		f := node.NewFnode("file.go", node.WithNewTemplate("test", "package foo"))
		d := node.NewDnode("dir", node.WithSubNodes(f))
		err := d.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadFile(path.Join(tmpDir, d.Name(), f.Name()))
		require.NoError(t, err)
	})

	t.Run("fails when subnode build fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		f := node.NewFnode("file.go")
		d := node.NewDnode("dir", node.WithSubNodes(f))
		err := d.Build(tmpDir, "module_name")
		require.EqualError(t, err, `failed to build subnode "file.go": node template is nil`)

		_, err = os.ReadFile(path.Join(tmpDir, d.Name(), f.Name()))
		assert.True(t, os.IsNotExist(err))
	})
}

func TestFnodeBuild(t *testing.T) {
	t.Run("fails when does not have template", func(t *testing.T) {
		tmpDir := t.TempDir()
		f := node.NewFnode("file.go")
		err := f.Build(tmpDir, "module_name")
		assert.EqualError(t, err, "node template is nil")

		_, err = os.ReadFile(path.Join(tmpDir, f.Name()))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("fails when cannot execute template", func(t *testing.T) {
		tmpDir := t.TempDir()
		f := node.NewFnode("file.go", node.WithNewTemplate("test", "package foo {{ .Foo }}"))
		err := f.Build(tmpDir, "module_name")
		assert.Error(t, err)
		isValidError := strings.HasPrefix(err.Error(), "failed to execute template")
		assert.True(t, isValidError)

		_, err = os.ReadFile(path.Join(tmpDir, f.Name()))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("creates a file using node template", func(t *testing.T) {
		tmpDir := t.TempDir()
		f := node.NewFnode("file.go", node.WithNewTemplate("test", "package foo"))
		err := f.Build(tmpDir, "module_name")
		require.NoError(t, err)

		expected := "package foo"

		data, err := os.ReadFile(path.Join(tmpDir, f.Name()))
		require.NoError(t, err)
		assert.Equal(t, expected, string(data))
	})
}
