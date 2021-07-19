package layout

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFnode(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("package foo"))

	testCases := []struct {
		desc     string
		opts     []FnodeOption
		expected *Fnode
	}{
		{
			desc:     "has default file permissions and no template when created",
			expected: &Fnode{node: node{name: "test_file", permissions: 0644}},
		},
		{
			desc:     "has custom file permissions when created with permission option",
			opts:     []FnodeOption{WithFperm(0600)},
			expected: &Fnode{node: node{name: "test_file", permissions: 0600}},
		},
		{
			desc: "has custom template when created with new template option",
			opts: []FnodeOption{WithNewTemplate("test_new", "package foo")},
			expected: &Fnode{
				node:     node{name: "test_file", permissions: 0644},
				template: template.Must(template.New("test_new").Parse("package foo")),
			},
		},
		{
			desc: "has custom template when created with template option",
			opts: []FnodeOption{WithTemplate(tmpl)},
			expected: &Fnode{
				node:     node{name: "test_file", permissions: 0644},
				template: template.Must(template.New("test").Parse("package foo")),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := NewFnode("test_file", tC.opts...)
			assert.Equal(t, tC.expected, n)
		})
	}
}

func TestFnodeBuild(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("returns an error when does not have template", func(t *testing.T) {
		f := NewFnode("test_file_1")
		err := f.Build(tmpDir, "module_name")
		assert.EqualError(t, err, "node template is nil")

		_, err = os.ReadFile(path.Join(tmpDir, f.Name()))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("creates a file using node template", func(t *testing.T) {
		f := NewFnode("test_file_2", WithNewTemplate("test", "package foo"))
		err := f.Build(tmpDir, "module_name")
		require.NoError(t, err)

		expected := "package foo"

		data, err := os.ReadFile(path.Join(tmpDir, f.Name()))
		require.NoError(t, err)
		assert.Equal(t, expected, string(data))
	})
}

func TestFnodeWBuild(t *testing.T) {
	t.Run("writes module to template", func(t *testing.T) {
		mod := "cheftest"
		tmpl := template.Must(template.New("test").Parse(`package foo
import "{{ .Module }}/test/template"`))
		expected := fmt.Sprintf(`package foo
import "%s/test/template"`, mod)

		var out bytes.Buffer
		f := NewFnode("test_fnode", WithTemplate(tmpl))
		err := f.wbuild(&out, mod)
		assert.NoError(t, err)
		assert.Equal(t, expected, out.String())
	})
}

func TestDnode(t *testing.T) {
	t.Run("has default directory permissions and no children when created", func(t *testing.T) {
		n := NewDnode("test_dir")
		expected := &Dnode{node: node{name: "test_dir", permissions: 0755}}
		assert.Equal(t, expected, n)
	})

	t.Run("has custom directory permissions when created with permission option", func(t *testing.T) {
		n := NewDnode("test_dir", WithDperm(0700))
		expected := &Dnode{node: node{name: "test_dir", permissions: 0700}}
		assert.Equal(t, expected, n)
	})

	t.Run("has non empty children list when created with children option", func(t *testing.T) {
		f1 := NewFnode("test_file_1")
		d1 := NewDnode("test_dir_1")
		n := NewDnode("test_dir", WithSubNodes(f1, d1))
		expected := &Dnode{
			node: node{
				name:        "test_dir",
				permissions: 0755,
			},
			subnodes: []Node{
				&Fnode{node: node{name: "test_file_1", permissions: 0644}},
				&Dnode{node: node{name: "test_dir_1", permissions: 0755}},
			},
		}
		assert.Equal(t, expected, n)
	})
}

func TestDnodeBuild(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates node directory in a provided location", func(t *testing.T) {
		n := NewDnode("test_dir_1")
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name()))
		assert.NoError(t, err)
	})

	t.Run("creates a directory subnode", func(t *testing.T) {
		sn := NewDnode("sub_test_dir_2")
		n := NewDnode("test_dir_2", WithSubNodes(sn))
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadDir(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})

	t.Run("creates a file subnode", func(t *testing.T) {
		sn := NewFnode("test_file_1", WithNewTemplate("test", "package foo"))
		n := NewDnode("test_dir_3", WithSubNodes(sn))
		err := n.Build(tmpDir, "module_name")
		require.NoError(t, err)

		_, err = os.ReadFile(path.Join(tmpDir, n.Name(), sn.Name()))
		require.NoError(t, err)
	})
}
