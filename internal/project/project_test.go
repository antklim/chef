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
	testCases := []struct {
		desc string
		name string
		opts []project.Option
		root string
		cat  project.Category
		srv  project.Server
	}{
		{
			desc: "returns default project manager when no options provided",
			name: "ramen",
			cat:  project.CategoryService,
			srv:  project.ServerHTTP,
		},
		{
			desc: "returns project with custom options",
			name: "borsch",
			opts: []project.Option{
				project.WithCategory(project.CategoryPackage),
				project.WithRoot("/r"),
				project.WithServer(project.ServerGRPC),
			},
			root: "/r",
			cat:  project.CategoryPackage,
			srv:  project.ServerGRPC,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(tC.name, tC.opts...)
			assert.Equal(t, tC.name, p.Name())
			assert.Equal(t, tC.root, p.Root())
			assert.Equal(t, tC.cat, p.Category())
			assert.Equal(t, tC.srv, p.Server())
		})
	}
}

func TestProjectValidate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	err = os.Mkdir(path.Join(tmpDir, "chefsushi"), 0755)
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
			desc: "fails when provided root does not exist",
			name: "cheftempura",
			opts: []project.Option{project.WithRoot("tempura")},
			err:  "stat tempura: no such file or directory",
		},
		{
			desc: "fails when provided root is not a directory",
			name: "chefkarrage",
			opts: []project.Option{project.WithRoot(karrageFile)},
			err:  karrageFile + " is not a directory",
		},
		{
			desc: "fails when root contains file or directory with the project name",
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

func assertProjectLayout(t *testing.T, root string) {
	// root of the project should include: adapter, app, handler, provider, test and main.go
	d, err := os.ReadDir(path.Join(root, "cheftest"))
	require.NoError(t, err)
	assert.Len(t, d, 7)
}

func TestProjectInit(t *testing.T) {
	name := "cheftest"
	tmpDir, err := os.MkdirTemp("", name)
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	defer os.RemoveAll(name)

	testCases := []struct {
		desc   string
		root   string
		assert func(*testing.T, error)
	}{
		{
			desc: "inits default project in provided directory",
			root: tmpDir,
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
				assertProjectLayout(t, tmpDir)
			},
		},
		{
			desc: "inits default project in current directory",
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
				assertProjectLayout(t, "")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(name, project.WithRoot(tC.root))
			err := p.Validate()
			require.NoError(t, err)

			err = p.Init()
			tC.assert(t, err)
		})
	}
}
