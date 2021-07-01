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

func TestHttpHandler(t *testing.T) {
	h := httpHandler("health")
	assert.Equal(t, "health.go", h.Name())
	assert.Equal(t, fs.FileMode(0644), h.Permissions())
	assert.IsType(t, &template.Template{}, h.Template())
}

func TestServiceNodes(t *testing.T) {
	nodes := serviceNodes()
	assert.Len(t, nodes, 6)
	expectedNodes := []string{"adapter", "app", "handler", "provider", "server", "test"}
	for _, n := range expectedNodes {
		hasNode := _hasNodeWithName(nodes, n)
		assert.True(t, hasNode)
	}
}

func TestHttpServiceNodes(t *testing.T) {
	nodes := httpServiceNodes()
	t.Run("has correct components", func(t *testing.T) {
		assert.Len(t, nodes, 7)
		expectedNodes := []string{"adapter", "app", "handler", "provider", "server", "test", "main.go"}
		for _, n := range expectedNodes {
			hasNode := _hasNodeWithName(nodes, n)
			assert.True(t, hasNode)
		}
	})

	t.Run("main.go has correct imports", func(t *testing.T) {
		nn := _findNodeByName(nodes, "main.go")
		require.NotNil(t, nn)
		n, ok := nn.(fnode)
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

func _hasNodeWithName(nodes []Node, name string) bool {
	for _, n := range nodes {
		if n.Name() == name {
			return true
		}
	}
	return false
}

func _findNodeByName(nodes []Node, name string) Node {
	for _, n := range nodes {
		if n.Name() == name {
			return n
		}
	}
	return nil
}
