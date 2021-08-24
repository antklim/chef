package project

import (
	"fmt"
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
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
		a    func(*testing.T, map[string]component)
	}{
		{
			desc: "sets components for default project",
			p:    New("test"),
			a: func(t *testing.T, c map[string]component) {
				assert.Empty(t, c)
			},
		},
		{
			desc: "sets components for http service project",
			p:    New("test1", WithCategory("srv"), WithServer("http")),
			a: func(t *testing.T, c map[string]component) {
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

// TODO (ref): merge public and internal tests
func TestProjectInit(t *testing.T) {
	name := "project"
	tmpDir := t.TempDir()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defLoc := path.Join(cwd, name)
	tl := layout.New()

	testCases := []struct {
		desc          string
		loc           string
		hasComponents bool
		opts          []Option
	}{
		{
			desc: "inits project with default options",
			loc:  defLoc,
		},
		{
			desc: "inits project with with custom location",
			loc:  path.Join(tmpDir, name),
			opts: []Option{WithRoot(tmpDir)},
		},
		{
			desc:          "inits project with layout determied by server",
			loc:           defLoc,
			hasComponents: true,
			opts:          []Option{WithServer("http")},
		},
		{
			desc: "inits project with custom layout",
			loc:  defLoc,
			opts: []Option{WithLayout(tl)},
		},
		{
			desc:          "inits project with custom layout taking priority over category",
			loc:           defLoc,
			hasComponents: true,
			opts:          []Option{WithLayout(tl), WithServer("http")},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := New(name, tC.opts...)
			err := p.Init()
			assert.NoError(t, err)
			assert.Equal(t, tC.loc, p.loc)
			if tC.hasComponents {
				assert.NotEmpty(t, p.components)
			} else {
				assert.Empty(t, p.components)
			}
		})
	}
	// TODO (ref): inits project with default layout in the existing project directory
}

func TestProjectRegisterComponent(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("package foo"))

	t.Run("new project does not have registered components", func(t *testing.T) {
		p := New("project")
		assert.Empty(t, p.components)
	})

	t.Run("returns error when project not inited", func(t *testing.T) {
		p := New("project")
		componentName := "handler"
		err := p.RegisterComponent(componentName, "handler", tmpl)
		require.EqualError(t, err, "project not inited")
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when template is nil", func(t *testing.T) {
		p := New("project")
		err := p.Init()
		require.NoError(t, err)
		componentName := "handler"
		err = p.RegisterComponent(componentName, "handler", nil)
		require.EqualError(t, err, "nil component template")
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when location does not exist", func(t *testing.T) {
		l := layout.New(node.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.Init()
		require.NoError(t, err)

		componentName := "handler"
		err = p.RegisterComponent(componentName, "other/handler", tmpl)
		assert.EqualError(t, err, `"other/handler" does not exist`)
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("returns error when location is not a directory", func(t *testing.T) {
		l := layout.New(node.NewFnode("handler"))
		p := New("project", WithLayout(l))
		err := p.Init()
		require.NoError(t, err)

		componentName := "handler"
		err = p.RegisterComponent(componentName, "handler", tmpl)
		assert.EqualError(t, err, `"handler" not a directory`)
		assert.NotContains(t, p.components, componentName)
	})

	t.Run("adds component to the list of components", func(t *testing.T) {
		l := layout.New(node.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.Init()
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
		l := layout.New(node.NewDnode("handler"))
		p := New("project", WithLayout(l))
		err := p.Init()
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
	// TODO (ref): in all error cases validate that no new nodes added to project layout
	// TODO (ref): in succes cases validate node added to layout (node should not have file extensions)
	testTmpl := template.Must(template.New("test").Parse("package foo"))

	testProject := func() (*Project, error) {
		l := layout.New(node.NewDnode("handler"))
		p := New("project", WithLayout(l), WithRoot(t.TempDir()))
		if err := p.Init(); err != nil {
			return nil, err
		}

		if err := p.RegisterComponent("http_handler", "handler", testTmpl); err != nil {
			return nil, err
		}

		return p, nil
	}

	t.Run("returns error when project is not inited", func(t *testing.T) {
		p := New("project")
		err := p.EmployComponent("foo", "bar")
		assert.EqualError(t, err, "project not inited")
	})

	t.Run("returns error when trying to add unknow component type", func(t *testing.T) {
		p, err := testProject()
		require.NoError(t, err)
		err = p.EmployComponent("foo", "bar")
		assert.EqualError(t, err, `unregistered component "foo"`)
	})

	t.Run("returns error when project layout does not exist", func(t *testing.T) {
		p, err := testProject()
		require.NoError(t, err)
		err = p.EmployComponent("http_handler", "echo.go")
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("returns error when trying to add unknown file extension", func(t *testing.T) {
		p, err := testProject()
		require.NoError(t, err)
		err = p.EmployComponent("http_handler", "echo.cpp")
		assert.EqualError(t, err, `unknown file extension ".cpp"`)
	})

	t.Run("adds new component node to a project layout", func(t *testing.T) {
		p, err := testProject()
		require.NoError(t, err)
		loc, err := p.Build()
		require.NoError(t, err)

		projectRoot, err := os.ReadDir(loc)
		assert.NoError(t, err)
		assert.Len(t, projectRoot, 1)
		assert.Equal(t, projectRoot[0].Name(), "handler")
		assert.True(t, projectRoot[0].IsDir())

		handlersDir, err := os.ReadDir(path.Join(loc, "handler"))
		assert.NoError(t, err)
		assert.Empty(t, handlersDir)

		err = p.EmployComponent("http_handler", "echo")
		assert.NoError(t, err)

		handlersDir, err = os.ReadDir(path.Join(loc, "handler"))
		assert.NoError(t, err)
		assert.Len(t, handlersDir, 1)
		assert.Equal(t, "echo.go", handlersDir[0].Name())
	})

	t.Run("returns error when component with the given name already exists", func(t *testing.T) {
		p, err := testProject()
		require.NoError(t, err)
		_, err = p.Build()
		require.NoError(t, err)

		err = p.EmployComponent("http_handler", "echo")
		assert.NoError(t, err)

		err = p.EmployComponent("http_handler", "echo")
		assert.EqualError(t, err, `add node failed: node "handler" already has subnode "echo.go"`)
	})
}
