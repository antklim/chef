package layout

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testProjectName = "XYZ"

func assertLayout(t *testing.T, root string) {
	{
		// root of the project should include: server, and main.go
		de, err := os.ReadDir(path.Join(root, testProjectName))
		require.NoError(t, err)
		assert.Len(t, de, 2)
	}

	{
		// root/server should include .gitkeep
		de, err := os.ReadDir(path.Join(root, testProjectName, "server"))
		require.NoError(t, err)
		assert.Len(t, de, 1)
	}
}

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	var server = dnode{
		name:        "server",
		permissions: dperm,
	}

	var root = dnode{
		name:        testProjectName,
		permissions: dperm,
		children: []Node{
			srvMain,
			server,
		},
	}

	err = Builder(tmpDir, root)
	require.NoError(t, err)
	assertLayout(t, tmpDir)
}
