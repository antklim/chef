package project_test

import (
	"os"
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
		p := project.New("p", project.WithRoot("/r"))
		err := p.Init()
		assert.EqualError(t, err, "set location failed: stat /r: no such file or directory")
	})

	t.Run("inits project", func(t *testing.T) {
		p := project.New("p")
		err := p.Init()
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

// TODO: test default layouts
// func TestDefaultLayouts(t *testing.T) {
// 	defLayouts := map[string][]string{
// 		layout.ServiceLayout:     {"adapter", "app", "handler", "provider", "server", "test"},
// 		layout.HTTPServiceLayout: {"adapter", "app", "handler", "provider", "server", "test", "main.go"},
// 	}

// 	for layoutName, layoutNodes := range defLayouts {
// 		l := layout.Get(layoutName)
// 		require.NotNil(t, l)
// 		for _, n := range layoutNodes {
// 			node := l.Get(n, layout.Root)
// 			assert.NotNil(t, node)
// 		}
// 	}
// }
