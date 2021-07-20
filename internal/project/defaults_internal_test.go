package project

import (
	"io/fs"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestHttpEndpoint(t *testing.T) {
	h := httpEndpoint("health")
	assert.Equal(t, "health.go", h.Name())
	assert.Equal(t, fs.FileMode(0644), h.Permissions())
	assert.IsType(t, &template.Template{}, h.Template())
}
