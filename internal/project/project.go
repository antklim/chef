package project

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/pkg/errors"
)

// TODO: write layout settings to yaml
// TODO: read layout settings from yaml
// TODO: test/build generated go code
// TODO: support functionality of bring your own templates
// TODO: init project with go.mod (when Go lang selected)
// TODO: make layout and components pluggable

const (
	categoryUnknown = "unknown"
	categoryService = "srv"
)

var categories = map[string]string{
	"srv":     categoryService,
	"service": categoryService,
}

func category(v string) string {
	cv := strings.ToLower(v)
	cat, ok := categories[cv]
	if !ok {
		return categoryUnknown
	}
	return cat
}

const (
	serverUnknown = "unknown"
	serverNone    = ""
	serverHTTP    = "http"
)

var servers = map[string]string{
	"":     serverNone,
	"http": serverHTTP,
}

func server(v string) string {
	sv := strings.ToLower(v)
	srv, ok := servers[sv]
	if !ok {
		return serverUnknown
	}
	return srv
}

const (
	defaultCategory = categoryService
	defaultServer   = serverNone
	defaultExt      = ".go" // default file extension
)

var (
	errEmptyProjectName     = errors.New("name cannot be empty")
	errComponentTemplateNil = errors.New("nil component template")
	errNotInited            = errors.New("project not inited")
	errInvalidNodeName      = errors.New("periods not allowed in a file name")
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

// Project stores the information to maintain components, layout and nodes.
type Project struct {
	inited     bool
	name       string
	opts       projectOptions
	loc        string
	lout       *layout.Layout
	components map[string]Component
}

// New project creates a new instance of a project.
func New(name string, opt ...Option) *Project {
	name = strings.TrimSpace(name)
	opts := defaultProjectOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	p := &Project{
		name:       name,
		opts:       opts,
		components: make(map[string]Component),
	}
	return p
}

// Init orchestrates process of project initialization.
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

// Build creates project layout and returns project location and any occurred
// build error. In case of the error the location is available, an empty string
// returned instead.
func (p *Project) Build() (string, error) {
	if !p.inited {
		return "", errNotInited
	}
	if err := p.build(); err != nil {
		return "", errors.Wrap(err, "build failed")
	}
	return p.loc, nil
}

// RegisterComponent adds a component to the project.
//
// After component registered, new layout nodes can be added to project
// using `EmployComponent` method.
func (p *Project) RegisterComponent(c Component) error {
	if !p.inited {
		return errNotInited
	}

	if c.Tmpl == nil {
		return errComponentTemplateNil
	}

	n := p.lout.FindNode(c.Loc)
	if n == nil {
		return fmt.Errorf("%q does not exist", c.Loc)
	}
	if _, ok := n.(node.Adder); !ok {
		return fmt.Errorf("%q cannot have subnodes", c.Loc)
	}

	p.components[c.Name] = c

	return nil
}

// EmployComponent employs registered component to add new node to a project
// layout.
func (p *Project) EmployComponent(component, name string) error {
	if !p.inited {
		return errNotInited
	}

	if strings.Index(name, ".") != strings.LastIndex(name, ".") {
		return errInvalidNodeName
	}

	nname, tname := name, name // node and template element name

	// TODO (feat): extension should be configurable
	switch ext := path.Ext(name); ext {
	case "":
		nname += defaultExt
	case defaultExt:
		tname = strings.TrimSuffix(name, defaultExt)
	default:
		return fmt.Errorf("unknown file extension %q", ext)
	}

	// TODO (feat): add node file extension based on project language preferences
	c, ok := p.components[component]
	if !ok {
		return fmt.Errorf("unregistered component %q", component)
	}

	// TODO (feat): nodes should be added by name. File name extensions should be added
	// at build time depending on template/component.

	n := node.NewFnode(nname, node.WithTemplate(c.Tmpl))
	if err := p.lout.AddNode(n, c.Loc); err != nil {
		return errors.Wrap(err, "failed to add node to layout")
	}

	data := struct {
		Name, Path string
	}{
		Name: tname,
		Path: "/" + tname,
	}

	return n.Build(path.Join(p.loc, c.Loc), data)
}

// Components returns a list of registered components sorted by component name.
func (p *Project) Components() []Component {
	names := make([]string, 0, len(p.components))
	for name := range p.components {
		names = append(names, name)
	}

	sort.Strings(names)

	components := make([]Component, 0, len(p.components))
	for _, name := range names {
		components = append(components, p.components[name])
	}

	return components
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
		return fmt.Errorf("category %q: layout not found", p.opts.cat)
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

	// TODO: get full path of the root (currently it uses relative path)
	fi, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", root)
	}

	p.loc = path.Join(root, p.name)
	return nil
}

func (p Project) validate() error {
	if p.name == "" {
		return errEmptyProjectName
	}

	if c := category(p.opts.cat); c == categoryUnknown {
		return fmt.Errorf("unknown category %q", p.opts.cat)
	}

	if s := server(p.opts.srv); s == serverUnknown {
		return fmt.Errorf("unknown server %q", p.opts.srv)
	}

	return nil
}

// Option sets project options such as root location, category, etc.
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

// WithRoot returns an Option that sets project root location.
func WithRoot(r string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.root = strings.TrimSpace(r)
	})
}

// WithCategory returns an Option that sets project category.
func WithCategory(c string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.cat = c
	})
}

// WithServer returns an Option that sets project server type.
func WithServer(s string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.srv = s
	})
}

// WithModule returns an Option that sets module name (for Go projects).
func WithModule(m string) Option {
	return newFuncOption(func(o *projectOptions) {
		o.mod = m
	})
}

// WithLayout returns an Option that sets project layout.
func WithLayout(l *layout.Layout) Option {
	return newFuncOption(func(o *projectOptions) {
		o.lout = l
	})
}
