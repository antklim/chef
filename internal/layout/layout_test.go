package layout_test

import (
	"os"
	"path"
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertLayout(t *testing.T, root string) {
	{
		// root of the project should include: server, and main.go
		de, err := os.ReadDir(path.Join(root, projectName))
		require.NoError(t, err)
		assert.Len(t, de, 2)
	}

	{
		// root/server should include .gitkeep
		de, err := os.ReadDir(path.Join(root, projectName, "server"))
		require.NoError(t, err)
		assert.Len(t, de, 1)
	}
}

type dirNode struct {
	name        string
	permissions uint32
	children    []layout.Node
}

func (n dirNode) Name() string {
	return n.name
}

func (n dirNode) Permissions() uint32 {
	return n.permissions
}

func (n dirNode) Children() []layout.Node {
	return n.children
}

var server = dirNode{
	name:        "server",
	permissions: 0755,
}

var root = dirNode{
	name:        projectName,
	permissions: 0755,
	children: []layout.Node{
		layout.SrvMain,
		server,
	},
}

const projectName = "XYZ"

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	err = layout.Builder(tmpDir, root)
	require.NoError(t, err)
	assertLayout(t, tmpDir)
}
