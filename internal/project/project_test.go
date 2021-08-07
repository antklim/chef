package project_test

import (
	"os"
	"path"
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	name := "borsch"
	opts := []project.Option{
		project.WithCategory("pkg"),
		project.WithRoot("/r"),
		project.WithServer("grpc"),
	}

	p := project.New(name, opts...)
	assert.Equal(t, name, p.Name())
	loc, err := p.Location()
	require.NoError(t, err)
	assert.Equal(t, "/r/borsch", loc)
}

func TestProjectValidate(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.Mkdir(path.Join(tmpDir, "chefsushi"), 0755)
	require.NoError(t, err)

	karrageFile := path.Join(tmpDir, "karrage")
	_, err = os.Create(karrageFile)
	require.NoError(t, err)

	testCases := []struct {
		desc string
		name string
		opts []project.Option
		err  string
	}{
		{
			desc: "fails when project name is an empty string",
			err:  "project name cannot be empty",
		},
		{
			desc: "fails when project category is unknown",
			name: "cheffoo",
			opts: []project.Option{project.WithCategory("foo")},
			err:  "project category foo is unknown",
		},
		{
			desc: "fails when project server is unknown",
			name: "chefbar",
			opts: []project.Option{project.WithServer("bar")},
			err:  "project server bar is unknown",
		},
		{
			desc: "fails when provided root directory does not exist",
			name: "cheftempura",
			opts: []project.Option{project.WithRoot("tempura")},
			err:  "stat tempura: no such file or directory",
		},
		{
			desc: "fails when provided root directory is not a directory",
			name: "chefkarrage",
			opts: []project.Option{project.WithRoot(karrageFile)},
			err:  karrageFile + " is not a directory",
		},
		{
			desc: "fails when root directory contains file or directory with the project name",
			name: "chefsushi",
			opts: []project.Option{project.WithRoot(tmpDir)},
			err:  `file or directory "chefsushi" already exists`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(tC.name, tC.opts...)
			err := p.Validate()
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestProjectInit(t *testing.T) {
	t.Run("propagates validation errors", func(t *testing.T) {
		p := project.New("")
		err := p.Init()
		assert.EqualError(t, err, "validation failed: project name cannot be empty")
	})

	t.Run("propagates set location errors", func(t *testing.T) {
		t.Skip()
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
	// testCases := []struct {
	// 	desc string
	// }{
	// 	{
	// 		desc: "",
	// 	},
	// }
	// for _, tC := range testCases {
	// 	t.Run(tC.desc, func(t *testing.T) {

	// 	})
	// }

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
