package layout

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLayout struct {
	nodes []Node
}

func (l testLayout) Nodes() []Node {
	return l.nodes
}

func (testLayout) Schema() string {
	return "testLayout"
}

var _ Layout = testLayout{}

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	l := testLayout{nodes: []Node{newdnode("server"), srvMain}}

	err = Builder(tmpDir, "XYZ", l)
	require.NoError(t, err)

	d, err := os.ReadDir(path.Join(tmpDir, "XYZ"))
	require.NoError(t, err)
	assert.Len(t, d, 2)
}
