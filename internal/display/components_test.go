package display_test

import (
	"bytes"
	"testing"

	"github.com/antklim/chef/internal/display"
	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestComponentsList(t *testing.T) {
	t.Run("displays a components list", func(t *testing.T) {
		components := []project.Component{
			{Name: "header", Loc: "internal/header", Desc: "header component"},
			{Name: "test", Loc: "test", Desc: "project tests"},
		}

		var buf bytes.Buffer
		err := display.ComponentsList(&buf, components)
		assert.NoError(t, err)

		expected := "registered components:\n" +
			"NAME\tLOCATION\tDESCRIPTION\nheader\tinternal/header\theader component\ntest\ttest\t\tproject tests\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("displays an information message when a components list is empty", func(t *testing.T) {
		var buf bytes.Buffer
		err := display.ComponentsList(&buf, nil)
		assert.NoError(t, err)
		assert.Equal(t, "registered components:\n\tproject does not have registered components\n", buf.String())
	})
}

func TestComponentsEmploy(t *testing.T) {
	var buf bytes.Buffer
	err := display.ComponentsEmploy(&buf, "health.go", "http_handler")
	assert.NoError(t, err)
	assert.Equal(t, `successfully added "health.go" as "http_handler" component`+"\n", buf.String())
}
