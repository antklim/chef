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
	assert.NotEmpty(t, buf.String())
}
