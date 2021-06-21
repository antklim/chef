package template_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout/template"
	"github.com/stretchr/testify/assert"
)

func TestTemplateRegistry(t *testing.T) {
	testCases := []struct {
		desc string
		name string
	}{
		{
			desc: "has an http router template",
			name: template.HTTPRouter,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tmpl := template.Get(tC.name)
			assert.NotNil(t, tmpl)
		})
	}
}
