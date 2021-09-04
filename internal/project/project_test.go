package project_test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/antklim/chef/internal/project"
	testapi "github.com/antklim/chef/test/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectInitFails(t *testing.T) {
	tmpDir := t.TempDir()
	foofile := path.Join(tmpDir, "foofoo")
	_, err := os.Create(foofile)
	require.NoError(t, err)

	testCases := []struct {
		desc string
		name string
		opts []project.Option
		err  string
	}{
		{
			desc: "when the project name is an empty string",
			err:  "validation failed: name cannot be empty",
		},
		{
			desc: "when the project category is unknown",
			name: "cheftest",
			opts: []project.Option{project.WithCategory("foo")},
			err:  `validation failed: unknown category "foo"`,
		},
		{
			desc: "when the project server is unknown",
			name: "cheftest",
			opts: []project.Option{project.WithServer("foo")},
			err:  `validation failed: unknown server "foo"`,
		},
		{
			desc: "when root directory does not exist",
			name: "cheftest",
			opts: []project.Option{project.WithRoot("foo")},
			err:  "set location failed: stat foo: no such file or directory",
		},
		{
			desc: "when root is not a directory",
			name: "cheftest",
			opts: []project.Option{project.WithRoot(foofile)},
			err:  `set location failed: "` + foofile + `" is not a directory`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(tC.name, tC.opts...)
			err := p.Init()
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestProjectInit(t *testing.T) {
	name := "project"
	l := layout.New()

	testCases := []struct {
		desc          string
		hasComponents bool
		opts          []project.Option
	}{
		{
			desc: "inits project with default options",
		},
		{
			desc:          "inits project with layout determied by server",
			hasComponents: true,
			opts:          []project.Option{project.WithServer("http")},
		},
		{
			desc: "inits project with custom layout",
			opts: []project.Option{project.WithLayout(l)},
		},
		{
			desc:          "inits project with custom layout taking priority over category",
			hasComponents: true,
			opts:          []project.Option{project.WithLayout(l), project.WithServer("http")},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := project.New(name, tC.opts...)
			err := p.Init()
			assert.NoError(t, err)

			if tC.hasComponents {
				assert.NotEmpty(t, p.ComponentsNames())
			} else {
				assert.Empty(t, p.ComponentsNames())
			}
		})
	}
}

func TestProjectBuildFails(t *testing.T) {
	name := "cheftest" // test project name

	tmpDir := t.TempDir()
	ppath := path.Join(tmpDir, name)
	err := os.Mkdir(ppath, 0755)
	require.NoError(t, err)

	testCases := []struct {
		desc string
		pgen func() (*project.Project, error)
		err  string
	}{
		{
			desc: "when project not inited",
			pgen: func() (*project.Project, error) {
				return project.New(name), nil
			},
			err: "project not inited",
		},
		{
			desc: "when root contains file or directory with the project name",
			pgen: func() (*project.Project, error) {
				p := project.New(name, project.WithRoot(tmpDir))
				err := p.Init()
				return p, err
			},
			err: fmt.Sprintf("build failed: mkdir %s: file exists", ppath),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p, err := tC.pgen()
			require.NoError(t, err)

			loc, err := p.Build()
			assert.EqualError(t, err, tC.err)
			assert.Empty(t, loc)
		})
	}
}

func TestProjectBuild(t *testing.T) {
	t.Run("builds project", func(t *testing.T) {
		tmpDir := t.TempDir()
		p := project.New("project", project.WithRoot(tmpDir))
		err := p.Init()
		require.NoError(t, err)

		loc, err := p.Build()
		require.NoError(t, err)

		nodes, err := os.ReadDir(loc)
		require.NoError(t, err)
		assert.NotEmpty(t, nodes)
	})
}

