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

func TestProjectInit(t *testing.T) {
	t.Run("propagates validation errors", func(t *testing.T) {
		p := project.New("")
		err := p.Init()
		assert.EqualError(t, err, "validation failed: project name cannot be empty")
	})

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
