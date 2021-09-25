package project

import (
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/antklim/chef/internal/chef"
	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectOptions(t *testing.T) {
	tl := layout.New()
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
		{
			desc: "project created with custom layout",
			opts: []Option{WithLayout(tl)},
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "",
				lout: tl,
			},
		},
		{
			desc: "project created from notation",
			opts: []Option{WithNotation(chef.Notation{Category: "srv", Server: "http", Module: "cheftest"})},
			expected: projectOptions{
				root: "",
				cat:  "srv",
				srv:  "http",
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

func TestProjectSetComponents(t *testing.T) {
	testCases := []struct {
		desc string
		p    *Project
		a    func(*testing.T, map[string]Component)
	}{
		{
			desc: "sets components for default project",
			p:    New("test"),
			a: func(t *testing.T, c map[string]Component) {
				assert.Empty(t, c)
			},
		},
		{
			desc: "sets components for http service project",
			p:    New("test1", WithCategory("srv"), WithServer("http")),
			a: func(t *testing.T, c map[string]Component) {
				assert.NotEmpty(t, c)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Nil(t, tC.p.lout)
			tC.p.setComponents()
			tC.a(t, tC.p.components)
		})
	}
}

func TestProjectSetLayout(t *testing.T) {
	tl := layout.New()
	testCases := []struct {
		desc string
		p    *Project
		a    func(*testing.T, *layout.Layout)
	}{
		{
			desc: "sets default project layout",
			p:    New("test"),
			a: func(t *testing.T, l *layout.Layout) {
				assert.NotNil(t, l)
				assert.NotEqual(t, tl, l)
			},
		},
		{
			desc: "sets http service layout",
			p:    New("test1", WithCategory("srv"), WithServer("http")),
			a: func(t *testing.T, l *layout.Layout) {
				assert.NotNil(t, l)
				assert.NotEqual(t, tl, l)
			},
		},
		{
			desc: "sets custom layout",
			p:    New("test2", WithLayout(tl)),
			a: func(t *testing.T, l *layout.Layout) {
				assert.Equal(t, tl, l)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Nil(t, tC.p.lout)
			err := tC.p.setLayout()
			require.NoError(t, err)
			tC.a(t, tC.p.lout)
		})
	}

	t.Run("returns error when unknown layout requested", func(t *testing.T) {
		p := New("test", WithCategory("test"))
		err := p.setLayout()
		assert.EqualError(t, err, `category "test": layout not found`)
		assert.Nil(t, p.lout)
	})
}

func TestProjectInit(t *testing.T) {
	name := "project"
	tmpDir := "/tmp"
	cwd, _ := os.Getwd()
	dloc := path.Join(cwd, name)

	testCases := []struct {
		desc string
		loc  string
		opts []Option
	}{
		{
			desc: "by default inits project in current direcoty",
			loc:  dloc,
		},
		{
			desc: "inits project with with custom location",
			loc:  path.Join(tmpDir, name),
			opts: []Option{WithRoot(tmpDir)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := New(name, tC.opts...)
			err := p.Init()
			assert.NoError(t, err)
			assert.Equal(t, tC.loc, p.loc)
		})
	}
}

func TestProjectRegisterComponent(t *testing.T) {
	t.Run("overrides an existing component", func(t *testing.T) {
		tmpl := template.Must(template.New("test").Parse("package foo"))
		l := layout.New(node.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.Init()
		require.NoError(t, err)

		c := NewComponent("http_handler", "handler", "HTTP Handler", tmpl)
		err = p.RegisterComponent(c)
		require.NoError(t, err)
		assert.Equal(t, tmpl, p.components[c.Name].Tmpl)

		otherTmpl := template.Must(template.New("test2").Parse("package bar"))
		c = NewComponent("http_handler", "handler", "HTTP Handler", otherTmpl)
		err = p.RegisterComponent(c)
		require.NoError(t, err)
		assert.Equal(t, otherTmpl, p.components[c.Name].Tmpl)
	})
}
