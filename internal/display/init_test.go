package display_test

import (
	"bytes"
	"testing"

	"github.com/antklim/chef/internal/display"
	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestProjectInit(t *testing.T) {
	components := []project.Component{
		{Name: "header", Loc: "internal/header", Desc: "header component"},
	}

	var buf bytes.Buffer
	err := display.ProjectInit(&buf, "/tmp/cheftest", components)
	assert.NoError(t, err)

	expected := "project successfully inited at /tmp/cheftest\n\n" +
		"registered components:\n" +
		"NAME\tLOCATION\tDESCRIPTION\nheader\tinternal/header\theader component\n"
	assert.Equal(t, expected, buf.String())
}
