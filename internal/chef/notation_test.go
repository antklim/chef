package chef_test

import (
	"bytes"
	"testing"

	"github.com/antklim/chef/internal/chef"
	"github.com/stretchr/testify/assert"
)

const testChefTemplate = `version: "1.0"

project:
  name: dogs-and-cats
  description: Simple HTTP service in Go
  language: go
  template: https://github.com/antklim/chef-go-template
  components:
    handler: handler/http/handler.go`

func TestNotationWrite(t *testing.T) {
	n := chef.Notation{
		Version: "1.0",
		Project: chef.Project{
			Name:        "dogs-and-cats",
			Description: "Simple HTTP service in Go",
			Language:    "go",
			Template:    "https://github.com/antklim/chef-go-template",
			Components: map[string]string{
				"handler": "handler/http/handler.go",
			},
		},
	}

	var buf bytes.Buffer
	err := n.Write(&buf)
	assert.NoError(t, err)

	assert.YAMLEq(t, testChefTemplate, buf.String())
}

func TestReadNotation(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(testChefTemplate)

	expected := chef.Notation{
		Version: "1.0",
		Project: chef.Project{
			Name:        "dogs-and-cats",
			Description: "Simple HTTP service in Go",
			Language:    "go",
			Template:    "https://github.com/antklim/chef-go-template",
			Components: map[string]string{
				"handler": "handler/http/handler.go",
			},
		},
	}
	notation, err := chef.ReadNotation(&buf)
	assert.NoError(t, err)
	assert.Equal(t, expected, notation)
}
