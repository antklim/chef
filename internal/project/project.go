package project

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/antklim/chef/internal/layout"
)

type Category string

const (
	CategoryApp Category = "app"
	CategoryPkg Category = "pkg"
)

type Server string

const (
	ServerHTTP Server = "http"
	ServerGRPC Server = "grpc"
)

const (
	defaultCategory = CategoryApp
	defaultServer   = ServerHTTP
)

type projectOptions struct {
	root string
	cat  Category
	srv  Server
}

var defaultProjectOptions = projectOptions{
	cat: defaultCategory,
	srv: defaultServer,
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
		return errors.New("project name required: empty name provided")
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

// Init initializes the project layout.
func (p Project) Init() error {
	rl := layout.Node{
		Name:     p.name,
		Children: layout.Default,
	}

	root, err := p.root()
	if err != nil {
		return err
	}

	if err := layout.Builder(root, rl); err != nil {
		return err
	}

	return nil
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Root() string {
	return p.opts.root
}

func (p Project) Category() Category {
	return p.opts.cat
}

func (p Project) Server() Server {
	return p.opts.srv
}

func (p Project) root() (root string, err error) {
	root = p.opts.root
	if root == "" {
		root, err = os.Getwd()
	}
	return
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

func WithCategory(c Category) Option {
	return newFuncOption(func(o *projectOptions) {
		o.cat = c
	})
}

func WithServer(s Server) Option {
	return newFuncOption(func(o *projectOptions) {
		o.srv = s
	})
}
