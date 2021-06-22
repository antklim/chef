package layout

import (
	"io/fs"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
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
	assert.Len(t, nodes, 7)
	expectedNodes := []string{"adapter", "app", "handler", "provider", "server", "test", "main.go"}
	for _, n := range expectedNodes {
		hasNode := _hasNodeWithName(nodes, n)
		assert.True(t, hasNode)
	}
}

func _hasNodeWithName(nodes []Node, name string) bool {
	for _, n := range nodes {
		if n.Name() == name {
			return true
		}
	}
	return false
}
