package layout

import (
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFnodeBuild(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	buildfnode := func(name, t string) (fnode, error) {
		var tmpl *template.Template

		if t != "" {
			tmpl = template.Must(template.New("test").Parse(t))
		}

		f := fnode{node: node{name: name, permissions: 0644}, template: tmpl}
		return f, f.Build(tmpDir)
	}

	t.Run("creates nothing when does not have template", func(t *testing.T) {
		f, err := buildfnode("test_file_1", "")
		require.NoError(t, err)

		_, err = os.ReadFile(path.Join(tmpDir, f.Name()))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("creates a file using node template", func(t *testing.T) {
		f, err := buildfnode("test_file_2", `package foo`)
		require.NoError(t, err)

		expected := "package foo"

		data, err := os.ReadFile(path.Join(tmpDir, f.Name()))
		require.NoError(t, err)
		assert.Equal(t, expected, string(data))
	})
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

func TestDnodeBuild(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	t.Run("creates node directory in a provided location", func(t *testing.T) {
		n := newdnode("test_dir_1")
		err := n.Build(tmpDir)
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name()))
		assert.NoError(t, err)
	})

	t.Run("creates a directory subnode", func(t *testing.T) {
		sn := newdnode("sub_test_dir_2")
		n := newdnode("test_dir_2", withSubNodes(sn))
		err := n.Build(tmpDir)
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})

	t.Run("creates a file subnode", func(t *testing.T) {
		sn := fnode{
			node:     node{name: "test_file_1", permissions: 0644},
			template: template.Must(template.New("test").Parse("package foo")),
		}
		n := newdnode("test_dir_3", withSubNodes(sn))
		err := n.Build(tmpDir)
		require.NoError(t, err)

		_, err = os.ReadFile(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})
}
