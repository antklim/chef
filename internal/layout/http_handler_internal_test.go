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
