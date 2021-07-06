package template_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/antklim/chef/internal/layout/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateRegistry(t *testing.T) {
	testCases := []struct {
		desc string
		name string
	}{
		{
			desc: "has an http endpoint template",
			name: template.HTTPEndpoint,
		},
		{
			desc: "has an http router template",
			name: template.HTTPRouter,
		},
		{
			desc: "has an http server template",
			name: template.HTTPServer,
		},
		{
			desc: "has an http service template",
			name: template.HTTPService,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tmpl := template.Get(tC.name)
			assert.NotNil(t, tmpl)
		})
	}
}

func TestHttpEndpointTemplate(t *testing.T) {
	data := template.HTTPEndpointData{
		Name: "health",
		Path: "/health_path",
	}
	tmpl := template.Get(template.HTTPEndpoint)
	var out bytes.Buffer
	err := tmpl.Execute(&out, data)
	require.NoError(t, err)
	outs := out.String()

	t.Run("does not have <no value> patterns", func(t *testing.T) {
		assert.False(t, strings.Contains(outs, "<no value>"))
	})

	t.Run("declares route", func(t *testing.T) {
		assert.True(t, strings.Contains(outs, `const healthRoute = "/health_path"`))
	})

	t.Run("adds route to router", func(t *testing.T) {
		assert.True(t, strings.Contains(outs, "router.Handle(healthRoute, healthHandler())"))
	})

	t.Run("defines route handler", func(t *testing.T) {
		assert.True(t, strings.Contains(outs, "func healthHandler() http.Handler"))
	})
}