func TestProjectRegisterComponentFails(t *testing.T) {
	name := "cheftest" // test project name
	tmpl := template.Must(template.New("test").Parse("package foo"))

	handlerComponent := project.NewComponent("http_handler", "handler", "", tmpl)
	noTmplComponent := project.NewComponent("http_handler", "handler", "", nil)

	testCases := []struct {
		desc string
		pgen func() (*project.Project, error)
		c    project.Component
		err  string
	}{
		{
			desc: "when project not inited",
			pgen: func() (*project.Project, error) {
				return project.New(name), nil
			},
			c:   noTmplComponent,
			err: "project not inited",
		},
		{
			desc: "when template is nil",
			pgen: func() (*project.Project, error) {
				p := project.New(name)
				err := p.Init()
				return p, err
			},
			c:   noTmplComponent,
			err: "nil component template",
		},
		{
			desc: "when location does not exist",
			pgen: func() (*project.Project, error) {
				l := layout.New(node.NewDnode("dir"))
				p := project.New(name, project.WithLayout(l))
				err := p.Init()
				return p, err
			},
			c:   handlerComponent,
			err: `"handler" does not exist`,
		},
		{
			desc: "when location cannot have subnodes",
			pgen: func() (*project.Project, error) {
				l := layout.New(node.NewFnode("handler"))
				p := project.New(name, project.WithLayout(l))
				err := p.Init()
				return p, err
			},
			c:   handlerComponent,
			err: `"handler" cannot have subnodes`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p, err := tC.pgen()
			require.NoError(t, err)

			err = p.RegisterComponent(tC.c)
			require.EqualError(t, err, tC.err)
			assert.NotContains(t, p.ComponentsNames(), "http_handler")
		})
	}
}

func TestProjectRegisterComponent(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("package foo"))
	l := layout.New(node.NewDnode("handler"))
	p := project.New("project", project.WithLayout(l))
	err := p.Init()
	require.NoError(t, err)

	testCases := []struct {
		desc          string
		loc           string
		componentName string
		c             project.Component
	}{
		{
			desc:          "registers a handler",
			loc:           "handler",
			componentName: "http_handler",
			c:             project.NewComponent("http_handler", "handler", "", tmpl),
		},
		{
			desc:          "registers other handler to the same location",
			loc:           "handler",
			componentName: "grpc_hander",
			c:             project.NewComponent("grpc_hander", "handler", "", tmpl),
		},
		{
			desc:          "registers component to root location",
			loc:           layout.Root,
			componentName: "main.go",
			c:             project.NewComponent("main.go", layout.Root, "", tmpl),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err = p.RegisterComponent(tC.c)
			require.NoError(t, err)
			assert.Contains(t, p.ComponentsNames(), tC.componentName)
		})
	}
}

func TestProjectEmployComponentFails(t *testing.T) {
	t.Run("when project is not inited", func(t *testing.T) {
		p := project.New("cheftest")
		err := p.EmployComponent("http_handler", "echo.go")
		assert.EqualError(t, err, "project not inited")
	})

	testCases := []struct {
		desc string
		comp string
		name string
		err  string
	}{
		{
			desc: "when adding unknow component type",
			comp: "foo",
			name: "bar",
			err:  `unregistered component "foo"`,
		},
		{
			desc: "when trying to add unknown file extension",
			comp: "http_handler",
			name: "echo.cpp",
			err:  `unknown file extension ".cpp"`,
		},
		{
			desc: "when node name contains additional periods",
			comp: "http_handler",
			name: "echo.bravo.go",
			err:  "periods not allowed in a file name",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p, err := testapi.ProjectFactory(project.WithRoot(t.TempDir()))
			require.NoError(t, err)

			err = p.EmployComponent(tC.comp, tC.name)
			assert.EqualError(t, err, tC.err)
		})
	}

	t.Run("when project layout does not exist", func(t *testing.T) {
		p, err := testapi.ProjectFactory(project.WithRoot(t.TempDir()))
		require.NoError(t, err)
		err = p.EmployComponent("http_handler", "echo.go")
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("when component with the given name already exists", func(t *testing.T) {
		p, err := testapi.ProjectFactory(project.WithRoot(t.TempDir()))
		require.NoError(t, err)
		_, err = p.Build()
		require.NoError(t, err)

		err = p.EmployComponent("http_handler", "echo")
		assert.NoError(t, err)

		err = p.EmployComponent("http_handler", "echo")
		assert.EqualError(t, err, `failed to add node to layout: failed to add node to "handler": node "echo.go" already exists`)
	})
}

func TestProjectEmployComponent(t *testing.T) {
	t.Run("adds new component node to a project layout", func(t *testing.T) {
		p, err := testapi.ProjectFactory(project.WithRoot(t.TempDir()))
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
}
