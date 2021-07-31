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
			err:  "project name required: empty name provided",
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
			err:  "file or directory chefsushi already exists",
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

func TestProjectBootstrap(t *testing.T) {
	testCases := []struct {
		desc string
		root string
	}{
		{
			desc: "inits project in provided directory",
			root: t.TempDir(),
		},
		{
			desc: "inits project in current directory",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New("cheftest", project.WithRoot(tC.root))
			err := p.Init()
			require.NoError(t, err)

			loc, err := p.Location()
			require.NoError(t, err)

			_, err = os.ReadDir(loc)
			require.NoError(t, err)

			os.RemoveAll(loc)
		})
	}
}

func TestProjectAdd(t *testing.T) {
	t.Run("returns error when trying to add unknow component type", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		p := project.New("cheftest")
		err := p.Add("foo", "bar")
		assert.EqualError(t, err, `could not add layout component: unknown component "foo"`)
	})

	t.Run("returns error when component with the given name already exists", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		// t.Log("not implemented")
		// t.Fail()
	})

	t.Run("adds new component node to a project layout", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		// t.Log("not implemented")
		// t.Fail()
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
