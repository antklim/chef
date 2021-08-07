package project

import (
	"fmt"
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testTmpl = template.Must(template.New("test").Parse("package foo"))

func testProject() (*Project, error) {
	l := layout.New("layout", layout.NewDnode("handler"))
	p := New("project", WithLayout(l))
	if err := p.setLayout(); err != nil {
		return nil, err
	}

	if err := p.RegisterComponent("http_handler", "handler", testTmpl); err != nil {
		return nil, err
	}
	return p, nil
}

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
	tl := layout.New("testLayout")
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
				lout: &tl,
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

func TestProjectSetLayout(t *testing.T) {
	tl := layout.New("testLayout")
	testCases := []struct {
		desc   string
		p      *Project
		schema string
	}{
		{
			desc:   "returns default project layout",
			p:      New("test"),
			schema: ServiceLayout,
		},
		{
			desc:   "returns http service layout",
			p:      New("test1", WithCategory("srv"), WithServer("http")),
			schema: HTTPServiceLayout,
		},
		{
			desc:   "returns custom layout",
			p:      New("test2", WithLayout(tl)),
			schema: "testLayout",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Nil(t, tC.p.lout)
			err := tC.p.setLayout()
			require.NoError(t, err)
			assert.Equal(t, tC.schema, tC.p.lout.Schema())
		})
	}

	t.Run("returns error when unknown layout requested", func(t *testing.T) {
		p := New("test", WithCategory("test"))
		err := p.setLayout()
		assert.EqualError(t, err, `layout for "test" category not found`)
		assert.Nil(t, p.lout)
	})
}

func TestProjectInit(t *testing.T) {
	name := "project"
	tmpDir := t.TempDir()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defLoc := path.Join(cwd, name)
	tl := layout.New("testLayout")

	testCases := []struct {
		desc string
		loc  string
		lout string
		opts []Option
	}{
		{
			desc: "inits project with default options",
			loc:  defLoc,
			lout: "srv",
		},
		{
			desc: "inits project with with custom location",
			loc:  path.Join(tmpDir, name),
			lout: "srv",
			opts: []Option{WithRoot(tmpDir)},
		},
		{
			desc: "inits project with layout determied by server",
			loc:  defLoc,
			lout: "srv_http",
			opts: []Option{WithServer("http")},
		},
		{
			desc: "inits project with custom layout",
			loc:  defLoc,
			lout: "testLayout",
			opts: []Option{WithLayout(tl)},
		},
		{
			desc: "inits project with custom layout taking priority over category",
			loc:  defLoc,
			lout: "testLayout",
			opts: []Option{WithLayout(tl), WithServer("http")},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := New(name, tC.opts...)
			err := p.Init()
			assert.NoError(t, err)

			assert.Equal(t, tC.loc, p.loc)
			assert.Equal(t, tC.lout, p.lout.Schema())
		})
	}
}

func TestProjectRegisterComponent(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("package foo"))

	t.Run("new project does not have registered components", func(t *testing.T) {
		p := New("project")
		assert.Empty(t, p.components)
	})

	t.Run("returns error when template is nil", func(t *testing.T) {
		p := New("project")
		componentName := "handler"
		err := p.RegisterComponent(componentName, "handler", nil)
		require.EqualError(t, err, "nil component template")
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when project does not have layout", func(t *testing.T) {
		p := New("project")
		componentName := "handler"
		err := p.RegisterComponent(componentName, "handler", tmpl)
		require.EqualError(t, err, "project does not have layout")
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when location does not exist", func(t *testing.T) {
		l := layout.New("layout", layout.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.setLayout()
		require.NoError(t, err)

		componentName := "handler"
		err = p.RegisterComponent(componentName, "other/handler", tmpl)
		assert.EqualError(t, err, `"other/handler" does not exist`)
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when location is not a directory", func(t *testing.T) {
		l := layout.New("layout", layout.NewFnode("handler"))
		p := New("project", WithLayout(l))
		err := p.setLayout()
		require.NoError(t, err)

		componentName := "handler"
		err = p.RegisterComponent(componentName, "handler", tmpl)
		assert.EqualError(t, err, `"handler" not a directory`)
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("adds component to the list of components", func(t *testing.T) {
		l := layout.New("layout", layout.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.setLayout()
		require.NoError(t, err)

		{
			// register handler
			componentName := "http_handler"
			err = p.RegisterComponent(componentName, "handler", tmpl)
			require.NoError(t, err)
			assert.Contains(t, p.components, componentName)
		}

		{
			// register other handler to the same location
			componentName := "grpc_hander"
			err = p.RegisterComponent(componentName, "handler", tmpl)
			require.NoError(t, err)
			assert.Contains(t, p.components, componentName)
		}

		{
			// register to root location
			componentName := "main.go"
			err = p.RegisterComponent(componentName, layout.Root, tmpl)
			require.NoError(t, err)
			assert.Contains(t, p.components, componentName)
		}
	})

	t.Run("overrides an existing component", func(t *testing.T) {
		l := layout.New("layout", layout.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.setLayout()
		require.NoError(t, err)

		componentName := "http_handler"
		err = p.RegisterComponent(componentName, "handler", tmpl)
		require.NoError(t, err)

		otherTmpl := template.Must(template.New("test2").Parse("package bar"))
		err = p.RegisterComponent(componentName, "handler", otherTmpl)
		require.NoError(t, err)

		cmp := p.components[componentName]
		assert.Equal(t, otherTmpl, cmp.template)
	})
}

func TestProjectEmployComponent(t *testing.T) {
	t.Run("returns error when trying to add unknow component type", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		p, err := testProject()
		require.NoError(t, err)
		err = p.Employ("foo", "bar")
		assert.EqualError(t, err, `unregistered component "foo"`)
	})

	t.Run("adds new component node to a project layout", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		p, err := testProject()
		require.NoError(t, err)
		err = p.Employ("http_handler", "echo")
		assert.NoError(t, err)
	})

	t.Run("returns error when component with the given name already exists", func(t *testing.T) {
		// TODO: validate that no new nodes added to project layout
		p, err := testProject()
		require.NoError(t, err)
		err = p.Employ("http_handler", "echo")
		assert.NoError(t, err)

		err = p.Employ("http_handler", "echo")
		assert.EqualError(t, err, `add node failed: node "handler" already has subnode "echo"`)
	})
}
