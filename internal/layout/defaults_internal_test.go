package layout

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpEndpoint(t *testing.T) {
	h := httpEndpoint("health")
	assert.Equal(t, "health.go", h.Name())
	assert.Equal(t, fs.FileMode(0644), h.Permissions())
	assert.IsType(t, &template.Template{}, h.Template())
}

func TestHttpServiceNodes(t *testing.T) {
	nodes := httpServiceNodes()

	t.Run("main.go has correct imports", func(t *testing.T) {
		nn := _findNodeByName(nodes, "main.go")
		require.NotNil(t, nn)
		n, ok := nn.(*Fnode)
		require.True(t, ok)

		mod := "cheftest"
		expected := fmt.Sprintf(`import (
	server "%s/server/http"
)`, mod)

		var out bytes.Buffer
		err := n.wbuild(&out, mod)
		assert.NoError(t, err)
		assert.True(t, strings.Contains(out.String(), expected))
	})
}

func _findNodeByName(nodes []Node, name string) Node {
	for _, n := range nodes {
		if n.Name() == name {
			return n
		}
	}
	return nil
}
