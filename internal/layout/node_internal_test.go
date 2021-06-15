package layout

import (
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFNodeBuild(t *testing.T) {
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
