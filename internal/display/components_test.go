package display_test

import (
	"bytes"
	"testing"

	"github.com/antklim/chef/internal/display"
	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestComponentsList(t *testing.T) {
	components := []project.Component{
		{Name: "header", Loc: "internal/header", Desc: "header component"},
		{Name: "test", Loc: "test", Desc: "project tests"},
	}

	var buf bytes.Buffer
	err := display.ComponentsList(&buf, components)
	assert.NoError(t, err)

	expected := "NAME\tLOCATION\tDESCRIPTION\nheader\tinternal/header\theader component\ntest\ttest\t\tproject tests\n"
	assert.Equal(t, expected, buf.String())
}
