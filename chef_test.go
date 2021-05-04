package chef_test

import (
	"os"
	"path"
	"testing"

	"github.com/antklim/chef"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	testCases := []struct {
		desc string
		name string
		opts []chef.Option
		proj chef.Project
	}{
		{
			desc: "returns default project manager when no options provided",
			name: "ramen",
			proj: chef.Project{
				Name:     "ramen",
				Category: "app",
				Server:   "http",
			},
		},
		{
			desc: "returns project with custom options",
			opts: []chef.Option{
				chef.WithCategory(chef.CategoryPkg),
				chef.WithRoot("/r"),
				chef.WithServer(chef.ServerGRPC),
			},
			name: "borsch",
			proj: chef.Project{
				Name:     "borsch",
				Root:     "/r",
				Category: "pkg",
				Server:   "grpc",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := chef.New(tC.name, tC.opts...)
			assert.Equal(t, tC.proj, p)
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
		opts []chef.Option
		err  string
	}{
		{
			desc: "fails when project name is an empty string",
			err:  "project name required: empty name provided",
		},
		{
			desc: "fails when provided root does not exist",
			name: "cheftempura",
			opts: []chef.Option{chef.WithRoot("tempura")},
			err:  "stat tempura: no such file or directory",
		},
		{
			desc: "fails when provided root is not a directory",
			name: "chefkarrage",
			opts: []chef.Option{chef.WithRoot(karrageFile)},
			err:  karrageFile + " is not a directory",
		},
		{
			desc: "fails when root contains file or directory with the project name",
			name: "chefsushi",
			opts: []chef.Option{chef.WithRoot(tmpDir)},
			err:  "file or directory chefsushi already exists",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := chef.New(tC.name, tC.opts...)
			err := p.Validate()
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestProjectInit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	testCases := []struct {
		desc   string
		name   string
		root   string
		assert func(*testing.T, error)
	}{
		{
			desc: "inits default project in provided directory",
			name: "cheftest",
			root: tmpDir,
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
				{
					// root of the project should include cmd, internal, test
					de, err := os.ReadDir(path.Join(tmpDir, "cheftest"))
					require.NoError(t, err)
					assert.Len(t, de, 3)
				}

				{
					// root/cmd should include main.go
					de, err := os.ReadDir(path.Join(tmpDir, "cheftest", "cmd"))
					require.NoError(t, err)
					assert.Len(t, de, 1)
				}

				{
					// root/internal should include app, adapter, provider, and server
					de, err := os.ReadDir(path.Join(tmpDir, "cheftest", "internal"))
					require.NoError(t, err)
					assert.Len(t, de, 4)
				}

				{
					// root/internal/server should include http
					de, err := os.ReadDir(path.Join(tmpDir, "cheftest", "internal", "server"))
					require.NoError(t, err)
					assert.Len(t, de, 1)
				}
				// TODO: all leaf directories should contain .gitkeep
			},
		},
		// TODO: inits default project in current directory
		// {
		// 	desc: "inits default project in current directory",
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := chef.New(tC.name, chef.WithRoot(tC.root))
			err := p.Validate()
			require.NoError(t, err)

			err = p.Init()
			tC.assert(t, err)
		})
	}
}