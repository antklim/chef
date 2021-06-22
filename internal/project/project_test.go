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
		project.WithCategory(project.CategoryPackage),
		project.WithRoot("/r"),
		project.WithServer(project.ServerGRPC),
	}

	p := project.New(name, opts...)
	assert.Equal(t, name, p.Name())
	loc, err := p.Location()
	require.NoError(t, err)
	assert.Equal(t, "/r/borsch", loc)
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
		assert func(*testing.T)
	}{
		{
			desc: "inits default project in provided directory",
			root: tmpDir,
			assert: func(t *testing.T) {
				assertProjectLayout(t, tmpDir)
			},
		},
		{
			desc: "inits default project in current directory",
			assert: func(t *testing.T) {
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
			require.NoError(t, err)
			tC.assert(t)
		})
	}
}
