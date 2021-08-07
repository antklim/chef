package project

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/antklim/chef/internal/layout"
	"github.com/pkg/errors"
)

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

// TODO: split project into different project types, move concrete types to subdirectories
// type API interface {
// 	Init() error
// 	Employ(Component) error // employ component

// 	Location() (string, error)
// 	Name() string
// }

type Component struct{}

func (c Component) String() string {
	return "foo"
}

const (
	categoryUnknown = "unknown"
	categoryService = "srv"
)

func category(v string) string {
	switch strings.ToLower(v) {
	case "srv", "service":
		return categoryService
	default:
		return categoryUnknown
	}
}

const (
	serverUnknown = "unknown"
	serverNone    = ""
	serverHTTP    = "http"
)

func server(v string) string {
	switch strings.ToLower(v) {
	case "":
		return serverNone
	case "http":
		return serverHTTP
	default:
		return serverUnknown
	}
}

const (
	defaultCategory = categoryService
	defaultServer   = serverNone
)

const (
	// ServiceLayout an abstract service layout name.
	ServiceLayout = "srv"
	// HTTPServiceLayout an http service layout name.
	HTTPServiceLayout = "srv_http"
)

var (
	errEmptyProjectName     = errors.New("project name cannot be empty")
	errComponentTemplateNil = errors.New("nil component template")
	errNotInited            = errors.New("project not inited")
)

type component struct {
	loc      string
	name     string
	template *template.Template
}

type projectOptions struct {
	root string
	cat  string
	srv  string
	mod  string
	lout *layout.Layout
}

var defaultProjectOptions = projectOptions{
	cat: defaultCategory,
	srv: defaultServer,
}

// Project manager.
type Project struct {
	inited     bool
	name       string
	opts       projectOptions
	loc        string
	lout       *layout.Layout
	components map[string]component
}

// New project.
func New(name string, opt ...Option) *Project {
	name = strings.TrimSpace(name)
	opts := defaultProjectOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	p := &Project{
		name:       name,
		opts:       opts,
		components: make(map[string]component),
	}
	return p
}

func (p Project) validate() error {
	if p.name == "" {
		return errEmptyProjectName
	}

	if c := category(p.opts.cat); c == categoryUnknown {
		return fmt.Errorf("project category %s is unknown", p.opts.cat)
	}

	if s := server(p.opts.srv); s == serverUnknown {
		return fmt.Errorf("project server %s is unknown", p.opts.srv)
	}

	return nil
}

// Init orchestrates project validation and build steps.
func (p *Project) Init() error {
	if err := p.validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}
	if err := p.setLocation(); err != nil {
		return errors.Wrap(err, "set location failed")
	}
	if err := p.setLayout(); err != nil {
		return errors.Wrap(err, "set layout failed")
	}
	p.inited = true
	return nil
}

// Build creates project layout nodes.
// returns location and build error.
func (p Project) Build() (string, error) {
	if !p.inited {
		return "", errNotInited
	}
	if err := p.build(); err != nil {
		return "", errors.Wrap(err, "build failed")
	}
	return "", errors.New("not implemented")
}

// Employ employs registered component to add new node to a project layout.
func (p Project) Employ(component, name string) error {
	if !p.inited {
		return errNotInited
	}

	// TODO: add node name extension based on project language preferences
	c, ok := p.components[component]
	if !ok {
		return fmt.Errorf("unregistered component %q", component)
	}

	n := layout.NewFnode(name, layout.WithTemplate(c.template))
	if err := p.lout.AddNode(n, c.loc); err != nil {
		return errors.Wrap(err, "add node failed")
	}
	return nil
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Location() string {
	return p.loc
}

func (p *Project) RegisterComponent(componentName, loc string, t *template.Template) error {
	if !p.inited {
		return errNotInited
	}

	if t == nil {
		return errComponentTemplateNil
	}

	n := p.lout.FindNode(loc)
	if n == nil {
		return fmt.Errorf("%q does not exist", loc)
	}
	if _, ok := n.(layout.Dir); !ok {
		return fmt.Errorf("%q not a directory", loc)
	}

	p.components[componentName] = component{
		loc:      loc,
		name:     componentName,
		template: t,
	}

	return nil
}

func (p Project) build() error {
	var dp fs.FileMode = 0755
	if err := os.Mkdir(p.loc, dp); err != nil {
		return err
	}

	return p.lout.Build(p.loc, p.opts.mod)
}

func (p *Project) setLocation() error {
	root := p.opts.root
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		root = cwd
	}

	fi, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", root)
	}

	fi, _ = os.Stat(path.Join(root, p.name))
	if fi != nil {
		return fmt.Errorf("file or directory %q already exists", p.name)
	}

	p.loc = path.Join(root, p.name)
	return nil
}

func (p *Project) setLayout() error {
	if p.opts.lout != nil {
		p.lout = p.opts.lout
		return nil
	}

	ln := category(p.opts.cat)

	if s := server(p.opts.srv); s != serverNone {
		ln += "_" + s
	}

	l := layout.Get(ln)
	if l == nil {
		return fmt.Errorf("layout for %q category not found", p.opts.cat)
	}

	p.lout = l
	return nil
}

type Option interface {
	apply(*projectOptions)
}

type funcOption struct {
	f func(*projectOptions)
}

func (fo *funcOption) apply(o *projectOptions) {
	fo.f(o)
}

func newFuncOption(f func(*projectOptions)) *funcOption {
	return &funcOption{f}
}

func WithRoot(r string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.root = strings.TrimSpace(r)
	})
}

func WithCategory(c string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.cat = c
	})
}

func WithServer(s string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.srv = s
	})
}

func WithModule(m string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.mod = m
	})
}

func WithLayout(l layout.Layout) Option {
	return newFuncOption(func(o *projectOptions) {
		o.lout = &l
	})
}
