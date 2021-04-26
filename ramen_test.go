package ramen_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/antklim/ramen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	testCases := []struct {
		desc string
		proj *ramen.Project
	}{
		{
			desc: "returns default project manager when no options provided",
			proj: &ramen.Project{
				Name:   "ramen",
				Taste:  "app",
				Server: "http",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := ramen.New()
			assert.Equal(t, tC.proj, p)
		})
	}
}

// TODO: move name validation to a New method
// TODO: creates project home directory in project root location
func TestProjectInit(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		assert func(*testing.T, error)
	}{
		{
			desc: "project name required",
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "project name required")
			},
		},
		// {
		// 	desc: "inits default project",
		// 	assert: func(t *testing.T, err error) {
		// 		require.NoError(t, err)
		// 	},
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := ramen.New()
			err := p.Init(tC.name)
			tC.assert(t, err)
		})
	}
}

func TestProjectLocation(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ramentest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	err = os.Mkdir(path.Join(tmpDir, "sushi"), 0755)
	require.NoError(t, err)

	testCases := []struct {
		desc   string
		name   string
		root   string
		assert func(*testing.T, string, error)
	}{
		{
			desc: "is root/name when project root provided by user",
			name: "miso",
			root: tmpDir,
			assert: func(t *testing.T, loc string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, path.Join(tmpDir, "miso"), loc)
			},
		},
		{
			desc: "fails when root contains file or directory with the project name",
			name: "sushi",
			root: tmpDir,
			assert: func(t *testing.T, loc string, err error) {
				assert.EqualError(t, err, "file or directory sushi already exists")
				assert.Equal(t, loc, "")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			loc, err := ramen.Location(tC.name, tC.root)
			tC.assert(t, loc, err)
		})
	}
}

func TestProjectRoot(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ramentest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	sushiDir := path.Join(tmpDir, "sushi")
	err = os.Mkdir(sushiDir, 0755)
	require.NoError(t, err)

	karrageFile := path.Join(tmpDir, "karrage")
	_, err = os.Create(karrageFile)
	require.NoError(t, err)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	testCases := []struct {
		desc   string
		name   string
		assert func(*testing.T, string, error)
	}{
		{
			desc: "is current working drrectory when no root provided",
			assert: func(t *testing.T, root string, err error) {
				require.NoError(t, err)
				assert.Equal(t, cwd, root)
			},
		},
		{
			desc: "is root directory when provided",
			name: sushiDir,
			assert: func(t *testing.T, root string, err error) {
				require.NoError(t, err)
				assert.Equal(t, sushiDir, root)
			},
		},
		{
			desc: "fails when provided root does not exist",
			name: "tempura",
			assert: func(t *testing.T, root string, err error) {
				require.EqualError(t, err, "stat tempura: no such file or directory")
				assert.Equal(t, "", root)
			},
		},
		{
			desc: "fails when provided root is not a directory",
			name: karrageFile,
			assert: func(t *testing.T, root string, err error) {
				require.EqualError(t, err, karrageFile+" is not a directory")
				assert.Equal(t, "", root)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			root, err := ramen.Root(tC.name)
			tC.assert(t, root, err)
		})
	}
}
