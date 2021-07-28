package project

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/antklim/chef/internal/layout"
	"github.com/pkg/errors"
)

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

// TODO: split project into different project types, move concrete types to subdirectories
// type API interface {
// 	Bootstrap() error
// 	Add(Component) error // add component

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
	errEmptyProjectName = errors.New("project name required: empty name provided")
)

type projectOptions struct {
	root string
	cat  string
	srv  string
	mod  string
}

var defaultProjectOptions = projectOptions{
	cat: defaultCategory,
	srv: defaultServer,
}

// Project manager.
type Project struct {
	name string
	opts projectOptions
	lout layout.Layout
}

// New project.
func New(name string, opt ...Option) Project {
	name = strings.TrimSpace(name)
	opts := defaultProjectOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	p := Project{
		name: name,
		opts: opts,
	}
	return p
}

func (p Project) Validate() error {
	if p.name == "" {
		return errEmptyProjectName
	}

	if c := category(p.opts.cat); c == categoryUnknown {
		return fmt.Errorf("project category %s is unknown", p.opts.cat)
	}

	if s := server(p.opts.srv); s == serverUnknown {
		return fmt.Errorf("project server %s is unknown", p.opts.srv)
	}

	root, err := p.root()
	if err != nil {
		return err
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
		return fmt.Errorf("file or directory %s already exists", p.name)
	}

	return nil
}

// Bootstrap orchestrates project validation and initialization steps.
func (p Project) Bootstrap() error {
	if err := p.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}
	if err := p.build(); err != nil {
		return errors.Wrap(err, "build failed")
	}
	return nil
}

// TODO: implement add component

// Add adds a new component node to a project
func (p Project) Add(component, name string) error {
	// TODO: add node name extension based on project language preferences

	if err := p.lout.AddComponent(component, name); err != nil {
		return errors.Wrap(err, "could not add layout component")
	}
	return errors.New("not implemented")
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Location() (string, error) {
	r, err := p.root()
	if err != nil {
		return "", err
	}
	return path.Join(r, p.name), nil
}

func (p Project) build() error {
	l, err := p.layout()
	if err != nil {
		return err
	}

	loc, err := p.Location()
	if err != nil {
		return err
	}

	var dp fs.FileMode = 0755
	if err := os.Mkdir(loc, dp); err != nil {
		return err
	}

	return l.Build(loc, p.opts.mod)
}

func (p Project) root() (root string, err error) {
	root = p.opts.root
	if root == "" {
		root, err = os.Getwd()
	}
	return
}

func (p Project) layout() (*layout.Layout, error) {
	ln := category(p.opts.cat)

	if s := server(p.opts.srv); s != serverNone {
		ln += "_" + s
	}

	l := layout.Get(ln)
	if l == nil {
		return nil, fmt.Errorf("not found layout for category %s", p.opts.cat)
	}

	return l, nil
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
