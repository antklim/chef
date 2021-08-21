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

// TODO: init project with go.mod (when Go lang selected)

// TODO: make layout and components pluggable

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
	defaultExt      = ".go" // default file extension
)

var (
	errEmptyProjectName     = errors.New("project name cannot be empty")
	errComponentTemplateNil = errors.New("nil component template")
	errNotInited            = errors.New("project not inited")
)

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
	p.setComponents()
	p.inited = true
	return nil
}

// Build creates project layout nodes.
// returns location and build error.
func (p *Project) Build() (string, error) {
	if !p.inited {
		return "", errNotInited
	}
	if err := p.build(); err != nil {
		return "", errors.Wrap(err, "build failed")
	}
	return p.loc, nil
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

// EmployComponent employs registered component to add new node to a project layout.
func (p *Project) EmployComponent(component, name string) error {
	// TODO: name should not contain more than one dot
	if !p.inited {
		return errNotInited
	}

	nname, tname := name, name // node and template element name

	// TODO: extension should be configurable
	switch ext := path.Ext(name); ext {
	case "":
		nname += defaultExt
	case defaultExt:
		tname = strings.TrimSuffix(name, defaultExt)
	default:
		return fmt.Errorf("unknown file extension %q", ext)
	}

	// TODO: add node file extension based on project language preferences
	c, ok := p.components[component]
	if !ok {
		return fmt.Errorf("unregistered component %q", component)
	}

	// TODO: nodes should be added by name. File name extensions should be added
	// at build time depending on template/component.

	n := layout.NewFnode(nname, layout.WithTemplate(c.template))
	if err := p.lout.AddNode(n, c.loc); err != nil {
		return errors.Wrap(err, "add node failed")
	}

	data := struct {
		Name, Path string
	}{
		Name: tname,
		Path: "/" + tname,
	}

	return n.Build(path.Join(p.loc, c.loc), data)
}

func (p *Project) build() error {
	var dp fs.FileMode = 0755
	if err := os.Mkdir(p.loc, dp); err != nil {
		return err
	}

	return p.lout.Build(p.loc, p.opts.mod)
}

func (p *Project) setComponents() {
	f := componentsFactory(category(p.opts.cat), server(p.opts.srv))
	if f == nil {
		return
	}

	for n, c := range f.makeComponents() {
		p.components[n] = c
	}
}

func (p *Project) setLayout() error {
	if p.opts.lout != nil {
		p.lout = p.opts.lout
		return nil
	}

	f := layoutFactory(category(p.opts.cat), server(p.opts.srv))
	if f == nil {
		return fmt.Errorf("layout for %q category not found", p.opts.cat)
	}

	p.lout = f.makeLayout()
	return nil
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

	p.loc = path.Join(root, p.name)
	return nil
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

func WithLayout(l *layout.Layout) Option {
	return newFuncOption(func(o *projectOptions) {
		o.lout = l
	})
}
