package project_test

import (
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
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
	// builds project in provided directory
	// builds project in current directory

	t.Run("fails when", func(t *testing.T) {
		testCases := []struct {
			desc string
			pf   func() *project.Project
			err  string
		}{
			{
				desc: "project does not have layout",
				pf: func() *project.Project {
					return project.New("project")
				},
				err: "project does not have layout",
			},
			// {
			// 	desc: "project could not be build",
			// 	pf: func() project.Project {
			// 		return project.New("project")
			// 	},
			// 	err: "---",
			// },
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				loc, err := tC.pf().Build()
				assert.EqualError(t, err, tC.err)
				assert.Empty(t, loc)
			})
		}
	})

	t.Run("builds project", func(t *testing.T) {})
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
