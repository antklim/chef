package layout

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutBuilder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cheftest")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	l := New("testLayout", []Node{newdnode("server"), srvMain})

	err = Builder(tmpDir, "XYZ", l)
	require.NoError(t, err)

	d, err := os.ReadDir(path.Join(tmpDir, "XYZ"))
	require.NoError(t, err)
	assert.Len(t, d, 2)
}
