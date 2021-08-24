package project_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectInitFails(t *testing.T) {
	tmpDir := t.TempDir()
	foofile := path.Join(tmpDir, "foofoo")
	_, err := os.Create(foofile)
	require.NoError(t, err)

	testCases := []struct {
		desc string
		name string
		opts []project.Option
		err  string
	}{
		{
			desc: "when the project name is an empty string",
			err:  "validation failed: name cannot be empty",
		},
		{
			desc: "when the project category is unknown",
			name: "cheftest",
			opts: []project.Option{project.WithCategory("foo")},
			err:  `validation failed: unknown category "foo"`,
		},
		{
			desc: "when the project server is unknown",
			name: "cheftest",
			opts: []project.Option{project.WithServer("foo")},
			err:  `validation failed: unknown server "foo"`,
		},
		{
			desc: "when root directory does not exist",
			name: "cheftest",
			opts: []project.Option{project.WithRoot("foo")},
			err:  "set location failed: stat foo: no such file or directory",
		},
		{
			desc: "when root is not a directory",
			name: "cheftest",
			opts: []project.Option{project.WithRoot(foofile)},
			err:  `set location failed: "` + foofile + `" is not a directory`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(tC.name, tC.opts...)
			err := p.Init()
			assert.EqualError(t, err, tC.err)
		})
	}

	// t.Run("propagates set location errors", func(t *testing.T) {})
	// t.Run("propagates set layout errors", func(t *testing.T) {})
}

func TestProjectInit(t *testing.T) {
	t.Run("propagates set location errors", func(t *testing.T) {
		p := project.New("project", project.WithRoot("/r"))
		err := p.Init()
		assert.EqualError(t, err, "set location failed: stat /r: no such file or directory")
	})

	t.Run("inits project", func(t *testing.T) {
		p := project.New("project")
		err := p.Init()
		assert.NoError(t, err)
	})

	t.Run("inits existing project", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := os.Mkdir(path.Join(tmpDir, "project"), 0755)
		require.NoError(t, err)

		p := project.New("project")
		err = p.Init()
		assert.NoError(t, err)
	})
}

func TestProjectBuild(t *testing.T) {
	t.Run("returns error when project not inited", func(t *testing.T) {
		p := project.New("project")
		loc, err := p.Build()
		assert.EqualError(t, err, "project not inited")
		assert.Empty(t, loc)
	})

	t.Run("returns error when root directory contains file or directory with the project name", func(t *testing.T) {
		tmpDir := t.TempDir()
		ppath := path.Join(tmpDir, "project")
		err := os.Mkdir(ppath, 0755)
		require.NoError(t, err)

		p := project.New("project", project.WithRoot(tmpDir))
		err = p.Init()
		require.NoError(t, err)

		expected := fmt.Sprintf("build failed: mkdir %s: file exists", ppath)
		_, err = p.Build()
		assert.EqualError(t, err, expected)
	})

	t.Run("builds project", func(t *testing.T) {
		tmpDir := t.TempDir()
		p := project.New("project", project.WithRoot(tmpDir))
		err := p.Init()
		require.NoError(t, err)

		loc, err := p.Build()
		require.NoError(t, err)

		nodes, err := os.ReadDir(loc)
		require.NoError(t, err)
		assert.NotEmpty(t, nodes)
	})
}
