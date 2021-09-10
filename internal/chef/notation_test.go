package chef_test

import (
	"bytes"
	"testing"

	"github.com/antklim/chef/internal/chef"
	"github.com/stretchr/testify/assert"
)

func TestNotationWrite(t *testing.T) {
	n := chef.Notation{Category: "srv", Server: "http"}

	var buf bytes.Buffer
	err := n.Write(&buf)
	assert.NoError(t, err)

	expected := `version: unknown
category: srv
server: http`
	assert.YAMLEq(t, expected, buf.String())
}

func TestReadNotation(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`version: 1.0
category: srv
server: http`)

	expected := chef.Notation{
		Category: "srv",
		Server:   "http",
	}
	notation, err := chef.ReadNotation(&buf)
	assert.NoError(t, err)
	assert.Equal(t, expected, notation)
}
