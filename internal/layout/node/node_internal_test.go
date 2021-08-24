package node

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

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

func TestFnodeWBuild(t *testing.T) {
	t.Run("writes module to template", func(t *testing.T) {
		data := struct {
			Module string
		}{
			Module: "cheftest",
		}
		tmpl := template.Must(template.New("test").Parse(`package foo
import "{{ .Module }}/test/template"`))
		expected := fmt.Sprintf(`package foo
import "%s/test/template"`, data.Module)

		var out bytes.Buffer
		f := NewFnode("test_fnode", WithTemplate(tmpl))
		err := f.wbuild(&out, data)
		assert.NoError(t, err)
		assert.Equal(t, expected, out.String())
	})
}
