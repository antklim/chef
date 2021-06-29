package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/antklim/chef/internal/layout"
)

// TODO: in imports replace chef/... with the project name

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

const (
	defaultCategory = CategoryService
	defaultServer   = ServerNone
)

var (
	errEmptyProjectName = errors.New("project name required: empty name provided")
)

type projectOptions struct {
	root string
	cat  string
	srv  string
}

var defaultProjectOptions = projectOptions{
	cat: string(defaultCategory),
	srv: string(defaultServer),
}

// Project manager.
type Project struct {
	name string
	opts projectOptions
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

	if c := NewCategory(p.opts.cat); c.IsUnknown() {
		return fmt.Errorf("project category %s is unknown", p.opts.cat)
	}

	if s := NewServer(p.opts.srv); s.IsUnknown() {
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
		return err
	}
	return p.build()
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

	return l.Build(loc)
}

func (p Project) root() (root string, err error) {
	root = p.opts.root
	if root == "" {
		root, err = os.Getwd()
	}
	return
}

func (p Project) layout() (*layout.Layout, error) {
	ln := p.opts.cat

	if p.opts.srv != "" {
		ln += "_" + p.opts.srv
	}

	l := layout.Get(ln)
	if l == nil {
		return nil, fmt.Errorf("not found layout with name %s", ln)
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
