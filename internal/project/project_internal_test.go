package project

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCategory(t *testing.T) {
	testCases := []struct {
		v        string
		expected string
	}{
		{
			v:        "SRV",
			expected: categoryService,
		},
		{
			v:        "service",
			expected: categoryService,
		},
		{
			v:        "foo",
			expected: categoryUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %s category when %s provided", tC.expected, tC.v), func(t *testing.T) {
			actual := category(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestProjectServer(t *testing.T) {
	testCases := []struct {
		v        string
		expected string
	}{
		{
			v:        "",
			expected: serverNone,
		},
		{
			v:        "Http",
			expected: serverHTTP,
		},
		{
			v:        "foo",
			expected: serverUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %s server when %s provided", tC.expected, tC.v), func(t *testing.T) {
			actual := server(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestProjectOptions(t *testing.T) {
	testCases := []struct {
		desc     string
		opts     []Option
		expected projectOptions
	}{
		{
			desc: "project created with default options",
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "",
			},
		},
		{
			desc: "project created with the custom root",
			opts: []Option{WithRoot("/r")},
			expected: projectOptions{
				root: "/r",
				cat:  "srv",
				srv:  "",
			},
		},
		{
			desc: "project created with custom category",
			opts: []Option{WithCategory("cli")},
			expected: projectOptions{
				root: "",
				cat:  "cli",
				srv:  "",
			},
		},
		{
			desc: "project created with custom server",
			opts: []Option{WithServer("http")},
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "http",
			},
		},
		{
			desc: "project created with custom module",
			opts: []Option{WithModule("cheftest")},
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "",
				mod:  "cheftest",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := New("test", tC.opts...)
			assert.Equal(t, tC.expected, p.opts)
		})
	}
}

func TestLayout(t *testing.T) {
	testCases := []struct {
		desc   string
		p      Project
		schema string
	}{
		{
			desc:   "returns default project layout",
			p:      New("test"),
			schema: ServiceLayout,
		},
		{
			desc:   "returns http service layout",
			p:      New("test", WithCategory("srv"), WithServer("http")),
			schema: HTTPServiceLayout,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l, err := tC.p.layout()
			require.NoError(t, err)
			assert.Equal(t, tC.schema, l.Schema())
		})
	}

	t.Run("returns error when unknown layout requested", func(t *testing.T) {
		p := New("test", WithCategory("test"))
		l, err := p.layout()
		assert.EqualError(t, err, "not found layout for category test")
		assert.Nil(t, l)
	})
}

func TestRegisterComponent(t *testing.T) {
	// tmpl := template.Must(template.New("test").Parse("package foo"))

	t.Run("new project does not have registered components", func(t *testing.T) {
		p := New("project")
		assert.Empty(t, p.components)
	})

	// t.Run("adds component to the list of components", func(t *testing.T) {
	// 	l := New("layout", NewDnode("handler"))

	// 	componentName := "http_handler"
	// 	err := l.RegisterComponent(componentName, "handler", tmpl)
	// 	require.NoError(t, err)
	// 	assert.Contains(t, l.components, componentName)

	// 	componentName = "main.go"
	// 	err = l.RegisterComponent(componentName, Root, tmpl)
	// 	require.NoError(t, err)
	// 	assert.Contains(t, l.components, componentName)
	// })

	// t.Run("returns error when location does not exist", func(t *testing.T) {
	// 	l := New("layout", NewDnode("handler"))
	// 	componentName := "handler"
	// 	err := l.RegisterComponent(componentName, "other/handler", tmpl)
	// 	assert.EqualError(t, err, `"other/handler" does not exist`)
	// 	assert.NotContains(t, l.components, componentName)
	// })

	// t.Run("returns error when location is not a directory", func(t *testing.T) {
	// 	l := New("layout", NewFnode("handler"))
	// 	componentName := "handler"
	// 	err := l.RegisterComponent(componentName, "handler", tmpl)
	// 	assert.EqualError(t, err, `"handler" not a directory`)
	// 	assert.NotContains(t, l.components, componentName)
	// })

	// t.Run("returns error when template is nil", func(t *testing.T) {
	// 	l := New("layout", NewDnode("handler"))
	// 	componentName := "handler"
	// 	err := l.RegisterComponent(componentName, "handler", nil)
	// 	require.EqualError(t, err, "component template is nil")
	// 	assert.NotContains(t, l.components, componentName)
	// })

	// t.Run("registers components", func(t *testing.T) {
	// 	l := New("layout", NewDnode("handler"))
	// 	componentName := "http_handler"
	// 	err := l.RegisterComponent(componentName, "handler", tmpl)
	// 	require.NoError(t, err)

	// 	// registers other component to the same location
	// 	componentName = "grpc_hander"
	// 	err = l.RegisterComponent(componentName, "handler", tmpl)
	// 	require.NoError(t, err)
	// })

	// t.Run("overrides an existing component", func(t *testing.T) {
	// 	l := New("layout", NewDnode("handler"))
	// 	componentName := "http_handler"
	// 	err := l.RegisterComponent(componentName, "handler", tmpl)
	// 	require.NoError(t, err)

	// 	otherTmpl := template.Must(template.New("test2").Parse("package bar"))
	// 	err = l.RegisterComponent(componentName, "handler", otherTmpl)
	// 	require.NoError(t, err)

	// 	cmp := l.components[componentName]
	// 	assert.Equal(t, otherTmpl, cmp.template)
	// })
}
