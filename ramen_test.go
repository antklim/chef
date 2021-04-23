package ramen_test

import (
	"os"
	"testing"

	"github.com/antklim/ramen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: creates project home directory in project location

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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := ramen.New()
			err := p.Init(tC.name)
			tC.assert(t, err)
		})
	}
}

// TODO: when no permissions to read/write to project location root then return error
// TODO: when project location root already has a directory equal to project name then return error

func TestProjectLocation(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	testCases := []struct {
		desc   string
		name   string
		assert func(*testing.T, string, error)
	}{
		{
			desc: "project location is cwd/name when no project root provided by user",
			name: "miso",
			assert: func(t *testing.T, loc string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, cwd+"/miso", loc)
			},
		},
		{
			desc: "returns error when cwd contains file or directory with the project name",
			name: "cmd",
			assert: func(t *testing.T, loc string, err error) {
				assert.EqualError(t, err, "file or directory cmd already exists")
				assert.Equal(t, loc, "")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			loc, err := ramen.Location(tC.name)
			tC.assert(t, loc, err)
		})
	}
}
